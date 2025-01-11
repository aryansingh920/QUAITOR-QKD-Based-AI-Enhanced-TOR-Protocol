package test

import (
	// "fmt"
	"testing"
	"time"

	"tor-protocol/config"
	"tor-protocol/protocol"
)

func TestTorLikeNetwork(t *testing.T) {
    cfg := config.GetConfig()

    // Bootstrap the nodes
    nodes := protocol.BootstrapNodes(cfg)
    t.Logf("Bootstrapped %d nodes", len(nodes))

    for i, n := range nodes {
        go func(node *protocol.Node) {
            err := node.Start()
            if err != nil {
                t.Errorf("Error starting node %d (%s): %v", i, node.ID, err)
            } else {
                t.Logf("Node %s started successfully", node.ID)
            }
        }(n)
    }

    // Wait a bit
    time.Sleep(2 * time.Second)

    // Build circuit
    circuitLength := 5
    circuit, err := protocol.BuildCircuit(nodes, circuitLength)
    if err != nil {
        t.Fatalf("BuildCircuit error: %v", err)
    }
    t.Logf("Built circuit of length %d: %v", circuitLength, getNodeIDs(circuit))

    // Send message
    testMessage := "Test message from unit test"
    err = protocol.SendThroughCircuit(circuit, testMessage)
    if err != nil {
        t.Errorf("SendThroughCircuit error: %v", err)
    } else {
        t.Log("Message sent successfully through the circuit")
    }

    // Clean up
    time.Sleep(2 * time.Second)
    for _, n := range nodes {
        n.Stop()
    }
    t.Log("All nodes stopped successfully")
}

// Helper function to extract node IDs for logging
func getNodeIDs(nodes []*protocol.Node) []string {
    ids := make([]string, len(nodes))
    for i, n := range nodes {
        ids[i] = n.ID
    }
    return ids
}


func TestBuildCircuit(t *testing.T) {
    nodes := []*protocol.Node{
        {ID: "Node1", Port: 9001},
        {ID: "Node2", Port: 9002},
        {ID: "Node3", Port: 9003},
        {ID: "Node4", Port: 9004},
        {ID: "Node5", Port: 9005},
		{ID: "Node6", Port: 9006},
		{ID: "Node7", Port: 9007},
		{ID: "Node8", Port: 9008},
		{ID: "Node9", Port: 9009},
		{ID: "Node10", Port: 9010},
    }

    circuit, err := protocol.BuildCircuit(nodes, 3)
    if err != nil {
        t.Fatalf("Unexpected error while building circuit: %v", err)
    }

    if len(circuit) != 3 {
        t.Errorf("Expected circuit length 3, got %d", len(circuit))
    }

    if circuit[0].Role != protocol.EntryNode {
        t.Errorf("Expected first node to be ENTRY, got %s", circuit[0].Role)
    }

    if circuit[2].Role != protocol.ExitNode {
        t.Errorf("Expected last node to be EXIT, got %s", circuit[2].Role)
    }

    t.Logf("Circuit built successfully: %v", getNodeIDs(circuit))
}
