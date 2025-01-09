package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	consul "github.com/hashicorp/consul/api"
)

type Config struct {
	ConsulHTTPAddr  string `json:"consul_http_addr"`
	ConsulHTTPToken string `json:"consul_http_token"`
	Username        string `json:"username"`
	PrivateKeyPath  string `json:"private_key_path"`
}

type ConfigFile []map[string]Config

func loadConfig(filePath, profile string) (Config, error) {
	// Read config file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %v", err)
	}

	var configFile ConfigFile
	err = json.Unmarshal(data, &configFile)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing config file: %v", err)
	}

	// Find profile in config
	for _, configMap := range configFile {
		if config, exists := configMap[profile]; exists {
			return config, nil
		}
	}

	return Config{}, fmt.Errorf("profile '%s' not found in config", profile)
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <profile> <command> [args]", os.Args[0])
	}

	profile := os.Args[1]
	command := os.Args[2]

	// Load config file
	config, err := loadConfig("sfc.json", profile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Create client configuration
	clientConfig := consul.DefaultConfig()
	clientConfig.Address = config.ConsulHTTPAddr
	clientConfig.Token = config.ConsulHTTPToken

	// Create Consul client
	client, err := consul.NewClient(clientConfig)
	if err != nil {
		log.Fatalf("Error creating Consul client: %v", err)
	}

	if command == "ls" {
		// Retrieve list of nodes from Consul
		nodes, _, err := client.Catalog().Nodes(nil)
		if err != nil {
			log.Fatalf("Error retrieving node list: %v", err)
		}

		// Convert list of nodes to JSON
		jsonData, err := json.MarshalIndent(nodes, "", "  ")
		if err != nil {
			log.Fatalf("Error converting data to JSON: %v", err)
		}

		fmt.Println(string(jsonData))
	} else if command == "connect" {
		if len(os.Args) < 4 {
			log.Fatalf("Please provide a node name for the connect command")
		}
		nodeName := os.Args[3]

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
		cmd := exec.Command("ssh", "-i", config.PrivateKeyPath, config.Username+"@"+wanIP)
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
