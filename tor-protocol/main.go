// main.go
package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"tor-protocol/config"
	"tor-protocol/server"
)

func main() {
    // Seed random
    rand.Seed(time.Now().UnixNano())

    log.Printf("tor-protocol API\n")

    // Start server
    fmt.Println("At ServerMain: LoadConfig")
	config.LoadConfig()
    server.ServerMain()
}
