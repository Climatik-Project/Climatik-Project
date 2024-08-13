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
	service "github.com/Climatik-Project/Climatik-Project/internal/alert"
	mockConfig "github.com/Climatik-Project/Climatik-Project/internal/alert/tests"
)

const (
	labelKey              = "climatik-project.io"
	defaultPowerCapHigh   = 90
	defaultPowerCapMedium = 80
	defaultPowerCapLow    = 50
)

var (
	PrometheusURL = getEnv("PROM_URL", "http://prometheus:9090")
	log           = ctrl.Log.WithName("controller")
)

// PowerCappingConfigReconciler reconciles a PowerCappingConfig object
type PowerCappingConfigReconciler struct {
	client.Client
	Scheme           *runtime.Scheme
	PodInformer      cache.SharedIndexInformer
	PrometheusClient prom_v1.API
	AlertService     *service.AlertService
}

//+kubebuilder:rbac:groups=climatik-project.io,resources=powercappingconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=climatik-project.io,resources=powercappingconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=climatik-project.io,resources=powercappingconfigs/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=pods/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *PowerCappingConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Fetch the PowerCappingConfig instance
	// List all the power capping configs instances
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
	log.Info("Reconcile", "powerCappingConfig:", fmt.Sprintf("powerCappingConfig %v", powerCappingConfig))
	return ctrl.Result{}, nil
}

func (r *PowerCappingConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	log.Info("Setting up PowerCappingConfigReconciler")
	// Create a new Kubernetes client
	kubeClient, err := kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		log.Info("Failed to create Kubernetes client")
		return err
	}

	// Create a new shared informer factory
	factory := informers.NewSharedInformerFactory(kubeClient, 0)
	log.Info("Shared informer factory created")
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
	log.Info("event handler created")
	// Add the shared informer factory to the manager's runnable list
	err = mgr.Add(manager.RunnableFunc(func(ctx context.Context) error {
		factory.Start(ctx.Done())
		return nil
	}))
	if err != nil {
		return err
	}
	log.Info("Pod informer created")
	// start the informer
	go factory.Start(context.Background().Done())
	log.Info("Pod informer started")
	// Wait for the caches to be synced before starting the reconciler
	if ok := cache.WaitForCacheSync(context.Background().Done(), r.PodInformer.HasSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}
	log.Info("Cache synced")
	promClient, err := prom_api.NewClient(prom_api.Config{
		Address: PrometheusURL,
	})
	if err != nil {
		return err
	}
	r.PrometheusClient = prom_v1.NewAPI(promClient)
	log.Info("Prometheus client created", "url", PrometheusURL)
	return ctrl.NewControllerManagedBy(mgr).
		For(&powercappingv1alpha1.PowerCappingConfig{}).
		Complete(r)
}

func (r *PowerCappingConfigReconciler) getKeplerMetrics(ctx context.Context, podName, device string) (float64, string, error) {
	query := fmt.Sprintf(`kepler_container_%s_joules_total{pod='%s'}`, device, podName)
	result, warnings, err := r.PrometheusClient.Query(ctx, query, time.Now())
	if err != nil {
		return 0, "", err
	}
	if len(warnings) > 0 {
		log.Info("Prometheus query warnings", "warnings", warnings)
	}
	if result.Type() == model.ValVector {
		vector := result.(model.Vector)
		if len(vector) > 0 {
			for _, sample := range vector {
				for labelName, labelValue := range sample.Metric {
					if string(labelName) == device {
						return float64(sample.Value), string(labelValue), nil
					}
				}
			}
		}
	}
	return 0, "", fmt.Errorf("no data with device %s returned from Prometheus query", device)
}

