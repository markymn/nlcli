package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/markymn/nlcli/internal/history"
	"github.com/markymn/nlcli/internal/shell"
)

type Anthropic struct {
	apiKey string
	model  string
}

func NewAnthropic(apiKey, model string) *Anthropic {
	if model == "" {
		model = "claude-3-haiku-20240307"
	}
	return &Anthropic{apiKey: apiKey, model: model}
}

func (c *Anthropic) Name() string {
	return "Anthropic"
}

func (c *Anthropic) Model() string {
	return c.model
}

func (c *Anthropic) GetCommand(userInput, cwd string, shellType shell.ShellType, hist *history.History) (string, error) {
	prompt := fmt.Sprintf(`You are a command line assistant.
Target Shell: %s
OS: Windows
Current Directory: %s

Previous commands:
%s

User Request: %s

Output only the shell command to execute. No markdown, no explanations.`,
		shell.GetShellName(shellType), cwd, hist.Format(), userInput)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"model": c.model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"max_tokens": 100,
	})

	req, _ := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(reqBody))
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("api error: %s", resp.Status)
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("no response")
	}

	return strings.TrimSpace(result.Content[0].Text), nil
}

// FetchAnthropicModels retrieves available models from the Anthropic API
func FetchAnthropicModels(apiKey string) ([]string, error) {
	req, _ := http.NewRequest("GET", "https://api.anthropic.com/v1/models", nil)
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var models []string
	for _, m := range result.Data {
		models = append(models, m.ID)
	}
	return models, nil
}
