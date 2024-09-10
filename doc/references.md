# News

## [Why the AI Industry’s Thirst for New Data Centers Can’t Be Satisfied](https://www.wsj.com/tech/ai/why-the-ai-industrys-thirst-for-new-data-centers-cant-be-satisfied-93c7eff5?mod=tech_lead_pos1)

The AI industry's rush to build data centers due to increasing demand for artificial intelligence has led to significant challenges, including a shortage of necessary components, appropriate real estate, and power. Key statistics include a fivefold increase in the lead time for custom cooling systems, backup generator delivery extending to two years, and a desperate search for 15 megawatts of power by Hydra Host for a facility planning to operate 10,000 AI chips. The U.S. data-center space grew by 26% last year, with a record amount under construction, yet vacancy rates are negligible, indicating supply cannot meet demand. Amazon Web Services opens a new data center globally every three days, yet it still takes up to two years to construct large facilities. Critical components like transceivers now face extended delivery times, and labor shortages are impacting construction timelines.

## [Salesforce Calls for AI Emissions Regulations as Concerns Grow Over Tech Sector’s Carbon Footprint](https://www.wsj.com/articles/salesforce-calls-for-ai-emissions-regulations-as-concerns-grow-over-tech-sectors-carbon-footprint-dc9c016f?mod=tech_lead_pos4)

Salesforce is advocating for stricter environmental regulations on AI, emphasizing the need for transparency in emissions and energy usage associated with AI operations. The company suggests that all entities using general-purpose AI models should disclose their carbon footprint using standardized metrics. Concerns about AI's sustainability are highlighted by projections that AI energy consumption could triple, reaching 4.5% of total power generation. Notably, legislative efforts led by Senator Ed Markey aim to establish standards for assessing AI's environmental impacts, suggesting a path toward mitigating these concerns through more informed development and usage of AI technologies.


## [A Way for Energy Investors to Ride the AI Boom - WSJ](https://www.wsj.com/finance/investing/a-way-for-energy-investors-to-ride-the-ai-boom-bb0a607d?mod=hp_listc_pos2)

"McKinsey, BCG and S&P Global Commodity Insights all project electricity demand tied to data centers to increase at a compound annual growth rate of between 13% and 15% through 2030. PJM Interconnection, whose jurisdiction includes data center-heavy Virginia, expects total electricity demand to grow at an annual rate of 2.4% over the next 10 years, up from its year-ago forecast of 1.4%. 

This comes as the U.S. power market has been tightening for seven straight years, notes Steve Fleishman, equity analyst at Wolfe Research. Meanwhile, the time it takes for new capacity to go from the planning stage to commercial operation has only gotten longer as grid operators face long backlogs. 

Anyone looking to profit off the AI theme would do well to keep a basket of electricity-exposed stocks in their basket." 

## [Amid explosive demand, America is running out of power](https://www.washingtonpost.com/business/2024/03/07/ai-data-centers-power/)

The article discusses the growing power crisis in the United States, where unprecedented demand for electricity due to the rapid expansion of data centers and clean-tech manufacturing is pushing the power grid to its limits. Key statistics highlight the severity:

1. Georgia anticipates needing 18,000 megawatts more by 2040 due to data centers.
2. U.S. data centers could consume 6% of national electricity by 2026, up from 4% in 2022.
3. Utility projections for required power have nearly doubled recently.
4. The number of new transmission lines has dropped dramatically, from 4,000 miles added in 2013 to less than 1,000 miles annually now.

This surge in power demand is complicating efforts to transition to cleaner energy and is increasing reliance on aging infrastructure, threatening both economic growth and environmental targets. The situation is causing regulatory and financial conflicts over who should bear the cost of necessary grid upgrades and expansions.

## [What Nvidia's Blackwell efficiency gains mean for DC operators](https://www.theregister.com/2024/03/27/nvidia_blackwell_efficiency/?td=rt-3a)

Nvidia's new Blackwell GPUs mark a significant step forward in addressing datacenter power constraints amidst the global energy crisis. These GPUs boast power ratings up to 1,200W with unprecedented efficiency gains—about 1.7x higher than their predecessor Hopper and 3.2x that of Ampere. Despite their high power consumption, they offer increased performance per watt across various metrics. However, the substantial power and cooling requirements necessitate advanced cooling solutions like liquid cooling to manage the high density and thermal outputs effectively.

