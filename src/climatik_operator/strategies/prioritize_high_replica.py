from .base_strategy import BasePowerCappingStrategy
from typing import Dict


# MaximizeReplicasStrategy strategy is to prioritize those deployments that have bursty requests to better serve the requests without delay
class PrioritizeHighReplicaStrategy(BasePowerCappingStrategy):

    def calculate_max_replicas(self, current_replicas: Dict[str, int],
                               power_consumptions: Dict[str, float],
                               total_power_cap: float) -> Dict[str, int]:
        updated_max_replicas = {}
        # TODO implement this strategy
        return updated_max_replicas
