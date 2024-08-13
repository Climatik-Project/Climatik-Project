import logging
import os
import sys

# Import the required modules
src_path = os.path.join(os.path.dirname(__file__), '.')
sys.path.append(src_path)

import kopf
import kubernetes
from jsonschema import validate, ValidationError
from crd import POWER_CAPPING_CONFIG_SCHEMA
from strategies import get_power_capping_strategy
from prom_metrics import PowerCappingMetrics
from prometheus_api_client import PrometheusConnect

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[logging.StreamHandler(sys.stdout)])

# Set the logging level based on environment variable
log_level = os.getenv('LOG_LEVEL', 'INFO').upper()
logging.getLogger().setLevel(log_level)

prom_host = os.getenv('PROMETHEUS_HOST', 'http://localhost:9090')
high_power_usage_ratio = float(os.getenv('HIGH_POWER_USAGE_RATIO', '0.95'))
moderate_power_usage_ratio = float(
    os.getenv('MODERATE_POWER_USAGE_RATIO', '0.8'))

logging.info(f"Prometheus host: {prom_host}")
logging.info(f"High power usage ratio: {high_power_usage_ratio}")
logging.info(f"Moderate power usage ratio: {moderate_power_usage_ratio}")

# Obtain the selected power capping strategy from an environment variable
selected_strategy = os.getenv('POWER_CAPPING_STRATEGY', 'maximize_replicas')
logging.info(f"Selected power capping strategy: {selected_strategy}")

# Get the selected power capping strategy instance
power_capping_strategy = get_power_capping_strategy(selected_strategy)

# Create a Prometheus API client
prom = PrometheusConnect(url=prom_host, disable_ssl=True)
metrics = PowerCappingMetrics(url=prom_host)


@kopf.on.startup()
def start_metrics_server(**kwargs):
    logging.info("Starting metrics server")
    metrics.start_server()


@kopf.on.create('climatik-project.io', 'v1alpha1', 'powercappingconfigs')
def create_power_capping_config(spec, **kwargs):
    logging.info("Creating power capping config")
    try:
        # Validate the CRD spec against the JSON schema
        validate(spec, POWER_CAPPING_CONFIG_SCHEMA)
    except ValidationError as e:
        logging.error(f"Invalid PowerCappingConfig spec: {e.message}")
        raise kopf.PermanentError(
            f"Invalid PowerCappingConfig spec: {e.message}")

    # Retrieve the power capping configuration from the custom resource
    power_cap_limit = spec.get('powerCapLimit')
    scale_object_refs = spec.get('scaledObjectRefs', [])
    logging.info(f"Power cap limit: {power_cap_limit}")
    logging.info(f"Scale object references: {scale_object_refs}")

    # Iterate over the ScaleObjectRefs and update the KEDA ScaleObjects
    for scale_object_ref in scale_object_refs:
        api_version = scale_object_ref.get('apiVersion')
        kind = scale_object_ref.get('kind')
        name = scale_object_ref.get('metadata', {}).get('name')
        logging.info(f"Processing scale object: {name}")

        # Retrieve the KEDA ScaleObject
        api_instance = kubernetes.client.CustomObjectsApi()
        scale_object = api_instance.get_namespaced_custom_object(
            group=api_version.split('/')[0],
            version=api_version.split('/')[1],
            namespace=kwargs['namespace'],
            plural=f"{kind.lower()}s",
            name=name)

        # Update the ScaleObject with the power capping configuration
        max_replicas = calculate_max_replicas(power_cap_limit)
        logging.info(f"Max replicas for {name}: {max_replicas}")
        scale_object['spec']['maxReplicaCount'] = max_replicas

        # Update the ScaleObject in the Kubernetes cluster
        api_instance.patch_namespaced_custom_object(
            group=api_version.split('/')[0],
            version=api_version.split('/')[1],
            namespace=kwargs['namespace'],
            plural=f"{kind.lower()}s",
            name=name,
            body=scale_object)


@kopf.timer('climatik-project.io',
            'v1alpha1',
            'powercappingconfigs',
            interval=10.0)