## [AWS resource restrictions point to datacenter power issues](https://www.theregister.com/2024/04/09/aws_resource_restrictions/?td=rt-3a)
AWS is facing power limitations for its datacenters in Ireland, prompting the company to direct customers to other European regions for resource-intensive operations. In Ireland, datacenters consumed 31% more power from 2021 to 2022, accounting for 18% of the country's total electricity usage. This figure is projected to rise, potentially reaching 32% by 2026. As a result, AWS cannot deploy certain high-power resources, like GPU nodes, in its Dublin facilities due to these power constraints. The power crisis affects not only AWS but other datacenters across Europe and beyond.

## [Arm CEO warns AI's power appetite could devour 25% of US electricity by 2030](https://www.wsj.com/tech/ai/artificial-intelligences-insatiable-energy-needs-not-sustainable-arm-ceo-says-a11218c9?mod=hp_lead_pos10)
The article highlights concerns over the increasing electricity consumption by AI datacenters, potentially consuming up to 25% of the U.S. power grid by 2030. Currently, AI datacenters use about 4% of U.S. electricity. This surge is largely attributed to power-intensive large language models like ChatGPT. The International Energy Agency predicts that global power consumption for AI datacenters could be ten times higher than in 2022. This could pose serious sustainability challenges unless efficiency improvements are made, even as some datacenters look towards alternative energy sources like nuclear power to manage their increasing power needs.


# Research

## [Thunderbolt: Throughput-Optimized, Quality-of-Service-Aware Power Capping at Scale](https://www.usenix.org/system/files/osdi20-li_shaohong.pdf)
" As the demand for data center capacity continues to grow,
hyperscale providers have used power oversubscription to
increase efficiency and reduce costs. Power oversubscription
requires power capping systems to smooth out the spikes that
risk overloading power equipment by throttling workloads.
Modern compute clusters run latency-sensitive serving and
throughput-oriented batch workloads on the same servers,
provisioning resources to ensure low latency for the former
while using the latter to achieve high server utilization. When
power capping occurs, it is desirable to maintain low latency
for serving tasks and throttle the throughput of batch tasks.
To achieve this, we seek a system that can gracefully throttle
batch workloads and has task-level quality-of-service (QoS)
differentiation.
In this paper we present Thunderbolt, a hardware-agnostic
power capping system that ensures safe power oversubscription while minimizing impact on both long-running
throughput-oriented tasks and latency-sensitive tasks. It uses
a two-threshold, randomized unthrottling/multiplicative decrease control policy to ensure power safety with minimized
performance degradation. It leverages the Linux kernel's CPU
bandwidth control feature to achieve task-level QoS-aware
throttling. It is robust even in the face of power telemetry unavailability. Evaluation results at the node and cluster levels
demonstrate the system's responsiveness, effectiveness for
reducing power, capability of QoS differentiation, and minimal impact on latency and task health. We have deployed this
system at scale, in multiple production clusters. As a result,
we enabled power oversubscription gains of 9%–25%, where
none was previously possible."

## [Precise Power Capping for Latency-Sensitive Applications in Datacenter (IEEE Transaction on Sustainable Computing)](https://fangmingliu.github.io/files/tsc2018-datacenter-power-capping.pdf)
"Power capping is widely used in cloud datacenters to mitigate power over-provisioning problem, thus improve datacenter
capacity and cut off their operation cost. However, inappropriate or aggressive power capping may lead to performance degradation of
applications (especially latency-sensitive ones), and there are few effective methods that can accurately evaluate and control such
negative impact caused by aggressive power capping. In this paper, we propose Fine-Grained Differential Method (FGD) to
quantitatively analyze how inappropriate power capping degrades the performance of latency-sensitive applications. By using FGD, we
can minimize the provisioned power for each server by setting a precise power budget according to application's Service Level
Agreement (SLA). And we further propose Precise Power Capping (PPCapping) which is designed to increase the datacenter capacity
with a fixed power supply by means of FGD. Our research also provides an insight of precise tradeoff between applications' SLAs and
datacenter capacity. We verify FGD and PPCapping by using real world traces from Tencent's datecenter with 25328 servers. The
experimental results show that FGD can accurately analyze the impact of power capping on the performance of latency-sensitive
applications, and PPCapping can effectively increase datacenter capacity compared with the typical power provisioning strategy"

## [Machine Learning-Based Energy-efficient Workload Management for Data Centers](https://ieeexplore.ieee.org/document/10454842/)
"Cooling costs count for a significant part of the total energy consumption in data centers, and previous researchers mainly focused on investigating thermal-ware workload distribution strategies for CPU-intensive workloads. This paper introduces a novel machine learning-based approach that aims at reducing energy consumption through thermal-aware workload distribution to build energy-efficient data centers for GPU-intensive workload. To achieve this goal, the study employs the GPUCloudSim Plus simulator, which effectively models the distribution of GPU-intensive applications under diverse workloads and utilizations. The integration of advanced machine learning models allows for accurate temperature predictions and comprehensive evaluation of the proposed algorithm's performance. We evaluated our ThermalAwareGpu workload scheduling algorithm, and it saved up to 12.82% of computing cost compared to the baseline algorithms. Our future work will explore the estimation of data center cooling energy and conduct in-depth comparisons of different workload balancing algorithms on more intensive experiments.]

