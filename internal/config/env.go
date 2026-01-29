package config

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
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

func LoadSafetyLevel() int {
	data, err := os.ReadFile(envPath)
	if err != nil {
		return 1
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "SAFETY_LEVEL=") {
			levelStr := strings.TrimPrefix(line, "SAFETY_LEVEL=")
			if level, err := strconv.Atoi(levelStr); err == nil {
				return level
			}
		}
	}
	return 1
}

func SaveConfig(key, model string, safety int) error {
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}
	content := fmt.Sprintf("API_KEY=%s\nMODEL=%s\nSAFETY_LEVEL=%d\n", key, model, safety)
	return os.WriteFile(envPath, []byte(content), 0600)
}

func SaveAPIKey(key string) error {
	model, _ := LoadModel()
	safety := LoadSafetyLevel()
	return SaveConfig(key, model, safety)
}

func SaveModel(model string) error {
	key, _ := LoadAPIKey()
	safety := LoadSafetyLevel()
	return SaveConfig(key, model, safety)
}

func SaveSafetyLevel(level int) error {
	key, _ := LoadAPIKey()
	model, _ := LoadModel()
	return SaveConfig(key, model, level)
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

func RemoveFromPath() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	binDir := filepath.Dir(exe)

	if runtime.GOOS == "windows" {
		script := fmt.Sprintf(`$currentPath = [Environment]::GetEnvironmentVariable("Path", "User"); $pathParts = $currentPath -split ';' | Where-Object { $_ -and $_ -ne '%s' }; [Environment]::SetEnvironmentVariable("Path", ($pathParts -join ';'), "User")`, strings.ReplaceAll(binDir, "'", "''"))
		return exec.Command("powershell", "-NoProfile", "-Command", script).Run()
	}

	home, _ := os.UserHomeDir()
	shellConfigs := []string{".zshrc", ".bashrc", ".bash_profile"}

	for _, cfg := range shellConfigs {
		path := filepath.Join(home, cfg)
		if _, err := os.Stat(path); err != nil {
			continue
		}

		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		lines := strings.Split(string(data), "\n")
		var newLines []string
		skip := false
		for _, line := range lines {
			if strings.Contains(line, "# Added by nlcli installer") {
				skip = true
				continue
			}
			if skip && strings.Contains(line, "export PATH=") && strings.Contains(line, binDir) {
				skip = false
				continue
			}
			if skip {
				skip = false
			}
			newLines = append(newLines, line)
		}
		os.WriteFile(path, []byte(strings.Join(newLines, "\n")), 0644)
	}

	return nil
}
