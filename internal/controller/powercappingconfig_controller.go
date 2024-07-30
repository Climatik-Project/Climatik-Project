/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package controller

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	prom_api "github.com/prometheus/client_golang/api"
	prom_v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"

	"github.com/Climatik-Project/Climatik-Project/api/v1alpha1"
	powercappingv1alpha1 "github.com/Climatik-Project/Climatik-Project/api/v1alpha1"
	"github.com/Climatik-Project/Climatik-Project/internal/alert"
)

const (
	labelKey              = "powercapping.climatik.io"
	defaultPowerCapHigh   = 90
	defaultPowerCapMedium = 80
	defaultPowerCapLow    = 50
)

var (
	PrometheusURL = getEnv("PROM_URL", "http://prometheus:9090")
)

// PowerCappingConfigReconciler reconciles a PowerCappingConfig object
type PowerCappingConfigReconciler struct {
	client.Client
	Scheme           *runtime.Scheme
	Log              logr.Logger
	PodInformer      cache.SharedIndexInformer
	PrometheusClient prom_v1.API
	AlertService     *alert.AlertService
}

//+kubebuilder:rbac:groups=powercapping.climatik-project.ai,resources=powercappingconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=powercapping.climatik-project.ai,resources=powercappingconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=powercapping.climatik-project.ai,resources=powercappingconfigs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *PowerCappingConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("powercappingconfig", req.NamespacedName)

	// Fetch the PowerCappingConfig instance
	powerCappingConfig := &powercappingv1alpha1.PowerCappingConfig{}
	err := r.Get(ctx, req.NamespacedName, powerCappingConfig)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("PowerCappingConfig resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get PowerCappingConfig")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *PowerCappingConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Create a new Kubernetes client
	kubeClient, err := kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		return err
	}

	// Create a new shared informer factory
	factory := informers.NewSharedInformerFactory(kubeClient, 0)

	// Create a new pod informer
	r.PodInformer = factory.Core().V1().Pods().Informer()

	// Add event handlers to the pod informer
	r.PodInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: r.handlePodAdd,
		UpdateFunc: func(oldObj, newObj interface{}) {
			r.handlePodAdd(newObj)
		},
		DeleteFunc: r.handlePodDelete,
	})

	// Add the shared informer factory to the manager's runnable list
	err = mgr.Add(manager.RunnableFunc(func(ctx context.Context) error {
		factory.Start(ctx.Done())
		return nil
	}))
	if err != nil {
		return err
	}

	// Wait for the caches to be synced before starting the reconciler
	if ok := cache.WaitForCacheSync(context.Background().Done(), r.PodInformer.HasSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}
	promClient, err := prom_api.NewClient(prom_api.Config{
		Address: PrometheusURL,
	})
	if err != nil {
		return err
	}
	r.PrometheusClient = prom_v1.NewAPI(promClient)

	return ctrl.NewControllerManagedBy(mgr).
		For(&powercappingv1alpha1.PowerCappingConfig{}).
		Complete(r)
}

func (r *PowerCappingConfigReconciler) getKeplerMetrics(ctx context.Context, podName, label string) (float64, string, error) {
	query := fmt.Sprintf(`kepler_container_%s_joules_total{pod='%s'}`, podName, label)
	result, warnings, err := r.PrometheusClient.Query(ctx, query, time.Now())
	if err != nil {
		return 0, "", err
	}
	if len(warnings) > 0 {
		r.Log.Info("Prometheus query warnings", "warnings", warnings)
	}
	if result.Type() == model.ValVector {
		vector := result.(model.Vector)
		if len(vector) > 0 {
			for _, sample := range vector {
				for labelName, labelValue := range sample.Metric {
					if string(labelName) == label {
						return float64(sample.Value), string(labelValue), nil
					}
				}
			}
		}
	}
	return 0, "", fmt.Errorf("no data with gpu label returned from Prometheus query")
}

