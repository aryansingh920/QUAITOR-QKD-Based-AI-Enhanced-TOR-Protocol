# Simple Makefile to build or run nodes quickly.

.PHONY: all entry relay exit

all: entry relay exit

entry:
	go build -o bin/entry cmd/entry/main.go

relay:
	go build -o bin/relay cmd/relay/main.go

exit:
	go build -o bin/exit cmd/exit/main.go

test:
	go test ./...
