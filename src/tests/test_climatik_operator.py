import unittest
from unittest.mock import MagicMock, patch
from climatik_operator.operator import create_power_capping_config, monitor_power_usage, calculate_max_replicas
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
            'scaleObjectRefs': [{
                'apiVersion': 'keda.sh/v1alpha1',
                'kind': 'ScaleObject',
                'metadata': {
                    'name': 'test-scaleobject'
                }
            }]
        }

        # Call the create_power_capping_config function
        create_power_capping_config(spec, namespace='default')

        # Assert that the necessary API calls were made
        self.mock_api_client.get_namespaced_custom_object.assert_called_once()
        self.mock_api_client.patch_namespaced_custom_object.assert_called_once(
        )

    @patch('climatik_operator.operator.prom', new_callable=MagicMock)
    @patch('climatik_operator.operator.kubernetes.client.CustomObjectsApi')
    def test_monitor_power_usage(self, mock_custom_objects_api, mock_prom):
        # Mock the CustomObjectsApi and PrometheusConnect
        mock_custom_objects_api.return_value = self.mock_api_client
        mock_prom.custom_query.return_value = [{'value': [None, 800]}]

        # Define the test spec
        spec = {
            'powerCapLimit':
            1000,
            'scaleObjectRefs': [{
                'apiVersion': 'keda.sh/v1alpha1',
                'kind': 'ScaleObject',
                'metadata': {
                    'name': 'test-scaleobject'
                }
            }]
        }

        # Call the monitor_power_usage function
        monitor_power_usage(spec, namespace='default')

        # Assert that the necessary API calls were made
        mock_prom.custom_query.assert_called_once()
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
