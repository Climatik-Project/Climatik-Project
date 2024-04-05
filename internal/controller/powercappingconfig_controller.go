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

	"github.com/go-logr/logr"
	kedav1alpha1 "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	powercappingv1alpha1 "github.com/Climatik-Project/Climatik-Project/api/v1alpha1"
)

// PowerCappingConfigReconciler reconciles a PowerCappingConfig object
type PowerCappingConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
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

	// Retrieve the power capping configuration from the custom resource
	powerCapLimit := powerCappingConfig.Spec.PowerCapLimit
	scaledObjectRefs := powerCappingConfig.Spec.ScaledObjectRefs

	// Iterate over the ScaledObjectRefs and update the KEDA ScaledObjects
	for _, scaledObjectRef := range scaledObjectRefs {
		scaledObject := &kedav1alpha1.ScaledObject{}
		err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: scaledObjectRef.Metadata.Name}, scaledObject)
		if err != nil {
			if errors.IsNotFound(err) {
				log.Info("ScaledObject resource not found. Ignoring since object must be deleted")
				continue
			}
			log.Error(err, "Failed to get ScaledObject")
			return ctrl.Result{}, err
		}

		// Update the ScaledObject with the power capping configuration
		maxReplicas := int32(calculateMaxReplicas(powerCapLimit))
		scaledObject.Spec.MaxReplicaCount = &maxReplicas

		// Update the ScaledObject in the Kubernetes cluster
		err = r.Update(ctx, scaledObject)
		if err != nil {
			log.Error(err, "Failed to update ScaledObject")
			return ctrl.Result{}, err
		}
	}

	// Monitor power usage
	go r.monitorPowerUsage(ctx, powerCappingConfig)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PowerCappingConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&powercappingv1alpha1.PowerCappingConfig{}).
		Complete(r)
}
