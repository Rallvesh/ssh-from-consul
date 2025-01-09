package main

import (
	"fmt"
	"log"
	"os"

	consul "github.com/hashicorp/consul/api"
)

func main() {
	// Get Consul address from environment variable or use default value
	consulAddress := os.Getenv("CONSUL_ADDRESS")
	if consulAddress == "" {
		consulAddress = "127.0.0.1:8500"
	}

	// Create client configuration
	config := consul.DefaultConfig()
	config.Address = consulAddress

	// Create client
	client, err := consul.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating Consul client: %v", err)
	}

	// Retrieve list of nodes
	nodes, _, err := client.Catalog().Nodes(nil)
	if err != nil {
		log.Fatalf("Error retrieving node list: %v", err)
	}

	// Print node information
	for _, node := range nodes {
		fmt.Printf("Node: %s, Address: %s\n", node.Node, node.Address)
	}
}