func (r *PowerCappingConfigReconciler) handlePodAdd(obj interface{}) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return
	}
	powerCapLabel := pod.Labels[labelKey]
	if powerCapLabel != "" {
		// retrieve power cap crd using the label
		// calculate power cap based on the peak power usage
		powerCappingConfig := &powercappingv1alpha1.PowerCappingConfig{}
		err := r.Get(context.Background(), client.ObjectKey{
			Name:      powerCapLabel,
			Namespace: pod.Namespace,
		}, powerCappingConfig)
		if err != nil {
			if errors.IsNotFound(err) {
				log.Error(err, fmt.Sprintf("PowerCappingConfig %s not found in ns %s", powerCapLabel, pod.Namespace))
				return
			}
			log.Error(err, fmt.Sprintf("Failed to get PowerCappingConfig: %s in ns %s", powerCapLabel, pod.Namespace))
			return
		}

		// fetch the observation window from the CRD
		// watch the pod power usage for the observation window
		switch powerCappingConfig.Spec.PowerCappingSpec.Kind {
		case v1alpha1.RelativePowerCapOfPeakPowerConsumptionInPercentage:
			log.Info("RelativePowerCapOfPeakPowerConsumptionInPercentage", "sampleWindow", powerCappingConfig.Spec.PowerCappingSpec.RelativePowerCapInPercentageSpec.SampleWindow)
			duration := time.Duration(powerCappingConfig.Spec.PowerCappingSpec.RelativePowerCapInPercentageSpec.SampleWindow) * time.Second
			go func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				log.Info("Watching pod power usage", "pod", pod.Name, "duration", duration)
				peakPower, err := r.watchPodPowerUsage(ctx, pod.Name, duration)
				if err != nil {
					log.Error(err, "Failed to watch pod power usage")
					return
				}

				powerCapPercentage := getPowerCapPercentage(powerCapLabel)
				powerCap := r.calculatePowerCap(peakPower, powerCapPercentage)
				deviceLabels := r.getPodDevices(pod)
				r.createAlert(pod, int(powerCap), deviceLabels)
			}()
		}
	}
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
			log.Error(err, "Failed to get Kepler metrics")
			return nil
		}
		devices[v] = label
	}
	return devices
}

func (r *PowerCappingConfigReconciler) createAlert(pod *corev1.Pod, powerCap int, deviceLabels map[string]string) error {
	mockConfig := mockConfig.NewMockPowerCappingConfig()

	return r.AlertService.SendAlert(pod.Name, powerCap, deviceLabels, mockConfig)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		value = fallback
	}
	return value
}

func (r *PowerCappingConfigReconciler) watchPodPowerUsage(ctx context.Context, podName string, window time.Duration) (float64, error) {
	ticker := time.NewTicker(window)
	defer ticker.Stop()

	peakPower := float64(0)

	for {
		select {
		case <-ctx.Done():
			return peakPower, nil
		case <-ticker.C:
			currentPower, err := r.queryPodPeakPower(ctx, podName, window.String())
			if err != nil {
				return 0, err
			}
			log.Info("Current power usage", "power", currentPower)
			if currentPower > peakPower {
				peakPower = currentPower
			}
			ticker.Stop()
			return peakPower, nil
		}
	}
}

func (r *PowerCappingConfigReconciler) queryPodPeakPower(ctx context.Context, podName, window string) (float64, error) {
	// sample query: max_over_time(sum(rate(kepler_container_joules_total{pod_name=~"stress-7d796cb489-fbw68"}[1m]))[1m:])
	query := fmt.Sprintf(`max_over_time(sum(rate(kepler_container_joules_total{pod_name=~"%s"}[1m]))[%s:])`, podName, window)
	result, warnings, err := r.PrometheusClient.Query(ctx, query, time.Now())
	if err != nil {
		return 0, err
	}
	if len(warnings) > 0 {
		log.Info("Prometheus query warnings", "warnings", warnings)
	}
	log.Info("Prometheus query result", "result", result)
	vectorResult, ok := result.(model.Vector)
	if !ok {
		log.Info("Prometheus query result error", "result", result)
	}

	for _, sample := range vectorResult {
		log.Info("Prometheus query sample", "sample", sample)
		return float64(sample.Value), nil
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
