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


python main.py --mode clear-key
python main.py --mode encrypt --message "Hello Quantum World" --key-length 256
python main.py --mode decrypt --message "Ubc95VQoTasmXYqI3JK5Cv6FTA==" --key-length 256
python main.py --mode show-circuits
python main.py --mode show-circuits --viz-qubits 6