def monitor_power_usage(spec, status, patch, **kwargs):
    logging.info("Monitoring power usage")
    power_cap_limit = spec.get('powerCapLimit')
    scale_object_refs = spec.get('scaledObjectRefs', [])

    current_replicas = {}
    power_consumptions = {}

    crd_api_instance = kubernetes.client.CustomObjectsApi()
    api_instance = kubernetes.client.AppsV1Api()

    for scale_object_ref in scale_object_refs:
        api_version = scale_object_ref['apiVersion']
        kind = scale_object_ref['kind']
        name = scale_object_ref['metadata']['name']
        logging.info(f"Retrieving scale object: {name}")

        # Retrieve the KEDA ScaledObject
        scaled_object = crd_api_instance.get_namespaced_custom_object(
            group=api_version.split('/')[0],
            version=api_version.split('/')[1],
            namespace=kwargs['namespace'],
            plural=f"{kind.lower()}s",
            name=name)

        deployment_name = scaled_object['spec']['scaleTargetRef']['name']

        # Retrieve the current number of replicas from the deployment
        deployment = api_instance.read_namespaced_deployment(
            namespace=kwargs['namespace'], name=deployment_name)
        current_replicas[deployment_name] = deployment.spec.replicas

        # Retrieve the power consumption for the deployment
        power_consumption = get_power_consumption(deployment_name,
                                                  kwargs['namespace'])
        power_consumptions[deployment_name] = power_consumption

    logging.info(f"Power consumption: {power_consumptions}")
    logging.info(f"Current replicas: {current_replicas}")

    # Calculate the updated maxReplicas for each deployment based on the selected strategy
    total_power_consumption = sum(power_consumptions.values())
    updated_max_replicas = {}

    if total_power_consumption > 0:
        updated_max_replicas = power_capping_strategy.calculate_max_replicas(
            current_replicas, power_consumptions, power_cap_limit)

    # Update the maxReplicaCount for each ScaledObject
    for scale_object_ref in scale_object_refs:
        api_version = scale_object_ref['apiVersion']
        kind = scale_object_ref['kind']
        name = scale_object_ref['metadata']['name']

        # Retrieve the KEDA ScaledObject
        scaled_object = crd_api_instance.get_namespaced_custom_object(
            group=api_version.split('/')[0],
            version=api_version.split('/')[1],
            namespace=kwargs['namespace'],
            plural=f"{kind.lower()}s",
            name=name)

        deployment_name = scaled_object['spec']['scaleTargetRef']['name']
        logging.info(
            f"ScaleTargetRef: {scaled_object['spec']['scaleTargetRef']}")
        max_replicas = updated_max_replicas.get(
            deployment_name, 1)  # Default to 1 if not calculated

        # Update the maxReplicaCount in the ScaledObject
        scaled_object['spec']['maxReplicaCount'] = max_replicas

        # Update the ScaledObject in the Kubernetes cluster
        crd_api_instance.patch_namespaced_custom_object(
            group=api_version.split('/')[0],
            version=api_version.split('/')[1],
            namespace=kwargs['namespace'],
            plural=f"{kind.lower()}s",
            name=name,
            body=scaled_object)

    # Update Prometheus metrics
    metrics.update_scale_objects(len(scale_object_refs))

    forecast_power_consumption = {}

    for deployment_name, power_consumption in power_consumptions.items():
        metrics.update_replicas(deployment_name,
                                current_replicas[deployment_name])
        metrics.update_power_consumption(deployment_name, power_consumption)
        forecast_power_consumption[
            deployment_name] = power_consumption * updated_max_replicas.get(
                deployment_name, 1)  # Default to 1 if not calculated

    for deployment_name, forecast_power in forecast_power_consumption.items():
        metrics.update_forecast_power_consumption(deployment_name,
                                                  forecast_power)

    # Update the status with the current and forecast power consumption
    patch['status'] = {
        'currentPowerConsumption': sum(power_consumptions.values()),
        'forecastPowerConsumption': sum(forecast_power_consumption.values())
    }


def calculate_max_replicas(power_cap_limit):
    # Implement the logic to calculate the maximum replicas based on the power cap limit
    # This is just a placeholder, replace it with your actual calculation
    logging.info(
        f"Calculating max replicas for power cap limit: {power_cap_limit}")
    return int(power_cap_limit / 100)


def get_current_replica_from_scale_object(api_instance, namespace,
                                          scale_object):
    deployment = api_instance.read_namespaced_deployment(
        namespace=namespace,
        name=scale_object['spec']['scaleTargetRef']['name'])
    return deployment.spec.replicas


def get_power_consumption(deployment_name, namespace):
    # get kepler container joules total metric
    query = f'sum(rate(kepler_container_joules_total{{container_namespace="{namespace}", container_name="{deployment_name}"}}[5m]))'

    logging.info(f"Power consumption query: {query}")
    # Execute the Prometheus query
    result = prom.custom_query(query=query)
    logging.info(
        f"Power consumption query result for {deployment_name}: {result}")
    # Extract the power consumption value from the query result
    power_consumption = 0
    if result:
        power_consumption = float(result[0]['value'][1])

    return power_consumption
