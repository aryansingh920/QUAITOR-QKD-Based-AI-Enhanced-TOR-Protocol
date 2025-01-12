// main.go
package main

import (
	"log"
	"math/rand"
	"time"

	"tor-protocol/server"
)

func main() {
    // Seed random
    rand.Seed(time.Now().UnixNano())

    log.Printf("tor-protocol API\n")

    // Start server
    server.ServerMain()
}
