package relay

import (
	"fmt"
	"io"
	"net"
	"tor-protocol/config"
)

func StartRelay(name string) {
    listener, err := net.Listen("tcp", getRelayAddress(name))
    if err != nil {
        fmt.Println("Error starting relay:", err)
        return
    }
    defer listener.Close()
    fmt.Println(name, "started on", getRelayAddress(name))

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }
        go handleConnection(conn, name)
    }
}

func handleConnection(conn net.Conn, name string) {
    defer conn.Close()

    data, err := io.ReadAll(conn)
    if err != nil {
        fmt.Println("Error reading data:", err)
        return
    }
    decrypted, _ := Decrypt(data)
    fmt.Println(name, "received message:", string(decrypted))

    if name == "Relay2" {
        forwardToServer(decrypted)
    } else {
        forwardToNextRelay(decrypted)
    }
}

func forwardToNextRelay(data []byte) {
    conn, err := net.Dial("tcp", config.Relay2Address)
    if err != nil {
        fmt.Println("Error connecting to next relay:", err)
        return
    }
    defer conn.Close()
    conn.Write(data)
}

func forwardToServer(data []byte) {
    conn, err := net.Dial("tcp", config.ServerAddress)
    if err != nil {
        fmt.Println("Error connecting to server:", err)
        return
    }
    defer conn.Close()
    conn.Write(data)
}

func getRelayAddress(name string) string {
    if name == "Relay1" {
        return config.Relay1Address
    }
    return config.Relay2Address
}
