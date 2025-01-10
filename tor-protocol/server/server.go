package server

import (
	"fmt"
	"net"
	"tor-protocol/config"
)

func StartServer() {
    listener, err := net.Listen("tcp", config.ServerAddress)
    if err != nil {
        fmt.Println("Error starting server:", err)
        return
    }
    defer listener.Close()
    fmt.Println("Server started on", config.ServerAddress)

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()

    buf := make([]byte, 1024)
    n, err := conn.Read(buf)
    if err != nil {
        fmt.Println("Error reading from connection:", err)
        return
    }
    fmt.Println("Server received message:", string(buf[:n]))
}
