from kubernetes import client, config


def list_pods_in_deployment(namespace, deployment_name, debug=False):
    # Load the Kubernetes configuration
    config.load_kube_config()

    # Create an instance of the Kubernetes API client
    api_instance = client.AppsV1Api()

    try:
        # Get the deployment object
        deployment = api_instance.read_namespaced_deployment(
            deployment_name, namespace)
        # Get the selector for the deployment
        selector = deployment.spec.selector.match_labels
        # Create a label selector for the pods
        label_selector = ",".join(
            [f"{key}={value}" for key, value in selector.items()])
        # Get the pods matching the label selector
        pods = api_instance.list_namespaced_pod(namespace,
                                                label_selector=label_selector)
        # Print the names of the pods
        if debug:
            print(f"Pods in deployment {deployment_name}:")
            for pod in pods.items:
                print(pod.metadata.name)

        return pods

    except Exception as e:
        print(f"Error: {e}")

    return None
