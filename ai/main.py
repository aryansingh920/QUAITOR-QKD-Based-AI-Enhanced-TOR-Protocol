import numpy as np
from sklearn.ensemble import RandomForestRegressor
import requests

# Sample training data
X = np.array([[1, 10], [2, 20], [3, 30]])  # Features: node ID and load
y = np.array([0.5, 1.2, 0.8])  # Delays

# Train model
model = RandomForestRegressor()
model.fit(X, y)


def predict_delay(node_id, load):
    # Predict delay based on node ID and load
    return model.predict([[node_id, load]])[0]

# Update Go middleware with traffic data


def update_traffic_data():
    traffic_data = {}
    for node_id in range(1, 11):  # Example: nodes 1–10
        load = np.random.randint(10, 100)  # Mock load
        delay = predict_delay(node_id, load)
        traffic_data[node_id] = delay
    requests.post("http://middleware/traffic",
                  json=traffic_data)  # Send to middleware
