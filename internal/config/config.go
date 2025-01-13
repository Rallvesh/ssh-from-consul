package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Config структура для загрузки параметров профиля.
type Config struct {
	ConsulHTTPAddr  string `json:"consul_http_addr"`
	ConsulHTTPToken string `json:"consul_http_token"`
	Username        string `json:"username,omitempty"`
	PrivateKeyPath  string `json:"private_key_path,omitempty"`
}

// ConfigFile представляет массив JSON-конфигураций с профилями.
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

// getConfigPath возвращает путь к файлу конфигурации в ~/.config/sfc/sfc.json.
func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Cannot determine home directory: %v", err)
	}
	return filepath.Join(homeDir, ".config", "sfc", "sfc.json")
}

// ensureConfigFileExists создаёт конфиг, если его нет.
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

// LoadConfig загружает профиль из JSON-файла.
func LoadConfig(profile string) (Config, error) {
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
