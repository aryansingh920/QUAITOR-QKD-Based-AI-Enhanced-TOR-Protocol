package client

import (
	"fmt"
	"net"
	"tor-protocol/config"
	"tor-protocol/relay"
)

func StartClient() {
    conn, err := net.Dial("tcp", config.Relay1Address)
    if err != nil {
        fmt.Println("Error connecting to relay:", err)
        return
    }
    defer conn.Close()

    message := "Hello, Server!"
    encrypted, _ := relay.Encrypt([]byte(message))
    conn.Write(encrypted)
    fmt.Println("Client sent message:", message)
}
