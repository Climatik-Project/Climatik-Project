package controller

import (
	"context"

	"github.com/Climatik-Project/Climatik-Project/internal/planner"
	"google.golang.org/grpc"
)

func calculateOptimalReplicas(name, namespace []string, powerCap int) (map[string]int32, error) {
	conn, err := grpc.Dial("planner-service:9999", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := planner.NewPlannerClient(conn)

	request := &planner.CalculateOptimalReplicasRequest{
		Deployments: make([]*planner.Deployment, len(name)),
		PowerCap:    float64(powerCap),
	}
	for i, _ := range name {
		request.Deployments[i] = &planner.Deployment{
			Name:      name[i],
			Namespace: namespace[i],
		}
	}

	response, err := client.CalculateOptimalReplicas(context.Background(), request)
	if err != nil {
		return nil, err
	}

	optimalReplicas := make(map[string]int32)
	for _, deploymentReplicas := range response.DeploymentReplicas {
		optimalReplicas[deploymentReplicas.Name] = deploymentReplicas.OptimalReplicas
	}

	return optimalReplicas, nil
}
