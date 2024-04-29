# Power Capping Custom Resource Definition (CRD)

The power capping operator uses a Custom Resource Definition (CRD) to define the power capping policies and
configuration for LLM inference services. The CRD allows users to specify the power cap limit, rack-level constraints,
and other parameters to control the power consumption and workload placement.

## CRD Schema

The power capping CRD has the following schema:

```yaml
apiVersion: powercapping.example.com/v1alpha1
kind: PowerCappingPolicy
metadata:
  name: llm-inference-power-capping
spec:
  powerCapLimit: <power-cap-limit>
  rackConstraints:
    - rackId: <rack-id>
      maxPower: <max-power>
      maxTemperature: <max-temperature>
  performancePowerRatio:
    targetRatio: <target-ratio>
    tolerancePercentage: <tolerance-percentage>
  customAlgorithm:
    name: <custom-algorithm-name>
    parameters:
      <parameter-name>: <parameter-value>
```

### Fields

- `powerCapLimit` (required): Specifies the power cap limit in watts for the LLM inference service. The operator ensures
  that the total power consumption of the inference service does not exceed this limit.

- `rackConstraints` (optional): Defines the rack-level constraints for workload placement.
    - `rackId`: The unique identifier of the rack.
    - `maxPower`: The maximum power capacity of the rack in watts.
    - `maxTemperature`: The maximum allowed temperature for the rack in degrees Celsius.

- `performancePowerRatio` (optional): Specifies the target performance-power ratio and tolerance for the LLM inference
  service.
    - `targetRatio`: The desired ratio of performance to power consumption.
    - `tolerancePercentage`: The acceptable deviation from the target ratio in percentage.

- `customAlgorithm` (optional): Allows users to specify a custom algorithm for power capping and workload placement.
    - `name`: The name of the custom algorithm.
    - `parameters`: A key-value pair of parameters required by the custom algorithm.

## Examples

Here are a few examples of how to use the power capping CRD:

1. Basic power capping policy:
   ```yaml
   apiVersion: powercapping.example.com/v1alpha1
   kind: PowerCappingPolicy
   metadata:
     name: llm-inference-power-capping
   spec:
     powerCapLimit: 1000
   ```
   This example sets a power cap limit of 1000 watts for the LLM inference service.

2. Power capping policy with rack constraints:
   ```yaml
   apiVersion: powercapping.example.com/v1alpha1
   kind: PowerCappingPolicy
   metadata:
     name: llm-inference-power-capping
   spec:
     powerCapLimit: 1000
     rackConstraints:
       - rackId: rack-1
         maxPower: 5000
         maxTemperature: 30
       - rackId: rack-2
         maxPower: 4000
         maxTemperature: 28
   ```
   This example sets a power cap limit of 1000 watts and defines rack-level constraints for two racks, specifying their
   maximum power capacity and maximum allowed temperature.

3. Power capping policy with performance-power ratio:
   ```yaml
   apiVersion: powercapping.example.com/v1alpha1
   kind: PowerCappingPolicy
   metadata:
     name: llm-inference-power-capping
   spec:
     powerCapLimit: 1000
     performancePowerRatio:
       targetRatio: 0.8
       tolerancePercentage: 10
   ```
   This example sets a power cap limit of 1000 watts and specifies a target performance-power ratio of 0.8 with a
   tolerance of 10%.

4. Power capping policy with custom algorithm:
   ```yaml
   apiVersion: powercapping.example.com/v1alpha1
   kind: PowerCappingPolicy
   metadata:
     name: llm-inference-power-capping
   spec:
     powerCapLimit: 1000
     customAlgorithm:
       name: energy-aware-scheduling
       parameters:
         energyThreshold: 500
         migrationInterval: 300
   ```
   This example sets a power cap limit of 1000 watts and specifies a custom algorithm named "energy-aware-scheduling"
   with parameters for energy threshold and migration interval.

These examples demonstrate how to define power capping policies using the CRD. Users can customize the policies based on
their specific requirements and constraints.

## Applying the CRD

To apply the power capping CRD, save the desired configuration in a YAML file (e.g., `power-capping-policy.yaml`) and
apply it using kubectl:

```
kubectl apply -f power-capping-policy.yaml
```

The power capping operator will detect the newly created or updated CRD and apply the specified policies to the LLM
inference service.

For more information on using kubectl to manage CRDs, refer to
the [Kubernetes documentation](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/).