/*
Created on 11/01/2025

@author: Aryan

Filename: node.go

Relative Path: tor-protocol/protocol/node.go
*/

package protocol

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

// NodeRole represents the role of a Node: entry, middle (relay), or exit.
type NodeRole string

const (
    EntryNode NodeRole = "ENTRY"
    RelayNode NodeRole = "RELAY"
    ExitNode  NodeRole = "EXIT"
)

// Node represents a Tor-like node
type Node struct {
    ID         string
    Role       NodeRole
    Address    string
    Port       int
    listener   net.Listener
    stopCh     chan struct{}
    wg         sync.WaitGroup

    // For demonstration, we'll generate ephemeral keys per node
    ephemeralPublicKey  []byte
    ephemeralPrivateKey []byte
}

// Start begins listening on the specified port
func (n *Node) Start() error {
    var err error
    addr := fmt.Sprintf(":%d", n.Port)

    n.listener, err = net.Listen("tcp", addr)
    if err != nil {
        return fmt.Errorf("Node %s failed to listen on port %d: %v", n.ID, n.Port, err)
    }

    // Generate ephemeral key pair for this node
    n.ephemeralPublicKey, n.ephemeralPrivateKey, err = GenerateEphemeralKeyPair()
    if err != nil {
        log.Printf("[%s] Failed to generate ephemeral keys: %v\n", n.ID, err)
    } else {
        log.Printf("[%s] Ephemeral keys generated successfully\n", n.ID)
    }

    log.Printf("[%s] %s node listening on %s\n", n.ID, n.Role, addr)

    // accept connections
    n.wg.Add(1)
    go func() {
        defer n.wg.Done()
        for {
            conn, err := n.listener.Accept()
            if err != nil {
                select {
                case <-n.stopCh:
                    // we're shutting down
                    return
                default:
                    log.Printf("[%s] Accept error: %v", n.ID, err)
                    continue
                }
            }
            // Handle connection in a separate goroutine
            n.wg.Add(1)
            go func(c net.Conn) {
                defer n.wg.Done()
                n.handleConnection(c)
            }(conn)
        }
    }()
    return nil
}

// Stop stops listening and closes the node
func (n *Node) Stop() {
    close(n.stopCh)
    if n.listener != nil {
        _ = n.listener.Close()
    }
    n.wg.Wait()
}

// handleConnection handles incoming data from a Node or from a Client
func (n *Node) handleConnection(conn net.Conn) {
    defer conn.Close()

    buffer := make([]byte, 4096)
    for {
        // read data
        conn.SetDeadline(time.Now().Add(30 * time.Second))
        bytesRead, err := conn.Read(buffer)
        if err != nil {
            return
        }

        data := buffer[:bytesRead]

        // Parse the RelayCell from the raw data
        relayCell, err := ParseRelayCell(data)
        if err != nil {
            log.Printf("[%s] Error parsing relay cell: %v\n", n.ID, err)
            return
        }

        // In a real Tor, we'd decrypt layers here using ephemeral keys.
        // For demonstration, skip actual decrypt but show placeholder.
        // e.g., decryptedPayload := DecryptPayload(relayCell.Payload, n.ephemeralPrivateKey)

        // If this is the exit node, we "finalize" the message (like sending to the destination)
        // For demonstration, we just log it
        if n.Role == ExitNode {
            log.Printf("[%s] EXIT node final output: %s\n", n.ID, relayCell.Payload)
            // In a real Tor network, you'd forward to the actual destination here.
        } else {
            // If not exit, forward to the next node
            // The relayCell.NextAddr is the IP:port of the next node
            if err := forwardToNextNode(relayCell); err != nil {
                log.Printf("[%s] Error forwarding to next node: %v\n", n.ID, err)
            }
        }
    }
}

// GenerateRandomTraffic simulates dummy traffic to a random node at intervals.
// This is purely for demonstration to mimic "cover traffic" or "dummy" traffic.
func (n *Node) GenerateRandomTraffic(intervalSeconds int) {
    ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            // pick a random message length and random node
            randomMessage := generateRandomString(rand.Intn(20) + 5) // 5-25 chars
            randomPort := 9001 + rand.Intn(10)                       // from default range
            relayCell := &RelayCell{
                NextAddr: fmt.Sprintf("127.0.0.1:%d", randomPort),
                Payload:  []byte("[DUMMY_TRAFFIC] " + randomMessage),
            }
            // Attempt to dial and send
            err := forwardToNextNode(relayCell)
            if err != nil {
                log.Printf("[%s] Random traffic error: %v\n", n.ID, err)
            } else {
                log.Printf("[%s] Sent dummy traffic to %d\n", n.ID, randomPort)
            }
        case <-n.stopCh:
            return
        }
    }
}

// generateRandomString creates a pseudo-random string of length n.
func generateRandomString(n int) string {
    letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}
