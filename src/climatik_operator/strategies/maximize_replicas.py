from .base_strategy import BasePowerCappingStrategy
from typing import Dict


# MaximizeReplicasStrategy strategy is to maximize the number of replicas of all the deployments
class MaximizeReplicasStrategy(BasePowerCappingStrategy):

    def calculate_max_replicas(self, current_replicas: Dict[str, int],
                               power_consumptions: Dict[str, float],
                               total_power_cap: float) -> Dict[str, int]:
        updated_max_replicas = {}
        # TODO: Implement the logic to maximize the number of replicas based on the power cap limit
        return updated_max_replicas
