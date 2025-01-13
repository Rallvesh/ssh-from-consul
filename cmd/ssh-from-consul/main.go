package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	consul "github.com/hashicorp/consul/api"
	"github.com/rallvesh/ssh-from-consul/internal/config"
	"github.com/rallvesh/ssh-from-consul/internal/ssh"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s [profile] <command> [args]", os.Args[0])
	}

	profile := "default"
	commandIndex := 1

	// Если первый аргумент - это профиль, используем его
	if len(os.Args) > 2 {
		profile = os.Args[1]
		commandIndex = 2
	}

	command := os.Args[commandIndex]

	cfg, err := config.LoadConfig(profile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	clientConfig := consul.DefaultConfig()
	clientConfig.Address = cfg.ConsulHTTPAddr
	clientConfig.Token = cfg.ConsulHTTPToken

	client, err := consul.NewClient(clientConfig)
	if err != nil {
		log.Fatalf("Error creating Consul client: %v", err)
	}

	if command == "ls" {
		nodes, _, err := client.Catalog().Nodes(nil)
		if err != nil {
			log.Fatalf("Error retrieving node list: %v", err)
		}

		jsonData, err := json.MarshalIndent(nodes, "", "  ")
		if err != nil {
			log.Fatalf("Error converting data to JSON: %v", err)
		}

		fmt.Println(string(jsonData))
	} else if command == "connect" {
		if len(os.Args) <= commandIndex+1 {
			log.Fatalf("Please provide a node name for the connect command")
		}
		nodeName := os.Args[commandIndex+1]

		node, _, err := client.Catalog().Node(nodeName, nil)
		if err != nil || node == nil {
			log.Fatalf("Error retrieving node information: %v", err)
		}

		wanIP, exists := node.Node.TaggedAddresses["wan_ipv4"]
		if !exists || wanIP == "" {
			log.Fatalf("WAN IP not found for node: %s", nodeName)
		}

		fmt.Printf("Connecting to node: %s, WAN IP: %s\n", nodeName, wanIP)

		sshClient := ssh.SSHClient{
			Username:       cfg.Username,
			PrivateKeyPath: cfg.PrivateKeyPath,
		}

		err = sshClient.Connect(wanIP)
		if err != nil {
			log.Fatalf("Error connecting via SSH: %v", err)
		}
	} else {
		log.Fatalf("Unknown command: %s", command)
	}
}
