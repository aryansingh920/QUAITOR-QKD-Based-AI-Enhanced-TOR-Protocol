#!/bin/zsh

# Iterate through the range of ports
for port in {9000..9020}
do
  # Find the process using the port and extract its PID
  pid=$(lsof -ti:$port)

  # Check if a PID exists for the port
  if [ -n "$pid" ]; then
    # Kill the process
    kill -9 $pid
    echo "Killed process $pid on port $port"
  else
    echo "No process running on port $port"
  fi
done
