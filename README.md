# Power-Capped LLM Inference Service using Kubernetes

## 1. Overview

The purpose of this project is to create a scalable and power-efficient Large Language Model (LLM) inference service
using Kubernetes. The service utilizes a custom power capping operator that accepts a Custom Resource Definition (CRD)
to specify the power capping limit. The operator
uses [KEDA (Kubernetes Event-Driven Autoscaling)](https://github.com/kedacore/keda) to scale the LLM inference service
deployment based on the specified power cap. [Kepler](https://github.com/sustainable-computing-io/kepler), a power
monitoring tool, is used to monitor the power consumption of CPU and GPU resources on the server.

In addition to server-level power capping, the operator also considers rack-level heating issues and incorporates
techniques for monitoring, capping, and scheduling workloads to reduce cooling requirements at the rack level. By
leveraging rack-aware scheduling algorithms, the operator aims to minimize heat recirculation and optimize the placement
of workloads across servers and racks.

## 2. Motivation

Data centers face the challenge of efficiently utilizing their compute resources while ensuring that power and cooling
constraints are not exceeded. Overpower and overheat incidents can lead to hardware damage, service disruptions, and
increased operational costs. This project aims to provide a solution that enables data centers to evenly distribute
workloads in time and space, reducing the risk of overpower or overheat incidents.

By implementing a power capping operator in Kubernetes, data centers can dynamically manage the power consumption of LLM
inference workloads at both the server and rack levels. The operator optimizes workload placement and resource
allocation to minimize power consumption, reduce cooling requirements, and ensure compliance with power cap limits and
rack-level constraints.

## 3. Architecture

The power capping operator follows an architecture similar to the Kubernetes Vertical Pod Autoscaler (VPA) controller.
It consists of three main components:

1. Controller: Monitors the current and past resource and power consumption, and provides recommended actions for the
   Webhook based on the defined policies.
2. Webhook: Enforce SW or HW power capping tuning.
3. Forecast Model: Forecast workload and power consumption.

```mermaid
graph TD
    A[Prometheus] --> B[Climatik Forecast Model]
    A --> C[Climatik Controller]
    D[Kepler] --> A
    E[Workload Exporters] --> A
    B --> C
    C --> F[Prometheus Alert Manager]
    C --> G[Slack]
    C --> H[GitOps]
    F --> I[Climatik Webhook]
    G --> I
    H --> I

    subgraph "OpenShift Node 1"
        J[LLM Inferencing Pod]
        K[LLM Training Pod]
        L[CPU]
        M[GPU]
    end

    subgraph "OpenShift Node 2"
        N[LLM Inferencing Pod]
        O[LLM Training Pod]
        P[CPU]
        Q[GPU]
    end

    subgraph "Power Capping Policy CRD"
        R["Power Capping Threshold: 90%<br>Power Usage Observation<br>window: 4 Hours"]
        S["Power Capping Threshold: 80%<br>Power Usage Observation<br>window: 1 Hours"]
    end

    J --> R
    K --> S
    N --> R
    O --> S

    I -->|Pod replica and resource tuning| J
    I -->|Pod replica and resource tuning| K
    I -->|Pod replica and resource tuning| N
    I -->|Pod replica and resource tuning| O
    I -->|P/C state tuning| L
    I -->|P/C state tuning| M
    I -->|P/C state tuning| P
    I -->|P/C state tuning| Q

    C -->|Power Capping Recommendation Event Posting| R
    C -->|Power Capping Recommendation Event Posting| S
```

Out of the box, the power capping operator includes batteries for Power Oversubscription and Performance-Power Ratio
Optimization scenarios. These batteries serve as examples of how the system functions in simple scenarios. Data centers
can develop or purchase more advanced algorithms from the marketplace to cover specific needs and use cases.

## 4. Installation

To install the power capping operator, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/Climatik-Project/Climatik-Project
   ```

2. Create .env file in root folder with secrets

   ```bash
   SLACK_WEBHOOK_URL=<your-slack-webhook-url>
   GITHUB_USERNAME=<your-username>
   GITHUB_REPO=<your-repo-name>
   GITHUB_PAT=<your-github-pat>
   PROMETHEUS_HOST=http://localhost:9090
   SLACK_SIGNING_SECRET=<secret> # see README-slack-webhook-server.md
   SLACK_BOT_TOKEN=<secret> # see README-slack-webhook-server.md
   ```

3. Python Libraries:

   ```bash
   deactivate
   python -m venv venv
   source venv/bin/activate
   pip install -r python/climatik_operator/requirements.txt
   ```

4. Install the necessary CRDs and operators:

   ```bash
   make cluster-up
   make
   ```

5. Verify resources (Pod, Deployment, ScaledObject) exist:

   ```bash
   kubectl get pods --all-namespaces
   kubectl get pods -n operator-powercapping-system
   kubectl get deployments -n operator-powercapping-system
   kubectl get scaledobjects -n operator-powercapping-system
   kubectl describe scaledobject mistral-7b-scaleobject -n operator-powercapping-system
   kubectl describe scaledobject llama2-7b-scaleobject -n operator-powercapping-system
   kubectl describe pod -n operator-powercapping-system operator-powercapping-controller-manager
   kubectl describe pod -n operator-powercapping-system operator-powercapping-webhook-manager
   kubectl describe pod -n operator-powercapping-system llama2-7b
   kubectl describe pod -n operator-powercapping-system mistral-7b
   ```

6. Package Visibility Issue:
   when running

   ```bash
   kubectl describe pod -n operator-powercapping-system operator-powercapping-controller-manager
   kubectl describe pod -n operator-powercapping-system operator-powercapping-webhook-manager
   ```

   if see

   ```bash
   failed to authorize: failed to fetch anonymous token: unexpected status from GET request to URL, 401 Unauthorized
   ```

   Please go to your own github and change visibility of your package to public

7. Check logs for containers:

   For manager:

   ```bash
   kubectl logs -n operator-powercapping-system operator-powercapping-controller-manager-${pod unique id} -c manager
   ```

   ```bash
   kubectl exec -it -n operator-powercapping-system deployment/llama2-7b -- /bin/sh
   ps aux
   ```

   For All:

   ```bash
   kubectl logs -n operator-powercapping-system operator-powercapping-controller-manager-${pod unique id} --all-containers=true
   ```

   For ScaleObjects:

   ```bash
   kubectl get scaledobject --all-namespaces
   kubectl logs -n keda -l app=keda-operator
   ```

8. Test Operator Locally:

   ```bash
   cd python/climatik_operator && kopf run operator.py
   ```

9. Check CRD:

   ```bash
   kubectl get crd
   ```

10. Configure the power capping CRD with the desired power cap limit, rack-level constraints, and other parameters. Refer
   to the [CRD documentation](docs/crd.md) for more details.

## 5. Usage

To reduce the risk of interrupting production workloads, data centers can initially use the power capping operator as a
pure observability and recommendation tool after installation. The operator will provide alerts and recommendations
based on the defined policies and constraints. Data center operators can manually review these recommendations and
decide whether to take the suggested actions.

The power capping operator will log the system behaviors and provide a summary and comparison of the scenarios where the
recommended actions were taken or not taken. If the recommendations are accepted, the system will simulate the behavior
of not taking the actions, and vice versa. This allows data centers to make informed decisions based on real data and
gradually adopt the power capping operator to automatically manage more workloads and use cases.

It's important to note that the power capping operator only installs the necessary CRDs and operators, and allows for
configuration of the parameters. The LLM inference services themselves are deployed and managed by other systems like
KServe and vLLM. The power capping operator will only affect the scaling behavior of these services to reach the
optimization goals, such as energy capping or efficiency.

## 6. Documentation

- [Architecture](doc/architecture.md)
- [CRD Documentation](doc/crd.md)
- [Integration with Kubernetes Tools](doc/integrations.md)
- [Custom Algorithms Marketplace](doc/marketplace.md)

## 7. Contributing

Contributions to the project are welcome! If you find any issues or have suggestions for improvement, please open an
issue or submit a pull request on the [GitHub repository](https://github.com/Climatik-Project/Climatik-Project).

## 8. License

This project is licensed under the [Apache License 2.0](LICENSE).

## 9. Contact

For any questions or inquiries, please contact the project [MAINTAINERS](MAINTAINERS.md).