## [Data Center Power Oversubscription with a Medium Voltage Power Plane and Priority-Aware Capping](https://dl.acm.org/doi/10.1145/3373376.3378533)
As major web and cloud service providers continue to accelerate the demand for new data center capacity worldwide, the importance of power oversubscription as a lever to reduce provisioning costs has never been greater. Building on insights from Google-scale deployments, we design and deploy a new architecture across hardware and software to improve power oversubscription significantly. Our design includes (1) a new medium voltage power plane to enable larger power sharing domains (across tens of MW of equipment) and (2) a scalable, fast, and robust power capping service coordinating multiple priorities of workload on every node. Over several years of production deployment, our co-design has enabled power oversubscription of 25% or higher, saving hundreds of millions of dollars of data center costs, while preserving the desired availability and performance of all workloads.

## [Optimum control method of workload placement and air conditioners in a GPU server environment](https://ieeexplore.ieee.org/document/9829613)
High-performance computing datacenters have been rapidly growing, both in number and size. In addition to the conventional CPU servers, more GPU servers are being placed in datacenters to achieve high-speed processing such as image recognition. From our experiments, the power consumption during GPU operation is greatly affected by changes in the temperature of the server intake port. Therefore, to reduce the total power consumption of datacenters, thermal management of datacenters can address dominant problems associated with cooling such as the recirculation of hot air from the server outlets to their inlets and the appearance of hot spots. In this paper, we propose a workload placement method for environments where CPU servers and GPU servers coexist, and an optimum air-conditioning control method that cooperates with the workload placement method to reduce the total power consumption of the servers and air conditioners in datacenters. Experiment results in an actual machine environment showed that our proposed method has valid power-saving effects by adjusting cooling capacity tradeoffs between GPU servers and air conditioners.

## [Sustainable Supercomputing for AI: GPU Power Capping at HPC Scale](https://dl.acm.org/doi/10.1145/3620678.3624793)
In this paper, we study the effects of power-capping GPUs at a research supercomputing center on GPU temperature and power draw; we show significant decreases in both temperature and power draw, reducing power consumption and potentially improving hardware life-span, with minimal impact on job performance. 

## [POLCA: Power Oversubscription in LLM Cloud Providers](https://arxiv.org/abs/2308.12908)
"We propose POLCA, our framework for power oversubscription that is robust, reliable, and readily deployable for GPU clusters. Using open-source models to replicate the power patterns observed in production, we simulate POLCA and demonstrate that we can deploy 30% more servers in the same GPU cluster for inference, with minimal performance loss."

## [The ugly truth behind ChatGPT: AI is guzzling resources at planet-eating rates | Mariana Mazzucato](https://www.theguardian.com/commentisfree/article/2024/may/30/ugly-truth-ai-chatgpt-guzzling-resources-environment)
The tech industry, notably through its reliance on data centers, significantly contributes to global greenhouse emissions, surpassing even commercial flights in environmental impact. Technologies like large language models require substantial energy, with implications such as high water use for cooling. The expansion of data center operations by companies like Google and Meta not only increases energy consumption but also intensifies water scarcity issues in some regions. This escalation in resource use can lead to energy shortages, affecting essential services and infrastructure development. To combat these environmental challenges, governments should enforce regulations focusing on sustainable practices and human rights in mineral supply chains. Policymakers are urged to promote less harmful business models and adopt holistic strategies to limit the tech industry's impact on the climate, aiming to keep global temperature rise below 1.5C.

## [Towards Improved Power Management in Cloud GPUs](https://homes.cs.washington.edu/~patelp1/papers/gpupower-cal23.pdf)
—As modern server GPUs are increasingly power intensive, better power management mechanisms can significantly reduce the
power consumption, capital costs, and carbon emissions in large cloud datacenters. This paper uses diverse datacenter workloads to
study the power management capabilities of modern GPUs. We find that current GPU management mechanisms have limited
compatibility and monitoring support under cloud virtualization. They have sub-optimal, imprecise, and non-intuitive implementations of
Dynamic Voltage and Frequency Scaling (DVFS) and power capping. Consequently, efficient GPU power management is not widely
deployed in clouds today. To address these limitations, we make actionable recommendations for GPU vendors and researchers.

