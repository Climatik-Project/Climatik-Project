from abc import ABC, abstractmethod
from typing import Dict


class BasePowerCappingStrategy(ABC):

    @abstractmethod
    def calculate_max_replicas(self, current_replicas: Dict[str, int],
                               power_consumptions: Dict[str, float],
                               total_power_cap: float) -> Dict[str, int]:
        pass
