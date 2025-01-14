import numpy as np
from sklearn.ensemble import RandomForestRegressor
import requests
import argparse


def train_model():
    """
    Train the RandomForestRegressor model with sample data.
    """
    # Sample training data
    X = np.array([[1, 10], [2, 20], [3, 30]])  # Features: node ID and load
    y = np.array([0.5, 1.2, 0.8])  # Delays

    # Train model
    model = RandomForestRegressor()
    model.fit(X, y)
    return model


def predict_delay(model, node_id, load):
    """
    Predict delay based on node ID and load using the trained model.
    """
    return model.predict([[node_id, load]])[0]


def update_traffic_data(model, start_node, end_node, middleware_url):
    """
    Generate traffic data and send it to the Go middleware.
    """
    traffic_data = {}
    for node_id in range(start_node, end_node + 1):
        load = np.random.randint(10, 100)  # Mock load
        delay = predict_delay(model, node_id, load)
        traffic_data[node_id] = delay

    # Send traffic data to middleware
    response = requests.post(middleware_url, json=traffic_data)
    print(
        f"Response from middleware: {response.status_code} - {response.text}")


if __name__ == "__main__":
    # Set up argument parser
    parser = argparse.ArgumentParser(description="AI Traffic Update Script")
    parser.add_argument("--start_node", type=int,
                        required=True, help="Start node ID")
    parser.add_argument("--end_node", type=int,
                        required=True, help="End node ID")
    parser.add_argument("--middleware_url", type=str, required=True,
                        help="Middleware URL to send traffic data")

    args = parser.parse_args()

    # Train model
    model = train_model()

    # Update traffic data
    update_traffic_data(model, args.start_node,
                        args.end_node, args.middleware_url)
