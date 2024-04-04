# placeholder to simulate different timer intervals, replica incremental, power usage modeling, and high/low watermarks
import random
import time

# Simulation parameters
NUM_REPLICAS = 10
POWER_CAP = 1000
CONSUMPTION_MEAN = 100
CONSUMPTION_STD_DEV = 20
OPERATOR_ADJUSTMENT_PROBABILITY = 0.8
NUM_SIMULATIONS = 1000


def simulate_power_consumption(num_replicas):
    power_consumption = []
    for _ in range(num_replicas):
        consumption = random.gauss(CONSUMPTION_MEAN, CONSUMPTION_STD_DEV)
        power_consumption.append(consumption)
    return power_consumption


def adjust_replicas(power_consumption, power_cap):
    total_consumption = sum(power_consumption)
    if total_consumption > power_cap:
        if random.random() < OPERATOR_ADJUSTMENT_PROBABILITY:
            num_replicas = len(power_consumption)
            adjusted_replicas = int(num_replicas *
                                    (power_cap / total_consumption))
            return adjusted_replicas
    return len(power_consumption)


def run_simulation(timer_interval):
    success_count = 0
    for _ in range(NUM_SIMULATIONS):
        power_consumption = simulate_power_consumption(NUM_REPLICAS)
        adjusted_replicas = adjust_replicas(power_consumption, POWER_CAP)
        time.sleep(timer_interval)
        if adjusted_replicas <= NUM_REPLICAS:
            success_count += 1
    success_rate = success_count / NUM_SIMULATIONS
    return success_rate


# Simulate with different timer intervals
timer_intervals = [0.1, 0.5, 1.0, 2.0, 5.0]
for interval in timer_intervals:
    success_rate = run_simulation(interval)
    print(
        f"Timer Interval: {interval} seconds, Success Rate: {success_rate:.2f}"
    )
