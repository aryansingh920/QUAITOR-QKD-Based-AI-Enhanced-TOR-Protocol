package protocol

import (
	"fmt"
	"log"
	mrand "math/rand"
	"net"
	"time"
	"tor-protocol/config"
)

// RelayCell is a simple container for a message and next address
type RelayCell struct {
    NextAddr string
    Payload  []byte
}

// BuildCircuit picks a random set of nodes for the path and assigns roles (entry, relays, exit).
func BuildCircuit(nodes []*Node, pathLength int) ([]*Node, error) {
    if pathLength < 2 || pathLength > len(nodes) {
        return nil, fmt.Errorf("invalid path length: %d", pathLength)
    }

    // Shuffle nodes and pick the first 'pathLength' from the random set
    shuffled := make([]*Node, len(nodes))
    copy(shuffled, nodes)
    mrand.Shuffle(len(shuffled), func(i, j int) {
        shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
    })
    path := shuffled[:pathLength]

    // Assign roles: first => ENTRY, last => EXIT, middle => RELAY
    if pathLength == 2 {
        path[0].Role = EntryNode
        path[1].Role = ExitNode
    } else {
        path[0].Role = EntryNode
        for i := 1; i < pathLength-1; i++ {
            path[i].Role = RelayNode
        }
        path[pathLength-1].Role = ExitNode
    }

    return path, nil
}

// SendThroughCircuit sends a message through the circuit.
func SendThroughCircuit(circuit []*Node, message string) error {
    // We simulate onion routing by building relay cells.
    // In a real Tor, you'd wrap encrypted layers. Here, we just forward "raw" data for simplicity.

    if len(circuit) == 0 {
        return fmt.Errorf("circuit is empty")
    }

    // NextAddr for the first node is the second node's address
    // If only 2 nodes, first is entry, second is exit.
    // We'll build relay cells in a forward-chained manner:
    var prevNext string
    for i := 0; i < len(circuit)-1; i++ {
        nextNode := circuit[i+1]
        prevNext = fmt.Sprintf("127.0.0.1:%d", nextNode.Port)
    }

    // Create the final RelayCell for the entry node
    cell := RelayCell{
        NextAddr: prevNext,
        Payload:  []byte(message),
    }

    // Connect to the entry node and write
    entryNode := circuit[0]
    entryAddr := fmt.Sprintf("127.0.0.1:%d", entryNode.Port)
    conn, err := net.Dial("tcp", entryAddr)
    if err != nil {
        return fmt.Errorf("error dialing entry node: %v", err)
    }
    defer conn.Close()

    data, err := cell.Serialize()
    if err != nil {
        return fmt.Errorf("error serializing relay cell: %v", err)
    }

    // Send the cell to the entry node
    _, err = conn.Write(data)
    if err != nil {
        return fmt.Errorf("error writing to entry node: %v", err)
    }

    log.Printf("Sent message to entry node (%s): %q\n", entryNode.ID, message)
    return nil
}

// forwardToNextNode dials the next node in the chain and sends the relay cell
func forwardToNextNode(cell *RelayCell) error {
    if cell.NextAddr == "" {
        // This means we are at the exit node
        return nil
    }

    conn, err := net.DialTimeout("tcp", cell.NextAddr, 5*time.Second)
    if err != nil {
        return fmt.Errorf("dial next node: %v", err)
    }
    defer conn.Close()

    data, err := cell.Serialize()
    if err != nil {
        return fmt.Errorf("serialize: %v", err)
    }

    _, err = conn.Write(data)
    if err != nil {
        return fmt.Errorf("write: %v", err)
    }
    return nil
}

// ParseRelayCell converts raw bytes into a RelayCell structure
func ParseRelayCell(data []byte) (*RelayCell, error) {
    return DeserializeCell(data)
}


// BootstrapNodes creates a list of Node instances from the config
func BootstrapNodes(cfg *config.Config) []*Node {
    nodes := make([]*Node, len(cfg.Ports))
    for i, port := range cfg.Ports {
        node := &Node{
            ID:      fmt.Sprintf("Node%d", i+1),
            Role:    RelayNode, // default to RELAY, roles will be reassigned upon path creation
            Port:    port,
            stopCh:  make(chan struct{}),
        }
        nodes[i] = node
    }
    return nodes
}
