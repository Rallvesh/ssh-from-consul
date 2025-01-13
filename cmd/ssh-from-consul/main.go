package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	consul "github.com/hashicorp/consul/api"
)

type Config struct {
	ConsulHTTPAddr  string `json:"consul_http_addr"`
	ConsulHTTPToken string `json:"consul_http_token"`
	Username        string `json:"username,omitempty"`
	PrivateKeyPath  string `json:"private_key_path,omitempty"`
}

type ConfigFile []map[string]Config

const defaultConfig = `[  
  {  
    "default": {  
      "consul_http_addr": "http://127.0.0.1:8500",  
      "consul_http_token": "",  
      "private_key_path": "",  
      "username": ""  
    }  
  }  
]`

func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Cannot determine home directory: %v", err)
	}
	return filepath.Join(homeDir, ".config", "sfc", "sfc.json")
}

func ensureConfigFileExists(configPath string) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("Config file not found, creating default config at:", configPath)

		err := os.MkdirAll(filepath.Dir(configPath), 0755)
		if err != nil {
			log.Fatalf("Error creating config directory: %v", err)
		}

		err = os.WriteFile(configPath, []byte(defaultConfig), 0644)
		if err != nil {
			log.Fatalf("Error writing default config: %v", err)
		}
	}
}

func loadConfig(profile string) (Config, error) {
	configPath := getConfigPath()
	ensureConfigFileExists(configPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %v", err)
	}

	var configFile ConfigFile
	err = json.Unmarshal(data, &configFile)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing config file: %v", err)
	}

	for _, configMap := range configFile {
		if config, exists := configMap[profile]; exists {
			return config, nil
		}
	}

	return Config{}, fmt.Errorf("profile '%s' not found in config", profile)
}

func getDefaultUsername() string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("Error getting current user: %v", err)
	}
	return currentUser.Username
}

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

	config, err := loadConfig(profile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	clientConfig := consul.DefaultConfig()
	clientConfig.Address = config.ConsulHTTPAddr
	clientConfig.Token = config.ConsulHTTPToken

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

		sshArgs := []string{wanIP}

		// Если username указан в конфиге — используем его, иначе берем системный
		username := config.Username
		if username == "" {
			username = getDefaultUsername()
		}

		// Если указан ключ, добавляем его в ssh команду
		if config.PrivateKeyPath != "" {
			sshArgs = append([]string{"-i", config.PrivateKeyPath, username + "@" + wanIP})
		} else {
			sshArgs = append([]string{username + "@" + wanIP})
		}

		cmd := exec.Command("ssh", sshArgs...)
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
