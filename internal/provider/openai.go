package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/markymn/nlcli/internal/history"
	"github.com/markymn/nlcli/internal/shell"
)

type OpenAI struct {
	apiKey string
	model  string
}

func NewOpenAI(apiKey, model string) *OpenAI {
	if model == "" {
		model = "gpt-4o-mini"
	}
	return &OpenAI{apiKey: apiKey, model: model}
}

func (c *OpenAI) Name() string {
	return "OpenAI"
}

func (c *OpenAI) Model() string {
	return c.model
}

func (c *OpenAI) GetCommand(userInput, cwd string, shellType shell.ShellType, hist *history.History) (string, error) {
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

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
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
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}

// FetchOpenAIModels retrieves available models from the OpenAI API
func FetchOpenAIModels(apiKey string) ([]string, error) {
	req, _ := http.NewRequest("GET", "https://api.openai.com/v1/models", nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)

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
	sort.Strings(models)
	return models, nil
}
