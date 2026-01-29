package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var configDir string
var envPath string

func init() {
	home, _ := os.UserHomeDir()
	configDir = filepath.Join(home, ".nlcli")
	envPath = filepath.Join(configDir, ".env")
}

func LoadAPIKey() (string, error) {
	data, err := os.ReadFile(envPath)
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "API_KEY=") {
			return strings.TrimPrefix(line, "API_KEY="), nil
		}
	}
	return "", fmt.Errorf("API_KEY not found")
}

func LoadModel() (string, error) {
	data, err := os.ReadFile(envPath)
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "MODEL=") {
			return strings.TrimPrefix(line, "MODEL="), nil
		}
	}
	return "", fmt.Errorf("MODEL not found")
}

func SaveConfig(key, model string) error {
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}
	content := "API_KEY=" + key + "\nMODEL=" + model + "\n"
	return os.WriteFile(envPath, []byte(content), 0600)
}

func SaveAPIKey(key string) error {
	model, _ := LoadModel()
	return SaveConfig(key, model)
}

func SaveModel(model string) error {
	key, _ := LoadAPIKey()
	return SaveConfig(key, model)
}

func SetupAPIKey() (string, error) {
	fmt.Print("\nEnter your API key:\n> ")

	reader := bufio.NewReader(os.Stdin)
	key, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	key = strings.TrimSpace(key)
	if key == "" {
		return "", fmt.Errorf("empty API key")
	}

	return key, nil
}

func LoadOrSetupAPIKey() (string, error) {
	key, err := LoadAPIKey()
	if err == nil && key != "" {
		return key, nil
	}
	return SetupAPIKey()
}
