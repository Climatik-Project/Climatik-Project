# To test:
# curl -g 'http://127.0.0.1:9090/api/v1/query?query=sum(rate(kepler_container_joules_total[1m]))'

from flask import Flask, jsonify, request
import random
import time

app = Flask(__name__)


# Synthetic data generation
def generate_synthetic_data(start_time, end_time, step):
    data = []
    current_time = start_time
    while current_time <= end_time:
        value = random.uniform(
            0, 10)  # Generate random float value between 0 and 10
        data.append((current_time, value))
        current_time += step
    return data


# Prometheus query API endpoint
@app.route('/api/v1/query', methods=['GET', 'POST'])
def query():
    if request.method == 'GET':
        query = request.args.get('query')
    elif request.method == 'POST':
        data = request.get_json()
        query = data.get('query')
    else:
        return jsonify({
            'status': 'error',
            'errorType': 'bad_request',
            'error': 'Unsupported request method'
        })

    start_time = request.args.get('start',
                                  default=int(time.time()) - 3600,
                                  type=int)
    end_time = request.args.get('end', default=int(time.time()), type=int)
    step = request.args.get('step', default=15, type=int)
    print(
        f'Query: {query}, start: {start_time}, end: {end_time}, step: {step}')
    if query == 'sum(rate(kepler_container_joules_total[1m]))':
        # Generate synthetic data
        data = generate_synthetic_data(start_time, end_time, step)

        # Format the response
        result = {
            'status': 'success',
            'data': {
                'resultType':
                'vector',
                'result': [{
                    'metric': {},
                    'value': [end_time,
                              sum(value for _, value in data)]
                }]
            }
        }

        return jsonify(result)
    else:
        return jsonify({
            'status': 'error',
            'errorType': 'bad_data',
            'error': 'Unsupported query'
        })


if __name__ == '__main__':
    app.run(port=9090)
