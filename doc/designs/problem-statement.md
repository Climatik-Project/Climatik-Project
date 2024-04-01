# Problem Statement: Maximizing LLM Service Performance under Power Constraint

## Objective Function

The objective is to maximize the sum of the performance of all the LLM services, where the performance is measured by the token throughput divided by the request concurrency of each LLM service.

Let:
- $n$ be the number of LLM services
- $T_i$ be the token throughput of the $i$-th LLM service
- $C_i$ be the request concurrency of the $i$-th LLM service

The objective function can be expressed as:

$$
\text{maximize} \sum_{i=1}^{n} \frac{T_i}{C_i}
$$

## Power Constraint

The optimization problem is subject to a power constraint, which limits the total power consumption of all the LLM services.

Let:
- $P_i$ be the power consumption of the $i$-th LLM service
- $P_{max}$ be the maximum allowed power consumption

The power constraint can be expressed as:

$$
\sum_{i=1}^{n} P_i \leq P_{max}
$$

## Decision Variables

The decision variables in this optimization problem are the resource allocations and configurations of each LLM service, which impact the token throughput and power consumption.

Let:
- $R_i$ be the resource allocation (e.g., number of GPUs, CPU cores, memory) for the $i$-th LLM service
- $F_i$ be the configuration (e.g., GPU frequency, batch size) for the $i$-th LLM service

The token throughput $T_i$ and power consumption $P_i$ of each LLM service are functions of the resource allocation $R_i$ and configuration $F_i$:

$$
T_i = f(R_i, F_i) \\
P_i = g(R_i, F_i)
$$

## Optimization Problem Formulation

Combining the objective function, power constraint, and decision variables, the optimization problem can be formulated as:

$$
\begin{align*}
\text{maximize} & \sum_{i=1}^{n} \frac{T_i}{C_i} \\
\text{subject to} & \sum_{i=1}^{n} P_i \leq P_{max} \\
& T_i = f(R_i, F_i) \\
& P_i = g(R_i, F_i) \\
& R_i \in \text{feasible resource allocations} \\
& F_i \in \text{feasible configurations}
\end{align*}
$$

The goal is to find the optimal resource allocations $R_i$ and configurations $F_i$ for each LLM service that maximize the sum of the performance ratios $\frac{T_i}{C_i}$, while ensuring that the total power consumption does not exceed the maximum allowed power consumption $P_{max}$.

This optimization problem is a complex one, as the token throughput and power consumption functions $f(R_i, F_i)$ and $g(R_i, F_i)$ are typically non-linear and depend on various factors such as the LLM model architecture, hardware characteristics, and workload patterns. Solving this optimization problem requires sophisticated modeling techniques and optimization algorithms to explore the large solution space efficiently.