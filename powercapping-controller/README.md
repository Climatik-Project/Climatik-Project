# PowerCapping Controller

The PowerCapping Controller is a Kubernetes-native solution designed to dynamically check power usage and inform the Frequency Tuner DaemonSet to implement power capping for specific services or workloads in a Kubernetes environment. It uses Custom Resources (CR) called PowerCappingPolicy to define power capping limit and policies and the powercapping controller to monitor power usage and recommend actions.

## Description

This controller mainly implements:

1. Power Usage Monitor: Periodically monitors power usage and determines if capping is needed.
2. Trigger Action: Update the PowerCappingPolicy status with true if the power usage is approaching the limit.

The system uses Custom Resources to define power capping policies and trigger actions in Kubernetes clusters.

## Getting Started

### Prerequisites

- go version v1.22.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.
- Prometheus installed in your cluster for power usage monitoring.

### Deploying the Controller

1. Clone the repository and checkout the `kubeconNA` branch:
   ```sh
   git clone https://github.com/climatik-project/climatik-project.git
   cd climatik-project
   git checkout kubeconNA
   cd powercapping-controller
   ```

2. Set the image name and tag:
   ```sh
   export IMG=<your-registry>/powercapping-controller:tag
   ```

3. Build and push the Docker image:
   ```sh
   make docker-build docker-push
   ```

4. Install the Custom Resource Definitions (CRDs):
   ```sh
   make install
   ```

5. Deploy the controller with custom environment variables:
   ```sh
   kubectl create namespace climatik
   kubectl create configmap powercapping-config -n climatik \
     --from-literal=PROMETHEUS_URL=http://prometheus-server.monitoring:9090 \
     --from-literal=MONITOR_INTERVAL=1m
   
   make deploy IMG=<your-registry>/powercapping-controller:tag
   ```

   This will create a ConfigMap with the PROMETHEUS_URL and MONITOR_INTERVAL environment variables, and deploy the controller using these settings.

6. Verify the deployment:
   ```sh
   kubectl get pods -n climatik
   ```

### Using the PowerCapping Controller

1. Create a PowerCappingPolicy CR:
   ```sh
   kubectl apply -f manifests/powercappingpolicy-sample.yaml
   ```

2. Monitor the status of your PowerCappingPolicy:
   ```sh
   kubectl get powercappingpolicies
   kubectl describe powercappingpolicy <policy-name>
   ```

3. The controller will automatically monitor power usage and apply frequency changes as needed based on the defined policy.

## Uninstalling the Controller

To remove the PowerCapping Controller from your cluster:

1. Delete any PowerCappingPolicy CRs:
   ```sh
   kubectl delete -f manifests/powercappingpolicy-sample.yaml
   ```

2. Uninstall the controller:
   ```sh
   make undeploy
   ```

3. Remove the CRDs:
   ```sh
   make uninstall
   ```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

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

