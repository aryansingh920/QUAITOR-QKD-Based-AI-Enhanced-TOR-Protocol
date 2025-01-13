#!/bin/zsh


echo "Killing processes running on ports ${ports[@]}..."

for port in {8801..8820}; do
  # Find the process ID (PID) using the port
  pid=$(lsof -t -i tcp:$port)

  if [ -n "$pid" ]; then
    # Kill the process if a PID is found
    kill -9 $pid
    echo "Killed process $pid on port $port."
  else
    echo "No process found running on port $port."
  fi
done

echo "All specified ports have been handled."
