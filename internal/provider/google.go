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

type Google struct {
	apiKey string
	model  string
}

func NewGoogle(apiKey, model string) *Google {
	if model == "" {
		model = "gemini-2.5-flash"
	}
	return &Google{apiKey: apiKey, model: model}
}

func (g *Google) Name() string {
	return "Google"
}

func (g *Google) Model() string {
	return g.model
}

func (g *Google) GetCommand(userInput, cwd string, shellType shell.ShellType, hist *history.History) (string, error) {
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
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	})

	url := "https://generativelanguage.googleapis.com/v1beta/models/" + g.model + ":generateContent?key=" + g.apiKey
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
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
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response")
	}

	return strings.TrimSpace(result.Candidates[0].Content.Parts[0].Text), nil
}

func FetchGoogleModels(apiKey string) ([]string, error) {
	req, _ := http.NewRequest("GET", "https://generativelanguage.googleapis.com/v1beta/models?key="+apiKey, nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var models []string
	for _, m := range result.Models {
		name := strings.TrimPrefix(m.Name, "models/")
		models = append(models, name)
	}
	return models, nil
}
