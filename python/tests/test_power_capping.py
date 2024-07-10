import unittest
from unittest.mock import MagicMock, patch
from climatik_operator.operator import monitor_power_usage
from kubernetes import client
from prometheus_api_client import PrometheusConnect


class TestPowerCappingStrategies(unittest.TestCase):

    def setUp(self):
        # Create a mock Kubernetes API client
        self.mock_api_client = MagicMock()

        # Create a mock Prometheus API client
        self.mock_prom_client = MagicMock(spec=PrometheusConnect)

    @patch('climatik_operator.operator.prom', new_callable=MagicMock)
    @patch('climatik_operator.operator.kubernetes.client.CustomObjectsApi')
    @patch('climatik_operator.operator.kubernetes.client.AppsV1Api')
    def test_monitor_power_usage_high_consumption(self, mock_apps_api,
                                                  mock_custom_objects_api,
                                                  mock_prom):
        # Mock the CustomObjectsApi and PrometheusConnect
        mock_custom_objects_api.return_value = self.mock_api_client
        mock_apps_api.return_value = self.mock_api_client
        mock_prom.custom_query.return_value = [{
            'value': [None, 950]
        }]  # High power consumption

        # Define the test spec and status
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
        status = {}

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

        # Call the monitor_power_usage function
        monitor_power_usage(spec, status, namespace='default')

        # Assert that the necessary API calls were made
        mock_prom.custom_query.assert_called_once()
        self.assertEqual(
            self.mock_api_client.get_namespaced_custom_object.call_count, 2)
        self.mock_api_client.read_namespaced_deployment.assert_called_once()
        self.mock_api_client.patch_namespaced_custom_object.assert_called_once(
        )

        # Assert the status updates
        self.assertEqual(status['currentPowerConsumption'], 950)
        self.assertEqual(status['forecastPowerConsumption'], 950)

        # Assert that the maxReplicaCount is set to the current number of replicas
        patch_call_args = self.mock_api_client.patch_namespaced_custom_object.call_args[
            1]
        self.assertEqual(patch_call_args['body']['spec']['maxReplicaCount'], 1)

    @patch('climatik_operator.operator.prom', new_callable=MagicMock)
    @patch('climatik_operator.operator.kubernetes.client.CustomObjectsApi')
    @patch('climatik_operator.operator.kubernetes.client.AppsV1Api')
    def test_monitor_power_usage_moderate_consumption(self, mock_apps_api,
                                                      mock_custom_objects_api,
                                                      mock_prom):
        # Mock the CustomObjectsApi and PrometheusConnect
        mock_custom_objects_api.return_value = self.mock_api_client
        mock_apps_api.return_value = self.mock_api_client
        mock_prom.custom_query.return_value = [{
            'value': [None, 800]
        }]  # Moderate power consumption

        # Define the test spec and status
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
        status = {}

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

        # Call the monitor_power_usage function
        monitor_power_usage(spec, status, namespace='default')

        # Assert that the necessary API calls were made
        mock_prom.custom_query.assert_called_once()
        self.assertEqual(
            self.mock_api_client.get_namespaced_custom_object.call_count, 2)
        self.mock_api_client.read_namespaced_deployment.assert_called_once()
        self.mock_api_client.patch_namespaced_custom_object.assert_called_once(
        )

        # Assert the status updates
        self.assertEqual(status['currentPowerConsumption'], 800)
        self.assertEqual(int(status['forecastPowerConsumption']), int(800))

        # Assert that the maxReplicaCount is set to one above the current number of replicas
        patch_call_args = self.mock_api_client.patch_namespaced_custom_object.call_args[
            1]
        self.assertEqual(patch_call_args['body']['spec']['maxReplicaCount'], 1)

    @patch('climatik_operator.operator.prom', new_callable=MagicMock)
    @patch('climatik_operator.operator.kubernetes.client.CustomObjectsApi')
    @patch('climatik_operator.operator.kubernetes.client.AppsV1Api')
    def test_monitor_power_usage_low_consumption(self, mock_apps_api,
                                                 mock_custom_objects_api,
                                                 mock_prom):
        # Mock the CustomObjectsApi and PrometheusConnect
        mock_custom_objects_api.return_value = self.mock_api_client
        mock_apps_api.return_value = self.mock_api_client
        mock_prom.custom_query.return_value = [{
            'value': [None, 500]
        }]  # Low power consumption

        # Define the test spec and status
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
        status = {}

        # Call the monitor_power_usage function
        monitor_power_usage(spec, status, namespace='default')

        # Assert that the necessary API calls were made
        mock_prom.custom_query.assert_called_once()
        self.assertEqual(
            self.mock_api_client.get_namespaced_custom_object.call_count, 2)
        self.mock_api_client.read_namespaced_deployment.assert_called_once()
        self.mock_api_client.patch_namespaced_custom_object.assert_called_once(
        )

        # Assert the status updates
        self.assertEqual(status['currentPowerConsumption'], 500)
        self.assertIn('forecastPowerConsumption', status)


if __name__ == '__main__':
    unittest.main()
