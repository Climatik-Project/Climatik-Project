import os
import sys

#################################################################
# import internal src
src_path = os.path.join(os.path.dirname(__file__), '..')
sys.path.append(src_path)

import unittest
from unittest.mock import MagicMock, patch
from climatik_operator.operator import create_power_capping_config, calculate_max_replicas
from kubernetes import client
from prometheus_api_client import PrometheusConnect


class TestPowerCappingOperator(unittest.TestCase):

    def setUp(self):
        # Create a mock Kubernetes API client
        self.mock_api_client = MagicMock()

        # Create a mock Prometheus API client
        self.mock_prom_client = MagicMock(spec=PrometheusConnect)

    @patch('climatik_operator.operator.kubernetes.client.CustomObjectsApi')
    def test_create_power_capping_config(self, mock_custom_objects_api):
        # Mock the CustomObjectsApi
        mock_custom_objects_api.return_value = self.mock_api_client

        # Define the test spec
        spec = {
            'powerCapLimit':
            1000,
            'scaledObjectRefs': [{
                'apiVersion': 'keda.sh/v1alpha1',
                'kind': 'ScaledObject',
                'metadata': {
                    'name': 'test-scaledobject'
                }
            }]
        }

        # Define the test ScaledObject
        scaled_object = {
            'spec': {
                'scaleTargetRef': {
                    'name': 'test-deployment'
                },
                'maxReplicaCount': 5
            }
        }

        # Define the test deployment
        deployment = MagicMock(spec=client.V1Deployment)
        deployment.spec.replicas = 3

        # Mock the get_namespaced_custom_object and read_namespaced_deployment methods
        self.mock_api_client.get_namespaced_custom_object.return_value = scaled_object
        self.mock_api_client.read_namespaced_deployment.return_value = deployment

        # Call the create_power_capping_config function
        create_power_capping_config(spec, namespace='default')

        # Assert that the necessary API calls were made
        self.mock_api_client.get_namespaced_custom_object.assert_called_once()
        self.mock_api_client.patch_namespaced_custom_object.assert_called_once(
        )

    def test_calculate_max_replicas(self):
        # Define the test power cap limit
        power_cap_limit = 1000

        # Call the calculate_max_replicas function
        max_replicas = calculate_max_replicas(power_cap_limit)

        # Assert the expected result
        self.assertEqual(max_replicas, 10)


if __name__ == '__main__':
    unittest.main()
