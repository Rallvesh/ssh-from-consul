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

	// Get node name from environment variable or argument
	var nodeName string
	if len(os.Args) > 1 {
		nodeName = os.Args[1]
	} else {
		log.Fatalf("Please provide a node name as an argument")
	}

	// Create client configuration
	config := consul.DefaultConfig()
	config.Address = consulAddress

	// Create client
	client, err := consul.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating Consul client: %v", err)
	}

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
}
