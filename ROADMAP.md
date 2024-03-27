# Roadmap

## 1. Core Functionality

### 1.1 Power Capping Configuration
- Implement the ability to create and manage power capping configurations using Custom Resource Definitions (CRDs).
- Define the schema for the power capping configuration, including the power cap limit and scale object references.
- Validate the power capping configuration against the defined schema.

### 1.2 Monitoring Power Usage
- Integrate with Prometheus to retrieve power consumption metrics.
- Implement a periodic monitoring mechanism to check the power usage against the power cap limit.
- Adjust the scaling behavior of the associated scale objects based on the power usage levels.

### 1.3 Scaling Adjustments
- Retrieve the current number of replicas from the associated Kubernetes deployments.
- Implement logic to set the maximum replica count based on the power usage levels and power cap limit.
- Update the scale objects with the adjusted maximum replica count to enforce power capping.

## 2. Integration

### 2.1 Kubernetes Deployment Integration
- Enhance the operator to handle Kubernetes deployments as the scale target.
- Retrieve the current replica count from the deployment status.
- Adjust the maximum replica count of the deployments based on the power capping configuration.

### 2.2 KServe Integration
- Integrate the power capping operator with KServe.
- Ensure power capping is applied to KServe deployments and associated scale objects.
- Provide seamless integration with KServe components and workflows.

### 2.3 vLLM Integration
- Extend the power capping operator to support vLLM deployments.
- Apply power capping to vLLM deployments and associated scale objects in a similar manner to KServe.
- Ensure compatibility and seamless integration with vLLM components.

### 2.4 Real-Time Carbon Intensity Integration
- Identiy and integrate real-time carbon intensity data into the power capping operator.
- Fetch carbon intensity data from external APIs or data sources.
- Calculate carbon emission based on power usage and carbon intensity.
- Dynamically adjust the power cap based on carbon intensity to achieve target carbon capping.

## 3. Enhancements

### 3.1 Power Efficiency Aware LLM Inference Routing
- Implement power efficiency aware routing for LLM inference requests.
- Utilize a Layer 7 router to route requests to LLM inference services with the highest token/watts ratio.
- Integrate with the power capping operator to optimize power efficiency and maximize token throughput.

### 3.2 GPU Frequency Tuning for Optimal Power Efficiency
- Implement GPU frequency tuning capabilities to optimize power efficiency.
- Develop a Kubernetes job to periodically tune GPU frequencies based on power consumption metrics.
- Integrate with the power capping operator to adjust GPU frequencies and maximize token throughput within the power cap limit.

### 3.3 Power Usage Smoothing
- Enhance the power capping operator to support power usage smoothing techniques.
- Implement workload-aware service placement to distribute workloads evenly across power nodes.
- Utilize dynamic power profile reshaping to optimize power usage and maximize throughput.
- Integrate with KEDA for power-aware scaling of LLM inference services.

## 4. Monitoring and Observability

### 4.1 Prometheus Integration
- Integrate the power capping operator with Prometheus for monitoring power consumption metrics.
- Expose relevant metrics from the operator for Prometheus scraping.
- Utilize Prometheus queries to retrieve power usage data for decision-making.

### 4.2 Grafana Dashboard
- Develop a Grafana dashboard to visualize power consumption metrics and operator performance.
- Display real-time power usage, carbon emission, and scaling behavior of LLM inference services.
- Provide insights into power efficiency and carbon capping effectiveness.

## 5. Testing and Validation

### 5.1 Unit Testing
- Implement comprehensive unit tests for the power capping operator codebase.
- Cover core functionality, integration points, and enhancement features.
- Ensure high code coverage and maintain test quality.

### 5.2 Integration Testing
- Develop integration tests to validate the operator's behavior in a real cluster environment.
- Test integration with Kubernetes deployments, KServe, vLLM, and carbon intensity data sources.
- Verify the operator's ability to enforce power capping and achieve target carbon capping.

### 5.3 Performance Testing
- Conduct performance tests to evaluate the operator's scalability and efficiency.
- Measure the impact of power capping on LLM inference throughput and latency.
- Optimize the operator's performance based on the test results.

## 6. Documentation and User Guide

### 6.1 Operator Installation Guide
- Provide detailed instructions for installing the power capping operator in a Kubernetes cluster.
- Include prerequisites, deployment steps, and configuration options.
- Cover integration with required components and services.

### 6.2 User Guide and Examples
- Develop a comprehensive user guide for the power capping operator.
- Provide examples and tutorials on how to configure and use the operator effectively.
- Include best practices and troubleshooting tips for common scenarios.

### 6.3 API Reference
- Document the API reference for the power capping configuration CRD.
- Describe the available fields, their meanings, and usage guidelines.
- Provide examples of API usage and expected responses.

## 7. Community and Collaboration

### 7.1 Open Source Contribution
- Publish the power capping operator as an open-source project, submit to CNCF Sandbox.

### 7.2 Ecosystem Integration
- Collaborate with the Kubernetes, KServe, vLLM, and other communities to ensure seamless integration.
- Provide support for cross-project features and compatibility.
- Contribute back to the ecosystem by sharing knowledge and expertise.

## 8. Simulation and Validation
### 8.1 Power Capping Simulation with Real-World Traces
- Develop a simulation framework to evaluate the power capping capabilities using real-world power consumption traces.
- Collect power consumption data from production environments or representative workloads.
- Simulate the behavior of the power capping operator using the collected traces.
- Analyze the effectiveness of the power capping operator in managing power usage within the defined limits.
- Evaluate the impact of different timer settings and configurations on the success rate of power capping.