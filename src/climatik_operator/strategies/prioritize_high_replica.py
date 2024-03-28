from .base_strategy import BasePowerCappingStrategy
from typing import Dict


# MaximizeReplicasStrategy strategy is to prioritize those deployments that have bursty requests to better serve the requests without delay
class PrioritizeHighReplicaStrategy(BasePowerCappingStrategy):

    def calculate_max_replicas(self, current_replicas: Dict[str, int],
                               power_consumptions: Dict[str, float],
                               total_power_cap: float) -> Dict[str, int]:
        updated_max_replicas = {}
        # TODO implement this strategy
        # Placeholder
        for deployment_name, power_consumption in power_consumptions.items():
            # calculate the max replicas based on the power cap limit
            updated_max_replicas[deployment_name] = int(total_power_cap /
                                                        power_consumption)
        return updated_max_replicas
