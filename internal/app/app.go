package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"

	consul "github.com/hashicorp/consul/api"
	"github.com/rallvesh/ssh-from-consul/internal/config"
	"github.com/rallvesh/ssh-from-consul/internal/ssh"
)

// Версия программы
const version = "v1.0.0"

// Выводим текущую версию приложения
func ShowVersion() {
	fmt.Printf("%s %s\n", os.Args[0], version)
	fmt.Printf("on %s_%s\n", runtime.GOOS, runtime.GOARCH)
}

// HandleCommand обрабатывает команду `ls` или `connect` с параметрами.
func HandleCommand(command string, profile string) {
	// Загружаем конфиг для указанного профиля
	cfg, err := config.LoadConfig(profile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Конфигурируем клиент Consul
	clientConfig := consul.DefaultConfig()
	clientConfig.Address = cfg.ConsulHTTPAddr
	clientConfig.Token = cfg.ConsulHTTPToken

	client, err := consul.NewClient(clientConfig)
	if err != nil {
		log.Fatalf("Error creating Consul client: %v", err)
	}

	// Обрабатываем команду
	switch command {
	case "ls":
		listNodes(client)
	case "connect":
		if len(os.Args) < 4 {
			log.Fatalf("Please provide a node name for the connect command")
		}
		nodeName := os.Args[3]
		connectToNode(client, nodeName, cfg)
	default:
		log.Fatalf("Unknown command: %s", command)
	}
}

// listNodes получает и выводит список узлов в формате JSON.
func listNodes(client *consul.Client) {
	nodes, _, err := client.Catalog().Nodes(nil)
	if err != nil {
		log.Fatalf("Error retrieving node list: %v", err)
	}

	// Преобразуем в JSON
	jsonData, err := json.MarshalIndent(nodes, "", "  ")
	if err != nil {
		log.Fatalf("Error converting data to JSON: %v", err)
	}

	// Печатаем JSON
	fmt.Println(string(jsonData))
}

// connectToNode подключается к узлу через SSH, используя параметры из конфига.
func connectToNode(client *consul.Client, nodeName string, cfg config.Config) {
	// Получаем информацию об узле
	node, _, err := client.Catalog().Node(nodeName, nil)
	if err != nil || node == nil {
		log.Fatalf("Error retrieving node information: %v", err)
	}

	// Получаем WAN IP узла
	wanIP, exists := node.Node.TaggedAddresses["wan_ipv4"]
	if !exists || wanIP == "" {
		log.Fatalf("WAN IP not found for node: %s", nodeName)
	}

	fmt.Printf("Connecting to node: %s, WAN IP: %s\n", nodeName, wanIP)

	// Подключаемся через SSH
	sshClient := ssh.SSHClient{
		Username:       cfg.Username,
		PrivateKeyPath: cfg.PrivateKeyPath,
	}

	err = sshClient.Connect(wanIP)
	if err != nil {
		log.Fatalf("Error connecting via SSH: %v", err)
	}
}
