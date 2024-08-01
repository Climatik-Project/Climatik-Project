package runners

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// JobConfig holds the configuration for the Job
type JobConfig struct {
	JobName    string
	Namespace  string
	Percentage string
}

type KubernetesRunner struct {
	JobManifestPath string
}

func (r *KubernetesRunner) Run() error {
	// Load the Kubernetes configuration
	config, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("failed to load in-cluster config: %v", err)
	}

	// Create a new Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	// Read the job manifest
	manifest, err := os.ReadFile(r.JobManifestPath)
	if err != nil {
		return fmt.Errorf("failed to read job manifest: %v", err)
	}

	// Create a new JobConfig and set the percentage value
	jobConfig := JobConfig{
		Percentage: "80", //FIXME Set the percentage from alert here
	}

	// Parse the template
	tmpl, err := template.New("job").Parse(string(manifest))
	if err != nil {
		panic(err)
	}

	// Execute the template with the JobConfig data
	var jobBuffer bytes.Buffer
	if err := tmpl.Execute(&jobBuffer, jobConfig); err != nil {
		panic(err)
	}

	// Decode the Job from the YAML
	decoder := yaml.NewYAMLOrJSONDecoder(&jobBuffer, 1024)
	var job batchv1.Job
	if err := decoder.Decode(&job); err != nil {
		panic(err)
	}

	// Delete the existing job if it exists
	jobConfig.JobName = job.Name
	jobConfig.Namespace = job.Namespace
	jobsClient := clientset.BatchV1().Jobs(jobConfig.Namespace)
	err = jobsClient.Delete(context.TODO(), jobConfig.JobName, metav1.DeleteOptions{
		PropagationPolicy: func() *metav1.DeletionPropagation {
			policy := metav1.DeletePropagationForeground
			return &policy
		}(),
	})
	if err != nil {
		fmt.Printf("Job %s does not exist or could not be deleted: %v\n", jobConfig.JobName, err)
	} else {
		fmt.Printf("Deleted existing job %s\n", jobConfig.JobName)
	}
	// Create the Job in the Kubernetes cluster
	result, err := jobsClient.Create(context.TODO(), &job, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Job %s status %v", jobConfig.JobName, result)
	return nil
}
