from prometheus_client import Gauge, start_http_server

class PowerCappingMetrics:

    def __init__(self, port=8000):
        self.port = port
        self.scale_objects_gauge = Gauge(
            'power_capping_scaled_objects',
            'Number of ScaledObjects in the PowerCappingConfig CRD')
        self.replicas_gauge = Gauge('power_capping_replicas',
                                    'Number of replicas in each deployment',
                                    ['deployment'])
        self.power_consumption_gauge = Gauge(
            'power_capping_power_consumption',
            'Current power consumption of each deployment', ['deployment'])
        self.forecast_power_consumption_gauge = Gauge(
            'power_capping_forecast_power_consumption',
            'Forecast power consumption of each deployment', ['deployment'])

    def start_server(self):
        start_http_server(self.port)

    def update_scale_objects(self, count):
        self.scale_objects_gauge.set(count)

    def update_replicas(self, deployment, count):
        self.replicas_gauge.labels(deployment=deployment).set(count)

    def update_power_consumption(self, deployment, power_consumption):
        self.power_consumption_gauge.labels(
            deployment=deployment).set(power_consumption)

    def update_forecast_power_consumption(self, deployment,
                                          forecast_power_consumption):
        self.forecast_power_consumption_gauge.labels(
            deployment=deployment).set(forecast_power_consumption)
