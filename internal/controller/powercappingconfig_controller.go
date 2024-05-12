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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	powercappingv1alpha1 "github.com/Climatik-Project/Climatik-Project/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

const (
	labelKey            = "powercapping.climatik.io"
	resourceName        = "nvidia.com/gpu"
	defaultPowerCapHigh = 90

	defaultPowerCapMedium = 80
	defaultPowerCapLow    = 50
)

// PowerCappingConfigReconciler reconciles a PowerCappingConfig object
type PowerCappingConfigReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	Log         logr.Logger
	PodInformer cache.SharedIndexInformer
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

	// Monitor power usage
	go r.monitorPowerUsage(ctx, powerCappingConfig)

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

	return ctrl.NewControllerManagedBy(mgr).
		For(&powercappingv1alpha1.PowerCappingConfig{}).
		Complete(r)
}

func (r *PowerCappingConfigReconciler) handlePodAdd(obj interface{}) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return
	}

	// Check if the pod has the powercapping label
	powerCapLabel := pod.Labels[labelKey]
	if powerCapLabel != "" {
		// Check if the pod claims GPU resources
		containsGPU := false
		for _, container := range pod.Spec.Containers {
			if _, ok := container.Resources.Limits[resourceName]; ok {
				containsGPU = true
				break
			}
		}

		if containsGPU {
			// Handle the pod with GPU
			go r.handlePodWithGPU(pod, powerCapLabel)
		}
	}
}

func (r *PowerCappingConfigReconciler) handlePodDelete(obj interface{}) {
	// Handle pod deletion if needed
}

func (r *PowerCappingConfigReconciler) handlePodWithGPU(pod *corev1.Pod, powerCapLabel string) {
	// Get device info from the device plugin
	deviceInfo := getDeviceInfo(pod)

	// Monitor GPU power consumption and temperature
	powerConsumption, temperature := monitorGPU(deviceInfo)

	// Determine power cap percentage based on label
	powerCapPercentage := 0
	switch powerCapLabel {
	case "high":
		powerCapPercentage = defaultPowerCapHigh
	case "medium":
		powerCapPercentage = defaultPowerCapMedium
	case "low":
		powerCapPercentage = defaultPowerCapLow
	default:
		r.Log.Info("Invalid power cap label", "label", powerCapLabel)
		return
	}

	// Calculate power cap value based on observed peak power consumption
	powerCap := int(float64(powerConsumption) * float64(powerCapPercentage) / 100)

	// Set GPU power cap
	err := setGPUPowerCap(deviceInfo, powerCap)
	if err != nil {
		r.Log.Error(err, "Failed to set GPU power cap")
		return
	}

	r.Log.Info("Power capping applied", "pod", pod.Name, "label", powerCapLabel, "powerCap", powerCap, "temperature", temperature)
}

func getDeviceInfo(pod *corev1.Pod) interface{} {
	// Dummy implementation for getDeviceInfo
	deviceInfo := map[string]string{
		"gpu-1": "nvidia-gpu-1",
		"gpu-2": "nvidia-gpu-2",
	}
	return deviceInfo
}

func setGPUPowerCap(deviceInfo interface{}, powerCap int) error {
	// Dummy implementation for setGPUPowerCap
	return nil
}

func monitorGPU(deviceInfo interface{}) (int, float64) {
	// Dummy implementation for monitorGPU
	powerConsumption := 250
	temperature := 80.5
	return powerConsumption, temperature
}
