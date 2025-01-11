#!/bin/zsh

clear
# Default number of nodes to start
DEFAULT_NODES=5

# Check if `n` is passed as an argument; if not, use default
NODES_COUNT=${1:-$DEFAULT_NODES}

# Range of ports to search for availability
START_PORT=9000
END_PORT=9999

# Function to check if a port is available
is_port_available() {
  ! nc -z 127.0.0.1 "$1" &> /dev/null
}

# Find `n` available ports
available_ports=()
for ((port=$START_PORT; port<=$END_PORT; port++)); do
  if is_port_available "$port"; then
    available_ports+=($port)
    # Stop searching if we've found enough
    if [ "${#available_ports[@]}" -ge "$NODES_COUNT" ]; then
      break
    fi
  fi
done

# Verify we have enough available ports
if [ "${#available_ports[@]}" -lt "$NODES_COUNT" ]; then
  echo "Error: Not enough available ports found in the range $START_PORT-$END_PORT."
  exit 1
fi

echo "Selected ports: ${available_ports[*]}"

# Randomly choose one port to act as the client
client_port=${available_ports[$((RANDOM % NODES_COUNT))]}

# Start nodes
for port in "${available_ports[@]}"; do
  if [ "$port" -eq "$client_port" ]; then
    echo "Starting client on port $port..."
    go run main.go "$port" client &
  else
    echo "Starting node on port $port..."
    go run main.go "$port" &
  fi
done

# Wait for processes to complete (or you can exit early if desired)
echo "All nodes and client started. Press Ctrl+C to terminate."
wait
