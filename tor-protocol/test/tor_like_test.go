package test

import (
	"testing"
	"time"

	"tor-protocol/config"
	"tor-protocol/protocol"
)

func TestTorLikeNetwork(t *testing.T) {
    cfg := config.GetConfig()

    // Bootstrap the nodes
    nodes := protocol.BootstrapNodes(cfg)
    for i, n := range nodes {
        go func(node *protocol.Node) {
            err := node.Start()
            if err != nil {
                t.Errorf("Error starting node %d: %v", i, err)
            }
        }(n)
    }

    // Wait a bit
    time.Sleep(2 * time.Second)

    // Build circuit
    circuit, err := protocol.BuildCircuit(nodes, 3)
    if err != nil {
        t.Fatalf("BuildCircuit error: %v", err)
    }

    // Send message
    testMessage := "Test message"
    err = protocol.SendThroughCircuit(circuit, testMessage)
    if err != nil {
        t.Errorf("SendThroughCircuit error: %v", err)
    }

    // Clean up
    time.Sleep(2 * time.Second)
    for _, n := range nodes {
        n.Stop()
    }
}
