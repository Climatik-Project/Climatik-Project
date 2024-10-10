# Dynamic Power-Capped Service on Kubernetes

## 1. Overview

The purpose of this project is to design and implement a system that dynamically tunes GPU and CPU frequencies to achieve power capping for specific services or workloads in a Kubernetes environment. This system is particularly tailored for GPU-intensive tasks such as Large Language Model (LLM) inference services and LLM training workloads.

Key features of the system include:

1. **Dynamic Power Management**: The system continuously monitors power usage and adjusts GPU and CPU frequencies in real-time to maintain power consumption within specified limits.

2. **Workload-Specific Policies**: Different power caps can be set for various services or workloads, allowing fine-grained control over power consumption across the cluster.

3. **Kubernetes-Native Design**: The system is fully integrated with Kubernetes, using Custom Resources (CRs) and custom controllers to manage power capping policies and frequency adjustments.

4. **Flexible Algorithm Integration**: The system supports custom algorithms for determining frequency adjustments, allowing for sophisticated power management strategies.

The architecture consists of three main components:

1. **Power Usage Monitor**: Continuously monitors power consumption using data from DCGM Exporter via Prometheus, implementing a `powercapping-controller` to monitor the power usage of all resources for a service / deployment.

2. **Action Recommender**: Analyzes power usage data and recommends frequency scaling actions based on defined policies, implementing a `freqtuning-recommender` to recommend frequency scaling actions.

3. **Frequency Tuner DaemonSet**: Applies the recommended frequency changes on individual nodes, implementing a `freqtuner` to apply the recommended frequency changes on individual nodes.

These components work together to ensure that power-intensive workloads like LLM inference and training can operate efficiently within specified power constraints. The system uses two Custom Resources:

1. **PowerCappingPolicy**: Defines the power capping policy for a specific workload or service.
2. **NodeFrequencies**: Manages the frequency settings for GPUs and CPUs on specific nodes.

By dynamically adjusting GPU and CPU frequencies based on real-time power consumption data, this system enables organizations to maximize the performance of their LLM workloads while staying within power budget constraints. This approach is particularly valuable in environments where power efficiency is crucial, such as large-scale AI training clusters or edge computing scenarios running inference services.

## 2. Motivation

The motivation for this project is to address the challenges of power management in large-scale Kubernetes environments, particularly for GPU-intensive workloads like LLM inference and training. Traditional power management solutions often lack the flexibility and granularity required to optimize power consumption for these workloads. This project aims to fill this gap by providing a dynamic power capping service that can be easily integrated into existing Kubernetes clusters.

Key benefits of the PowerCappingFreqTuner project include:

1. **Power Efficiency**: By dynamically adjusting GPU and CPU frequencies, the system can significantly reduce power consumption, leading to cost savings and environmental benefits.

2. **Performance Optimization**: The system ensures that power-intensive workloads like LLM inference and training can operate efficiently within specified power constraints, maximizing performance while staying within power budget constraints.

3. **Flexibility**: The system is designed to be highly flexible, allowing for easy integration into existing Kubernetes clusters and customization of power management policies.

## 3. Architecture

The Climatik Project implements a dynamic power capping service using Kubernetes. The architecture consists of several key components:

1. Power Usage Monitor (Controller 1): A custom Kubernetes controller that monitors power usage and determines if capping is needed. It reads from and updates the PowerCappingPolicy CR, and receives data from Prometheus (fed by DCGM Exporter).
2. Action Recommender (Controller 2): Recommends scaling actions based on the power capping policy. It reads from the PowerCappingPolicy CR and creates/updates the NodeFrequencies CR with recommended actions.
3. Frequency Tuner DaemonSet (Controller 3): Applies frequency changes on individual nodes. It reads from the NodeFrequencies CR and updates its status after applying changes.

The system uses Custom Resources (CRs) to define power capping policies and manage node frequencies, providing a flexible and scalable approach to power management in Kubernetes clusters:

1. PowerCappingPolicy: A Custom Resource (CR) that defines the power capping policy for a specific workload or service. 
2. NodeFrequencies: A Custom Resource (CR) that manages the frequency settings for GPUs and CPUs on specific nodes. 

This architecture allows for dynamic power management, workload-specific policies, flexible algorithm integration, and seamless integration with Kubernetes environments, making it particularly useful for GPU-intensive workloads like LLM inference and training.

For a detailed description of the system architecture, including component interactions and workflow, please refer to our [design document](docs/design.md). This document provides:

- A system architecture diagram
- Detailed descriptions of Custom Resources (CRs)
- Explanations of the main controllers and their functions
- The overall system workflow
- Key benefits of the architecture

The design document offers a comprehensive overview of how the PowerCapping Controller works in conjunction with other components to achieve efficient power management for LLM inference workloads.

[... rest of the README.md content ...]

## 4. Installation

[... installation section ...]

## 5. Usage

[... usage section ...]

## 6. Documentation

[... documentation section ...]

## 7. Contributing

Contributions to the project are welcome! If you find any issues or have suggestions for improvement, please open an issue or submit a pull request on the GitHub repository.

For detailed information on how to contribute to this project, please refer to our [CONTRIBUTING.md](CONTRIBUTING.md) file.

**NOTE:** Run `make help` for more information on all potential `make` targets.

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html).

## 8. License

This project is licensed under the Apache License 2.0. For full details, see the [LICENSE](LICENSE) file.

## 9. Code of Conduct

The Climatik Project follows the [CNCF Code of Conduct](code-of-conduct.md).

## 10. Maintainers

For a list of project maintainers and their contact information, please see our [MAINTAINERS.md](MAINTAINERS.md) file.

## 11. Contact

For any questions or inquiries, please contact the project maintainers listed in [MAINTAINERS.md](MAINTAINERS.md).