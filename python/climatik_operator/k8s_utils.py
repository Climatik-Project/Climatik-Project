from kubernetes import client, config


def k8s_api_instance():
    # Load the Kubernetes configuration
    config.load_kube_config()
    # Create an instance of the Kubernetes API client
    return client.AppsV1Api()


def list_pods_in_deployment(namespace, deployment_name, debug=False):
    api_instance = k8s_api_instance()

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


def list_resource_requests_in_deployments(namespace,
                                          resource_name='nvidia.com/gpu',
                                          debug=False):
    api_instance = k8s_api_instance()

    # prepare the list of pods
    pods = []
    # Get the list of deployments
    deployments = api_instance.list_namespaced_deployment(namespace)
    # Iterate over the deployments
    for deployment in deployments.items:
        # Get the pods in the deployment
        pods = list_pods_in_deployment(namespace, deployment.metadata.name,
                                       debug)

        # Iterate over the pods
        for pod in pods.items:
            # Get the containers in the pod
            containers = pod.spec.containers
            # Iterate over the containers
            for container in containers:
                # Get the resource requests for the container
                resources = container.resources.requests
                # Check if the resource name is in the resource requests
                if resource_name in resources:
                    if debug:
                        # Print the resource request
                        print(
                            f"Resource request for container {container.name} in deployment {deployment.metadata.name}: {resources[resource_name]}"
                        )
                    pods.append(pod)

    return pods
