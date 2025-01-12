#!/bin/zsh

# Check if "main.go" exists
if [ ! -f "main.go" ]; then
  echo "Error: main.go not found in the current directory."
  exit 1
fi

# Start the first tab in the current window
port=9001
osascript <<EOF
tell application "Terminal"
    do script "cd $(pwd) && go run main.go $port"
end tell
EOF

# Open additional tabs for ports 9002 to 9010
for port in {8801..8810}; do
  osascript <<EOF
tell application "Terminal"
    activate
    tell application "System Events" to keystroke "t" using command down
    delay 0.2
    do script "cd $(pwd) && go run main.go $port" in front window
end tell
EOF
done

echo "Launched 10 tabs running 'go run main.go' on ports 9001 to 9010 in the same Terminal window."
