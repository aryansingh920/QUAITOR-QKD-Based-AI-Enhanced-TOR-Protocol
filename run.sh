#!/bin/zsh

# cd tor-protocol && go mod init tor-protocol

go mod tidy

go build ./...

go list -m

go run main.go