func (r *PowerCappingConfigReconciler) handlePodAdd(obj interface{}) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return
	}
	duration := time.Duration(60) * time.Second
	powerCapLabel := pod.Labels[labelKey]
	if powerCapLabel != "" {
		// retrieve power cap crd using the label
		// calculate power cap based on the peak power usage
		obj := r.Get(context.Background(), client.ObjectKey{
			Name:      powerCapLabel,
			Namespace: pod.Namespace,
		}, &powercappingv1alpha1.PowerCappingConfig{})
		if obj != nil {
			r.Log.Error(obj, "Failed to get PowerCappingConfig")
			return
		}
		powerCappingConfig, ok := obj.(*powercappingv1alpha1.PowerCappingConfig)
		if !ok {
			r.Log.Error(fmt.Errorf("Failed to cast PowerCappingConfig"), "Failed to cast PowerCappingConfig")
			return
		}

		r.Log.Info("PowerCappingConfig", "powerCapLabel", powerCapLabel)
		// fetch the observation window from the CRD
		// watch the pod power usage for the observation window
		switch powerCappingConfig.Spec.PowerCappingSpec.Kind {
		case v1alpha1.RelativePowerCapOfPeakPowerConsumptionInPercentage:
			r.Log.Info("RelativePowerCapOfPeakPowerConsumptionInPercentage")
			duration = time.Duration(powerCappingConfig.Spec.PowerCappingSpec.RelativePowerCapInPercentageSpec.SampleWindow) * time.Second
		}
	}

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		peakPower, err := r.watchPodPowerUsage(ctx, pod.Name, duration)
		if err != nil {
			r.Log.Error(err, "Failed to watch pod power usage")
			return
		}

		powerCapPercentage := getPowerCapPercentage(powerCapLabel)
		powerCap := r.calculatePowerCap(peakPower, powerCapPercentage)
		deviceLabels := r.getPodDevices(pod)
		r.createAlert(pod, int(powerCap), deviceLabels)
	}()
}

func (r *PowerCappingConfigReconciler) handlePodDelete(obj interface{}) {
	// Handle pod deletion if needed
}

func (r *PowerCappingConfigReconciler) getPodDevices(pod *corev1.Pod) map[string]string {
	devices := make(map[string]string, 2)
	ctx := context.Background()
	for _, v := range []string{"package", "gpu"} {
		_, label, err := r.getKeplerMetrics(ctx, pod.Name, v)
		if err != nil {
			r.Log.Error(err, "Failed to get Kepler metrics")
			return devices
		}
		devices[v] = label
	}
	return devices
}

func (r *PowerCappingConfigReconciler) createAlert(pod *corev1.Pod, powerCap int, deviceLabels map[string]string) error {
	return r.AlertService.SendAlert(pod.Name, powerCap, deviceLabels)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		value = fallback
	}
	return value
}

func (r *PowerCappingConfigReconciler) watchPodPowerUsage(ctx context.Context, podName string, window time.Duration) (float64, error) {
	ticker := time.NewTicker(time.Second * window)
	defer ticker.Stop()

	peakPower := float64(0)

	for {
		select {
		case <-ctx.Done():
			return peakPower, nil
		case <-ticker.C:
			currentPower, err := r.queryPodPeakPower(ctx, podName)
			if err != nil {
				return 0, err
			}
			if currentPower > peakPower {
				peakPower = currentPower
			}
		}
	}
}

func (r *PowerCappingConfigReconciler) queryPodPeakPower(ctx context.Context, podName string) (float64, error) {
	query := fmt.Sprintf(`max_over_time(rate(kepler_container_joules_total{pod='%s'}[1m]))`, podName)
	result, warnings, err := r.PrometheusClient.Query(ctx, query, time.Now())
	if err != nil {
		return 0, err
	}
	if len(warnings) > 0 {
		r.Log.Info("Prometheus query warnings", "warnings", warnings)
	}
	if result.Type() == model.ValVector {
		vector := result.(model.Vector)
		if len(vector) > 0 {
			return float64(vector[0].Value), nil
		}
	}
	return 0, fmt.Errorf("no data returned from Prometheus query")
}

func (r *PowerCappingConfigReconciler) calculatePowerCap(peakPower float64, powerCapPercentage int) float64 {
	return peakPower * float64(powerCapPercentage) / 100
}

func getPowerCapPercentage(label string) int {
	switch label {
	case "high":
		return defaultPowerCapHigh
	case "medium":
		return defaultPowerCapMedium
	case "low":
		return defaultPowerCapLow
	default:
		return 0
	}
}
