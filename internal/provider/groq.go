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

type Groq struct {
	apiKey string
	model  string
	client *http.Client
}

func NewGroq(apiKey, model string) *Groq {
	if model == "" {
		model = "llama-3.3-70b-versatile"
	}
	return &Groq{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{},
	}
}

func (g *Groq) Name() string {
	return "Groq"
}

func (g *Groq) Model() string {
	return g.model
}

type groqRequest struct {
	Model     string    `json:"model"`
	Messages  []groqMsg `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type groqMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (g *Groq) GetCommand(userInput, cwd string, shellType shell.ShellType, hist *history.History) (string, error) {
	prompt := fmt.Sprintf(`You are a command line assistant.
Target Shell: %s
OS: Windows
Current Directory: %s

Previous commands:
%s

User Request: %s

Output only the shell command to execute. No markdown, no explanations.`,
		shell.GetShellName(shellType), cwd, hist.Format(), userInput)

	reqBody := groqRequest{
		Model:     g.model,
		MaxTokens: 150,
		Messages: []groqMsg{
			{Role: "user", Content: prompt},
		},
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.apiKey)

	resp, err := g.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var result groqResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	if result.Error != nil {
		return "", fmt.Errorf("%s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}

func FetchGroqModels(apiKey string) ([]string, error) {
	req, _ := http.NewRequest("GET", "https://api.groq.com/openai/v1/models", nil)
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
