# Installation Guide: Power Capping Operator with KEDA and KServe

This installation sets up the power capping operator, KEDA, KServe, and observe the power capping capabilities in action.

## Step 1: Install KEDA

1. Add the KEDA Helm repository:
```
helm repo add kedacore https://kedacore.github.io/charts
```

2. Update the Helm repository:
```
helm repo update
```

3. Install KEDA using Helm:
```
helm install keda kedacore/keda --namespace keda --create-namespace
```

4. Verify that KEDA is running:
```
kubectl get pods -n keda
```

## Step 2: Install KServe

1. Install KServe using kubectl:
```
kubectl apply -f https://github.com/kserve/kserve/releases/download/v0.12.0/kserve.yaml
```

2. Wait for the KServe pods to be ready:
```
kubectl get pods -n kserve
```

## Step 3: Install the Power Capping Operator

1. Clone the power capping operator repository:
```
git clone https://github.com/Climatik-Project/Climatik-Project
```

2. Navigate to the chart directory:
```
cd Climatik-Project/deploy/climatik-operator/helm-chart
```

3. Install the power capping operator using Helm:
```
helm install --namespace power-capping-operator power-capping-operator .
```

4. Verify that the power capping operator is running:
```
kubectl get pods -n power-capping-operator
```


## Step 4: Create a KServe Inference Service

1. Create a KServe inference service with a KEDA ScaleObject. Save the following configuration to a file named `kserve-inference-service.yaml`:
```yaml
apiVersion: serving.kserve.io/v1beta1
kind: InferenceService
metadata:
  name: my-llm-inference-service
---
apiVersion: keda.sh/v1alpha1
kind: ScaleObject
metadata:
  name: my-llm-inference-service-scaleobject
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-llm-inference-service
  pollingInterval: 15
  cooldownPeriod: 30
  minReplicaCount: 1
  maxReplicaCount: 10
  triggers:
    - type: prometheus
      metadata:
        serverAddress: http://prometheus-server.prometheus.svc.cluster.local
        metricName: token_throughput_per_second
        query: token_throughput_per_second[1m]
        threshold: 200
```

2. Apply the KServe inference service and KEDA ScaleObject:
```
kubectl apply -f kserve-inference-service.yaml
```

## Step 5: Configure Power Capping for the KServe Inference Service

1. Create a PowerCappingConfig that includes the KServe inference service. Save the following configuration to a file named `kserve-power-capping.yaml`:
```yaml
apiVersion: powercapping.example.com/v1
kind: PowerCappingConfig
metadata:
  name: kserve-power-capping
spec:
  powerCapLimit: 1000
  scaleObjectRefs:
    - apiVersion: keda.sh/v1alpha1
      kind: ScaleObject
      metadata:
        name: my-llm-inference-service-scaleobject
        namespace: default
```

2. Apply the PowerCappingConfig:
```
kubectl apply -f kserve-power-capping.yaml
```

## Step 6: Observe Power Capping in Action

1. Monitor the logs of the power capping operator to observe the power capping actions:
```
kubectl logs -l app=power-capping-operator -n power-capping-operator -f
```

2. Generate load on the KServe inference service to trigger scaling

3. Observe the scaling behavior of the KServe inference service based on the load and the power capping configuration.

4. Check the metrics and logs of the power capping operator to see how it adjusts the scaling of the KServe inference service to maintain the power cap limit.

