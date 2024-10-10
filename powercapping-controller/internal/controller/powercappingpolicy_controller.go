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
	"time"

	climatikv1alpha1 "github.com/Climatik-Project/Climatik-Project/powercapping-controller/api/v1alpha1"
	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PowerCappingPolicyReconciler reconciles a PowerCappingPolicy object
type PowerCappingPolicyReconciler struct {
	client.Client
	Log           logr.Logger
	Scheme        *runtime.Scheme
	CheckInterval time.Duration
	PrometheusURL string
}

//+kubebuilder:rbac:groups=climatik.io,resources=powercappingpolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=climatik.io,resources=powercappingpolicies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=climatik.io,resources=powercappingpolicies/finalizers,verbs=update

func (r *PowerCappingPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("powercappingpolicy", req.NamespacedName)

	var policy climatikv1alpha1.PowerCappingPolicy
	if err := r.Get(ctx, req.NamespacedName, &policy); err != nil {
		log.Error(err, "Unable to fetch PowerCappingPolicy")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Implement power usage monitoring logic
	cappingRequired, err := r.checkPowerUsage(&policy)
	if err != nil {
		log.Error(err, "Failed to check power usage")
		return ctrl.Result{}, err
	}

	// Update status
	policy.Status.CappingActionRequired = cappingRequired
	policy.Status.LastUpdated = metav1.Now()
	if err := r.Status().Update(ctx, &policy); err != nil {
		log.Error(err, "Failed to update PowerCappingPolicy status")
		return ctrl.Result{}, err
	}

	// Requeue the request to periodically check power usage
	return ctrl.Result{RequeueAfter: r.CheckInterval}, nil
}

func (r *PowerCappingPolicyReconciler) checkPowerUsage(policy *climatikv1alpha1.PowerCappingPolicy) (bool, error) {
	client, err := api.NewClient(api.Config{
		Address: r.PrometheusURL,
	})
	if err != nil {
		return false, fmt.Errorf("error creating Prometheus client: %v", err)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Replace this query with the appropriate DCGM metric for power usage
	result, warnings, err := v1api.Query(ctx, "dcgm_power_usage", time.Now())
	if err != nil {
		return false, fmt.Errorf("error querying Prometheus: %v", err)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}

	// Process the query result
	var currentPowerUsage float64
	if vector, ok := result.(model.Vector); ok {
		if len(vector) > 0 {
			currentPowerUsage = float64(vector[0].Value)
		}
	}

	cappingThreshold := float64(policy.Spec.CappingThreshold) / 100.0
	cappingLimit := float64(policy.Spec.PowerCapLimit) * cappingThreshold

	return currentPowerUsage > cappingLimit, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PowerCappingPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&climatikv1alpha1.PowerCappingPolicy{}).
		Complete(r)
}
