from flask import Flask, render_template, request
from prometheus_client import Histogram, make_wsgi_app
from werkzeug.middleware.dispatcher import DispatcherMiddleware
from werkzeug.serving import run_simple

# Define the Flask app
app = Flask(__name__)

# Define the Prometheus metric
histogram = Histogram('http_requests_latency_seconds', 'HTTP request latency',labelnames=['user', 'code'],
                      buckets=[0.1, 0.2, 0.5, 1.0, 2.0, 5.0, 10.0])

# Setup a middleware to expose the Prometheus metrics
application = DispatcherMiddleware(app, {'/metrics': make_wsgi_app()})

# Define the index page with two columns and five rows
@app.route('/', methods=['GET', 'POST'])
def index():
    if request.method == 'POST':
        # Get the username and button information from the form
        username = request.form['username']
        button_type = request.form['button']

        # Simulate the request processing
        if button_type == 'fast':
            elapsed_time = 0.1
            statusCode = '200'
        elif button_type == 'slow':
            elapsed_time = 5.0
            statusCode = '200'
        else:
            elapsed_time = 10.0
            statusCode = '500'

        # Record the request latency and success/failure to the metric
        histogram.labels(user=username, code=statusCode).observe(elapsed_time, {'user': username, 'code': statusCode})

    return render_template('index.html')

if __name__ == '__main__':
    # app.run(host='0.0.0.0', port=8080)
    run_simple(hostname='0.0.0.0', port=8080, application=application, threaded=True)
