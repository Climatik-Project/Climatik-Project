# Power Cap Estimation and Management for Data Centers

## Introduction

With the anticipated rise in demand for data centers driven by the increasing reliance on generative AI and burgeoning data storage requirements, the energy consumption of these facilities has become a focal point of concern. This surge in energy demand places unprecedented pressure on the data centers' power supplies, local infrastructure, and the broader national power grid.

To mitigate this, data centers can leverage flexibility in job scheduling and server allocation, allowing for adjustments in power consumption. A collaborative solution involving Kepler and Climatik presents a promising approach. Kepler's role is to forecast power usage across all computing phases accurately, while Climatik optimizes the distribution of computing tasks to adhere to a predetermined energy consumption threshold.

## Challenges

The challenge arises from the dynamic nature of power caps in data centers, in practice the power cap for data centers is not a constant. It is dynamically influenced jointly by several factors:

- Data center cooling capabilities;
- Power supply capacity limited by data center circuit breakers;
- Available capacity on local distribution feeders, which is a function of the feeder's maximum capacity and demand from adjacent areas;
- Variations in peak demand charges imposed by utility companies;
- Fluctuating electricity prices, determined by time-of-use rates or wholesale market prices, based on existing power supply agreements;
- Changing emission factors associated with electricity transmission;
- Optimal demand response programs offered by utility providers are available.

Given these considerations, it's clear that the power cap for data centers must be regularly updated to reflect these diverse and changing factors.

## Objectives

Our project aims to create and validate an algorithm that dynamically optimizes data centers' power caps to minimize energy costs and reduce carbon emissions. This includes the following three consecutive objectives:

1. **Dynamic Power Cap Optimizer:** The development of an online algorithm that adjusts the power cap for a data center, taking into account the limiting factors outlined. While these factors are well-documented in contexts such as energy storage and demand scheduling, integrating them with data center operations presents unique challenges, especially given the complexity and scale of data center activities. A promising strategy involves employing data-driven methods to learn simplified models of job flexibility from historical data, circumventing the need for an exhaustive model of all operations.

2. **Price-Taker Case Study:** Here, we assume data center operations do not influence grid prices or emissions. The objective is to optimize the power cap using local constraints and historical grid data. Such data is widely available in regions with deregulated electricity markets, including but not limited to California, Texas, and the Northeastern United States. Depending on the specific market, the data has a time resolution of five minutes to one hour.

3. **Grid-Influencer Case Study:** This study incorporates a grid simulation to assess the impact of data center operations on the overall power system, including costs, market prices, and emissions. We will use publicly available grid data, including generation mix, renewable profiles, and transmission line topologies. These data are integrated into a power system economic dispatch simulation model based on linearized optimal power flow that generates locational marginal prices along the total system cost and emission at an hourly resolution. The focus is on measuring the tangible benefits of our algorithm in reducing costs and emissions, compared to the baseline scenario.

The results of this study include a power cap optimizer prototype implementation and benchmark results in both case studies to compare the total data center energy cost and carbon footprint under various settings.
