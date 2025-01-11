/*
Updated on 11/01/2025
@author: Aryan

Core "protocol" logic, including building circuits and sending data (requests/responses).
*/
package protocol

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

// RelayCell is our minimal container for data + next/prev addresses
type RelayCell struct {
    PrevAddr       string // the node that sent to me
    NextAddr       string // the node to send to next
    Payload        []byte // raw data
    IsExitRequest  bool   // if set, the next node is the final destination
    IsExitResponse bool   // if set, data is returning from final destination
}

// BuildCircuit picks a random path from the known nodes, sets roles, and returns that path.
func BuildCircuit(knownPorts []int, pathLength int) ([]*Node, error) {
    if pathLength < 2 || pathLength > len(knownPorts) {
        return nil, errors.New("invalid path length")
    }

    // Shuffle to get random selection
    shuffled := append([]int(nil), knownPorts...)
    rand.Shuffle(len(shuffled), func(i, j int) {
        shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
    })

    selectedPorts := shuffled[:pathLength]
    var nodes []*Node
    for i, p := range selectedPorts {
        node := &Node{
            ID:        fmt.Sprintf("Node%d", p),
            Port:      p,
            stopCh:    make(chan struct{}),
        }
        node.HtmlContent = fmt.Sprintf("<html><body><h1>Hello from port %d</h1></body></html>", p)

        // If user has overridden role via environment, use that (for demonstration).
        // Otherwise, assign roles automatically.
        envRole := os.Getenv("TOR_NODE_ROLE")
        if envRole != "" {
            // If the user sets a role, we do not override them with random or entry/exit logic.
            node.Role = NodeRole(envRole)
        } else {
            if i == 0 {
                node.Role = EntryNode
            } else if i == pathLength-1 {
                node.Role = ExitNode
            } else {
                node.Role = RelayNode
            }
        }
        nodes = append(nodes, node)
    }
    return nodes, nil
}

// forwardToAddress dials "addr" and sends the cell over TCP
func forwardToAddress(cell *RelayCell, addr string) error {
    conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
    if err != nil {
        return fmt.Errorf("dial next node (%s): %v", addr, err)
    }
    defer conn.Close()

    data, err := cell.Serialize()
    if err != nil {
        return fmt.Errorf("serialize cell: %v", err)
    }

    // write
    _, err = conn.Write(data)
    if err != nil {
        return fmt.Errorf("write to node (%s): %v", addr, err)
    }
    return nil
}

// SetupLogging configures log output to go to both a file and stdout (optional).
func SetupLogging(filename string) {
    // open log file
    logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatalf("Failed to open log file %s: %v", filename, err)
    }

    // Create multi-writer for stdout + file
    mw := io.MultiWriter(os.Stdout, logFile)
    log.SetOutput(mw)
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
