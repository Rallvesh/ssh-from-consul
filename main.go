package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	consul "github.com/hashicorp/consul/api"
)

func main() {
	// Get Consul address and token from environment variables
	consulAddress := os.Getenv("CONSUL_HTTP_ADDR")
	if consulAddress == "" {
		consulAddress = "127.0.0.1:8500"
	}

	consulToken := os.Getenv("CONSUL_HTTP_TOKEN")
	if consulToken == "" {
		consulToken = ""
	}

	command := os.Args[1]

	// Create client configuration
	config := consul.DefaultConfig()
	config.Address = consulAddress
	config.Token = consulToken

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

		// Convert to JSON
		jsonData, err := json.MarshalIndent(nodes, "", "  ")
		if err != nil {
			log.Fatalf("Error converting data to JSON: %v", err)
		}

		fmt.Println(string(jsonData))
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

		// Extract WAN IP
		wanIP, exists := node.Node.TaggedAddresses["wan_ipv4"]
		if !exists || wanIP == "" {
			log.Fatalf("WAN IP not found for node: %s", nodeName)
		}

		fmt.Printf("Connecting to node: %s, WAN IP: %s\n", nodeName, wanIP)

		// Connect via SSH
		cmd := exec.Command("ssh", wanIP)
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
