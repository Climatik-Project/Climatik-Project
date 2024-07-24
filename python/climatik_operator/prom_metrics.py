import logging
import os
import sys
from prometheus_client import Gauge, start_http_server
from prometheus_api_client import PrometheusConnect
import requests

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[logging.StreamHandler(sys.stdout)])

# Set the logging level based on environment variable
log_level = os.getenv('LOG_LEVEL', 'INFO').upper()
logging.getLogger().setLevel(log_level)


class PowerCappingMetrics:

    def __init__(self, port=9091, url="http://localhost:9090"):
        self.logger = logging.getLogger(self.__class__.__name__)
        self.logger.info("Initializing PowerCappingMetrics")

        self.port = port
        self.prom = PrometheusConnect(url=url, disable_ssl=True)

        self.metrics = [
            'kepler_container_gpu_joules_total',
            'kepler_container_dram_joules_total',
            'kepler_container_other_joules_total',
            'kepler_container_package_joules_total',
            'kepler_container_platform_joules_total',
            'kepler_container_uncore_joules_total',
            'kepler_node_core_joules_total', 'kepler_node_other_joules_total',
            'kepler_container_bpf_block_irq_total',
            'kepler_container_bpf_cpu_time_ms_total',
            'kepler_container_bpf_net_rx_irq_total',
            'kepler_container_bpf_net_tx_irq_total',
            'kepler_container_bpf_page_cache_hit_total',
            'kepler_container_joules_total'
        ]

        self.scale_objects_gauge = Gauge(
            'power_capping_scaled_objects',
            'Number of ScaledObjects in the PowerCappingConfig CRD')
        self.replicas_gauge = Gauge('power_capping_replicas',
                                    'Number of replicas in each deployment',
                                    ['deployment'])
        self.power_consumption_gauge = Gauge(
            'power_capping_power_consumption',
            'Current power consumption of each deployment', ['deployment'])
        self.cpu_time_gauge = Gauge('power_capping_cpu_time',
                                    'CPU time usage of each deployment',
                                    ['deployment'])
        self.core_joules_gauge = Gauge(
            'power_capping_core_joules',
            'Core energy consumption of each deployment', ['deployment'])
        self.dram_joules_gauge = Gauge(
            'power_capping_dram_joules',
            'DRAM energy consumption of each deployment', ['deployment'])
        self.gpu_joules_gauge = Gauge(
            'power_capping_gpu_joules',
            'GPU energy consumption of each deployment', ['deployment'])

        # Add new gauges for the additional metrics
        self.bpf_block_irq_gauge = Gauge(
            'power_capping_bpf_block_irq',
            'BPF block IRQ total of each deployment', ['deployment'])
        self.bpf_cpu_time_ms_gauge = Gauge(
            'power_capping_bpf_cpu_time_ms',
            'BPF CPU time in ms total of each deployment', ['deployment'])
        self.bpf_net_rx_irq_gauge = Gauge(
            'power_capping_bpf_net_rx_irq',
            'BPF net RX IRQ total of each deployment', ['deployment'])
        self.bpf_net_tx_irq_gauge = Gauge(
            'power_capping_bpf_net_tx_irq',
            'BPF net TX IRQ total of each deployment', ['deployment'])
        self.bpf_page_cache_hit_gauge = Gauge(
            'power_capping_bpf_page_cache_hit',
            'BPF page cache hit total of each deployment', ['deployment'])
        self.joules_gauge = Gauge('power_capping_joules',
                                  'Joules total of each deployment',
                                  ['deployment'])

    def start_server(self):
        self.logger.info(f"Starting metrics server on port {self.port}")
        start_http_server(self.port)
        self.update_metrics()

    def update_scale_objects(self, count):
        self.logger.info(f"Updating scale objects count: {count}")
        self.scale_objects_gauge.set(count)

    def update_replicas(self, deployment, count):
        self.logger.info(f"Updating replicas for {deployment}: {count}")
        self.replicas_gauge.labels(deployment=deployment).set(count)

    def update_power_consumption(self, deployment, power_consumption):
        self.logger.info(
            f"Updating power consumption for {deployment}: {power_consumption}"
        )
        self.power_consumption_gauge.labels(
            deployment=deployment).set(power_consumption)

    def update_cpu_time(self, deployment, cpu_time):
        self.logger.info(f"Updating CPU time for {deployment}: {cpu_time}")
        self.cpu_time_gauge.labels(deployment=deployment).set(cpu_time)

    def update_core_joules(self, deployment, core_joules):
        self.logger.info(
            f"Updating core joules for {deployment}: {core_joules}")
        self.core_joules_gauge.labels(deployment=deployment).set(core_joules)

    def update_dram_joules(self, deployment, dram_joules):
        self.logger.info(
            f"Updating DRAM joules for {deployment}: {dram_joules}")
        self.dram_joules_gauge.labels(deployment=deployment).set(dram_joules)

    def update_gpu_joules(self, deployment, gpu_joules):
        self.logger.info(f"Updating GPU joules for {deployment}: {gpu_joules}")
        self.gpu_joules_gauge.labels(deployment=deployment).set(gpu_joules)

    # Add update functions for the new metrics
    def update_bpf_block_irq(self, deployment, bpf_block_irq):
        self.logger.info(
            f"Updating BPF block IRQ for {deployment}: {bpf_block_irq}")
        self.bpf_block_irq_gauge.labels(
            deployment=deployment).set(bpf_block_irq)

    def update_bpf_cpu_time_ms(self, deployment, bpf_cpu_time_ms):
        self.logger.info(
            f"Updating BPF CPU time in ms for {deployment}: {bpf_cpu_time_ms}")
        self.bpf_cpu_time_ms_gauge.labels(
            deployment=deployment).set(bpf_cpu_time_ms)

    def update_bpf_net_rx_irq(self, deployment, bpf_net_rx_irq):
        self.logger.info(
            f"Updating BPF net RX IRQ for {deployment}: {bpf_net_rx_irq}")
        self.bpf_net_rx_irq_gauge.labels(
            deployment=deployment).set(bpf_net_rx_irq)

    def update_bpf_net_tx_irq(self, deployment, bpf_net_tx_irq):
        self.logger.info(
            f"Updating BPF net TX IRQ for {deployment}: {bpf_net_tx_irq}")
        self.bpf_net_tx_irq_gauge.labels(
            deployment=deployment).set(bpf_net_tx_irq)

    def update_bpf_page_cache_hit(self, deployment, bpf_page_cache_hit):
        self.logger.info(
            f"Updating BPF page cache hit for {deployment}: {bpf_page_cache_hit}"
        )
        self.bpf_page_cache_hit_gauge.labels(
            deployment=deployment).set(bpf_page_cache_hit)

    def update_joules(self, deployment, joules):
        self.logger.info(f"Updating joules for {deployment}: {joules}")
        self.joules_gauge.labels(deployment=deployment).set(joules)

    def update_metrics(self):
        self.logger.info("Updating metrics")
        for metric in self.metrics:
            try:
                metric_data = self.prom.get_current_metric_value(
                    metric_name=metric)
                if not metric_data:
                    self.logger.warning(f"No data found for metric: {metric}")
                    continue
                self.logger.info(f"Metric: {metric}")
                for data in metric_data:
                    deployment_name = data['metric'].get('pod_name')
                    value = float(data['value'][1])
                    # Example to update a specific gauge, adapt as needed
                    if metric == 'kepler_container_gpu_joules_total':
                        self.update_gpu_joules(deployment_name, value)
                    elif metric == 'kepler_container_dram_joules_total':
                        self.update_dram_joules(deployment_name, value)
                    elif metric == 'kepler_container_core_joules_total':
                        self.update_core_joules(deployment_name, value)
                    elif metric == 'kepler_container_bpf_block_irq_total':
                        self.update_bpf_block_irq(deployment_name, value)
                    elif metric == 'kepler_container_bpf_cpu_time_ms_total':
                        self.update_bpf_cpu_time_ms(deployment_name, value)
                    elif metric == 'kepler_container_bpf_net_rx_irq_total':
                        self.update_bpf_net_rx_irq(deployment_name, value)
                    elif metric == 'kepler_container_bpf_net_tx_irq_total':
                        self.update_bpf_net_tx_irq(deployment_name, value)
                    elif metric == 'kepler_container_bpf_page_cache_hit_total':
                        self.update_bpf_page_cache_hit(deployment_name, value)
                    elif metric == 'kepler_container_joules_total':
                        self.update_joules(deployment_name, value)
                    # Add similar updates for other metrics
            except requests.exceptions.RequestException as e:
                self.logger.error(f"Request failed for metric {metric}: {e}")
            except ValueError as e:
                print(f"Failed to decode JSON for metric {metric}: {e}")


# Example usage
if __name__ == "__main__":
    pcm = PowerCappingMetrics()
    pcm.start_server()
