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


func TestBuildCircuit(t *testing.T) {
    nodes := []*protocol.Node{
        {ID: "Node1", Port: 9001},
        {ID: "Node2", Port: 9002},
        {ID: "Node3", Port: 9003},
    }
    circuit, err := protocol.BuildCircuit(nodes, 2)
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    if len(circuit) != 2 {
        t.Errorf("Expected 2 nodes in the circuit, got %d", len(circuit))
    }
    if circuit[0].Role != protocol.EntryNode || circuit[1].Role != protocol.ExitNode {
        t.Errorf("Incorrect roles assigned: %+v", circuit)
    }
}
