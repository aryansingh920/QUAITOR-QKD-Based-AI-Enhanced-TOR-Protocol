#!/bin/zsh

# cd tor-protocol && go mod init tor-protocol

go mod tidy

go build ./...

go list -m

go run main.go

go clean -modcache

go test ./...

go test ./... -bench=. -benchmem -v

air init
air


python qkd/main.py --mode encrypt --message "Hello Quantum World" --key-length 256
python qkd/main.py --mode clear-key
python qkd/main.py --mode decrypt --message "Ubc95VQoTasmXYqI3JK5Cv6FTA==" --key-length 256
python qkd/main.py --mode show-circuits
python qkd/main.py --mode show-circuits --viz-qubits 6

python ai/main.py --start_node 1 --end_node 10 --middleware_url http://middleware/traffic
