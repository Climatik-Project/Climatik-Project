from .base_strategy import BasePowerCappingStrategy
from .maximize_replicas import MaximizeReplicasStrategy
from .prioritize_high_replica import PrioritizeHighReplicaStrategy


def get_power_capping_strategy(strategy_name: str) -> BasePowerCappingStrategy:
    strategies = {
        'maximize_replicas': MaximizeReplicasStrategy,
        'prioritize_high_replica': PrioritizeHighReplicaStrategy
    }
    strategy_class = strategies.get(strategy_name)
    if strategy_class:
        return strategy_class()
    else:
        raise ValueError(f"Invalid power capping strategy: {strategy_name}")
