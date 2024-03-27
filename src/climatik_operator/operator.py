import kopf
import kubernetes
import os
from jsonschema import validate, ValidationError
from .crd import POWER_CAPPING_CONFIG_SCHEMA

# Import the Prometheus API client
from prometheus_api_client import PrometheusConnect

# obtain prometheus host and power usage ratios from environment variables
prom_host = os.getenv('PROMETHEUS_HOST')
high_power_usage_ratio = float(os.getenv('HIGH_POWER_USAGE_RATIO', '0.95'))
moderate_power_usage_ratio = float(
    os.getenv('MODERATE_POWER_USAGE_RATIO', '0.8'))

# obtain prometheus host from environment variable
prom_host = os.getenv('PROMETHEUS_HOST')
if not prom_host:
    raise ValueError("PROMETHEUS_HOST environment variable is not set")

# Create a Prometheus API client
prom = PrometheusConnect(url=prom_host, disable_ssl=True)


@kopf.on.create('powercapping.climatik-project.ai', 'v1',
                'powercappingconfigs')
def create_power_capping_config(spec, **kwargs):
    try:
        # Validate the CRD spec against the JSON schema
        validate(spec, POWER_CAPPING_CONFIG_SCHEMA)
    except ValidationError as e:
        raise kopf.PermanentError(
            f"Invalid PowerCappingConfig spec: {e.message}")

    # Retrieve the power capping configuration from the custom resource
    power_cap_limit = spec.get('powerCapLimit')
    scale_object_refs = spec.get('scaleObjectRefs', [])

    # Iterate over the ScaleObjectRefs and update the KEDA ScaleObjects
    for scale_object_ref in scale_object_refs:
        api_version = scale_object_ref.get('apiVersion')
        kind = scale_object_ref.get('kind')
        name = scale_object_ref.get('metadata', {}).get('name')

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
        scale_object['spec']['maxReplicaCount'] = max_replicas

        # Update the ScaleObject in the Kubernetes cluster
        api_instance.patch_namespaced_custom_object(
            group=api_version.split('/')[0],
            version=api_version.split('/')[1],
            namespace=kwargs['namespace'],
            plural=f"{kind.lower()}s",
            name=name,
            body=scale_object)


@kopf.timer('powercapping.climatik-project.ai',
            'v1',
            'powercappingconfigs',
            interval=60.0)
def monitor_power_usage(spec, status, **kwargs):
    # Retrieve the power capping configuration from the custom resource
    power_cap_limit = spec.get('powerCapLimit')
    scale_object_refs = spec.get('scaleObjectRefs', [])

    # obtain kepler power consumption kepler_node_joules_total and apply irate to get power in watts from prometheus client
    power_consumption = prom.custom_query(
        query="irate(kepler_node_joules_total[1m])")[0]['value'][1]

    # Update the status with the current power consumption
    status['currentPowerConsumption'] = power_consumption

    # Check power usage against the power cap limit
    if power_consumption >= power_cap_limit * high_power_usage_ratio:
        # Power usage is at 95% of the power cap limit
        # Set maxReplicaCount to the current number of replicas
        for scale_object_ref in scale_object_refs:
            api_version = scale_object_ref.get('apiVersion')
            kind = scale_object_ref.get('kind')
            name = scale_object_ref.get('metadata', {}).get('name')

            # Retrieve the KEDA ScaleObject
            api_instance = kubernetes.client.CustomObjectsApi()
            scale_object = api_instance.get_namespaced_custom_object(
                group=api_version.split('/')[0],
                version=api_version.split('/')[1],
                namespace=kwargs['namespace'],
                plural=f"{kind.lower()}s",
                name=name)

            # Set maxReplicaCount to the current number of replicas
            current_replicas = get_current_replica_from_scale_object(
                api_instance=api_instance,
                namespace=kwargs['namespace'],
                scale_object=scale_object)
            scale_object['spec']['maxReplicaCount'] = current_replicas

            # Update the ScaleObject in the Kubernetes cluster
            api_instance.patch_namespaced_custom_object(
                group=api_version.split('/')[0],
                version=api_version.split('/')[1],
                namespace=kwargs['namespace'],
                plural=f"{kind.lower()}s",
                name=name,
                body=scale_object)

        # Update the status with the forecast power consumption
        status['forecastPowerConsumption'] = power_consumption

    elif power_consumption >= power_cap_limit * moderate_power_usage_ratio:
        # Power usage is at 80% of the power cap limit
        # Set maxReplicaCount to one above the current number of replicas
        for scale_object_ref in scale_object_refs:
            api_version = scale_object_ref.get('apiVersion')
            kind = scale_object_ref.get('kind')
            name = scale_object_ref.get('metadata', {}).get('name')

            # Retrieve the KEDA ScaleObject
            api_instance = kubernetes.client.CustomObjectsApi()
            scale_object = api_instance.get_namespaced_custom_object(
                group=api_version.split('/')[0],
                version=api_version.split('/')[1],
                namespace=kwargs['namespace'],
                plural=f"{kind.lower()}s",
                name=name)

            # Set maxReplicaCount to one above the current number of replicas
            current_replicas = get_current_replica_from_scale_object(
                api_instance=api_instance,
                namespace=kwargs['namespace'],
                scale_object=scale_object)
            scale_object['spec']['maxReplicaCount'] = current_replicas + 1

            # Update the ScaleObject in the Kubernetes cluster
            api_instance.patch_namespaced_custom_object(
                group=api_version.split('/')[0],
                version=api_version.split('/')[1],
                namespace=kwargs['namespace'],
                plural=f"{kind.lower()}s",
                name=name,
                body=scale_object)

        # Update the status with the forecast power consumption
        status['forecastPowerConsumption'] = power_consumption * (
            (current_replicas + 1) / current_replicas)


def calculate_max_replicas(power_cap_limit):
    # Implement the logic to calculate the maximum replicas based on the power cap limit
    # This is just a placeholder, replace it with your actual calculation @wenboown
    return int(power_cap_limit / 100)


def get_current_replica_from_scale_object(api_instance, namespace,
                                          scale_object):
    deployment = api_instance.read_namespaced_deployment(
        namespace=namespace,
        name=scale_object['spec']['scaleTargetRef']['name'])

    return deployment.spec.replicas
