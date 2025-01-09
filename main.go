package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	consul "github.com/hashicorp/consul/api"
)

func main() {
	// Get Consul address from environment variable or use default value
	consulAddress := os.Getenv("CONSUL_ADDRESS")
	if consulAddress == "" {
		consulAddress = "127.0.0.1:8500"
	}

	// Check command arguments
	if len(os.Args) < 2 {
		log.Fatalf("Please provide a command (ls or connect <node_name>)")
	}

	command := os.Args[1]

	// Create client configuration
	config := consul.DefaultConfig()
	config.Address = consulAddress

	// Create client
	client, err := consul.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating Consul client: %v", err)
	}

	if command == "ls" {
		// Retrieve list of nodes
		nodes, _, err := client.Catalog().Nodes(nil)
		if err != nil {
			log.Fatalf("Error retrieving node list: %v", err)
		}

		// Print node information
		for _, node := range nodes {
			fmt.Printf("Node: %s, Address: %s\n", node.Node, node.Address)
		}
	} else if command == "connect" {
		if len(os.Args) < 3 {
			log.Fatalf("Please provide a node name for the connect command")
		}
		nodeName := os.Args[2]

		// Retrieve node information
		node, _, err := client.Catalog().Node(nodeName, nil)
		if err != nil || node == nil {
			log.Fatalf("Error retrieving node information: %v", err)
		}

		// Get node IP
		nodeIP := node.Node.Address
		fmt.Printf("Connecting to node: %s, IP: %s\n", nodeName, nodeIP)

		// Connect via SSH
		cmd := exec.Command("ssh", nodeIP)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()
		if err != nil {
			log.Fatalf("Error connecting via SSH: %v", err)
		}
	} else {
		log.Fatalf("Unknown command: %s", command)
	}
}
