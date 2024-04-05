package controller

import (
	"context"
	"os"
	"strconv"
	"time"

	kedav1alpha1 "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	prometheusapi "github.com/prometheus/client_golang/api"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	powercappingv1alpha1 "github.com/Climatik-Project/Climatik-Project/api/v1alpha1"
)

func (r *PowerCappingConfigReconciler) monitorPowerUsage(ctx context.Context, powerCappingConfig *powercappingv1alpha1.PowerCappingConfig) {
	log := r.Log.WithValues("powercappingconfig", powerCappingConfig.Name)

	// Retrieve the power capping configuration from the custom resource
	powerCapLimit := powerCappingConfig.Spec.PowerCapLimit
	scaledObjectRefs := powerCappingConfig.Spec.ScaledObjectRefs

	// get prometheus endpoint from env var
	//TODO set PROMETHEUS_ENDPOINT env var in the deployment or cli options
	prometheusEndpoint := os.Getenv("PROMETHEUS_ENDPOINT")
	if prometheusEndpoint == "" {
		log.Error(nil, "PROMETHEUS_ENDPOINT env var is not set")
		return
	}
	// Create a Prometheus API client
	prometheusClient, err := prometheusapi.NewClient(prometheusapi.Config{
		Address: prometheusEndpoint,
	})
	if err != nil {
		log.Error(err, "Failed to create Prometheus client")
		return
	}

	api := prometheusv1.NewAPI(prometheusClient)

	for {
		// Obtain Kepler power consumption from Prometheus
		// The query is the sum of the power consumption of the CPU package and GPUs
		query := "sum(irate(kepler_node_package_joules_total[1m]))+sum(irate(kepler_node_gpu_joules_total[1m]))"
		result, _, err := api.Query(ctx, query, time.Now())
		if err != nil {
			log.Error(err, "Failed to query Prometheus")
			time.Sleep(60 * time.Second)
			continue
		}

		powerConsumption, err := strconv.ParseFloat(result.(model.Vector)[0].Value.String(), 64)
		if err != nil {
			log.Error(err, "Failed to parse power consumption value")
			time.Sleep(60 * time.Second)
			continue
		}

		// Update the status with the current power consumption
		powerCappingConfig.Status.CurrentPowerConsumption = int(powerConsumption)
		err = r.Status().Update(ctx, powerCappingConfig)
		if err != nil {
			log.Error(err, "Failed to update PowerCappingConfig status")
		}

		// Check power usage against the power cap limit
		if powerConsumption >= float64(powerCapLimit)*0.95 {
			// Power usage is at 95% of the power cap limit
			// Set maxReplicaCount to the current number of replicas
			for _, scaledObjectRef := range scaledObjectRefs {
				scaledObject := &kedav1alpha1.ScaledObject{}
				err := r.Get(ctx, client.ObjectKey{Namespace: powerCappingConfig.Namespace, Name: scaledObjectRef.Metadata.Name}, scaledObject)
				if err != nil {
					if errors.IsNotFound(err) {
						log.Info("ScaledObject resource not found. Ignoring since object must be deleted")
						continue
					}
					log.Error(err, "Failed to get ScaledObject")
					continue
				}

				currentReplicas, err := r.getCurrentReplicas(ctx, powerCappingConfig.Namespace, scaledObject)
				if err != nil {
					log.Error(err, "Failed to get current replicas")
					continue
				}

				scaledObject.Spec.MaxReplicaCount = &currentReplicas

				err = r.Update(ctx, scaledObject)
				if err != nil {
					log.Error(err, "Failed to update ScaledObject")
				}
			}
		} else if powerConsumption >= float64(powerCapLimit)*0.8 {
			// Power usage is at 80% of the power cap limit
			// Set maxReplicaCount to one above the current number of replicas
			for _, scaledObjectRef := range scaledObjectRefs {
				scaledObject := &kedav1alpha1.ScaledObject{}
				err := r.Get(ctx, client.ObjectKey{Namespace: powerCappingConfig.Namespace, Name: scaledObjectRef.Metadata.Name}, scaledObject)
				if err != nil {
					if errors.IsNotFound(err) {
						log.Info("ScaledObject resource not found. Ignoring since object must be deleted")
						continue
					}
					log.Error(err, "Failed to get ScaledObject")
					continue
				}

				currentReplicas, err := r.getCurrentReplicas(ctx, powerCappingConfig.Namespace, scaledObject)
				if err != nil {
					log.Error(err, "Failed to get current replicas")
					continue
				}

				maxReplicaCount := currentReplicas + 1
				scaledObject.Spec.MaxReplicaCount = &maxReplicaCount

				err = r.Update(ctx, scaledObject)
				if err != nil {
					log.Error(err, "Failed to update ScaledObject")
				}
			}
		}
		//TODO update forecasted power consumption in the status

		time.Sleep(60 * time.Second)
	}
}

func (r *PowerCappingConfigReconciler) getCurrentReplicas(ctx context.Context, deploymentNamespace string, scaledObject *kedav1alpha1.ScaledObject) (int32, error) {
	deploymentName := scaledObject.Spec.ScaleTargetRef.Name

	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Namespace: deploymentNamespace, Name: deploymentName}, deployment)
	if err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("Deployment resource not found. Ignoring since object must be deleted")
			return 0, nil
		}
		r.Log.Error(err, "Failed to get Deployment")
		return 0, err
	}

	return *deployment.Spec.Replicas, nil
}

func calculateMaxReplicas(powerCapLimit int) int {
	// Implement the logic to calculate the maximum replicas based on the power cap limit
	// This is just a placeholder, replace it with your actual calculation
	return powerCapLimit / 100
}
