package main

import (
	"tor-protocol/client"
	"tor-protocol/relay"
	"tor-protocol/server"
)

func main() {
    go server.StartServer()
    go relay.StartRelay("Relay1")
    go relay.StartRelay("Relay2")
    client.StartClient()
}
