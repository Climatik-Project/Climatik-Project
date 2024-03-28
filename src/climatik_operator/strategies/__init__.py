from .maximize_replicas import MaximizeReplicasStrategy


def get_power_capping_strategy(strategy_name: str) -> BasePowerCappingStrategy:
    strategies = {'maximize_replicas': MaximizeReplicasStrategy}
    strategy_class = strategies.get(strategy_name)
    if strategy_class:
        return strategy_class()
    else:
        raise ValueError(f"Invalid power capping strategy: {strategy_name}")
