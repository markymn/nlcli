package provider

import (
	"strings"

	"github.com/markymn/nlcli/internal/history"
	"github.com/markymn/nlcli/internal/shell"
)

type Provider interface {
	Name() string
	Model() string
	GetCommand(userInput, cwd string, shellType shell.ShellType, hist *history.History) (string, error)
}

func DetectProvider(apiKey string) (primary string, fallbacks []string) {
	switch {
	case strings.HasPrefix(apiKey, "sk-proj-"):
		return "openai", nil
	case strings.HasPrefix(apiKey, "sk-ant-"):
		return "anthropic", nil
	case strings.HasPrefix(apiKey, "sk-"):
		return "openai", []string{"anthropic"}
	case strings.HasPrefix(apiKey, "AIza"):
		return "google", nil
	case strings.HasPrefix(apiKey, "gsk_"):
		return "groq", nil
	default:
		// Attempt to verify against all providers
		if provider := VerifyKey(apiKey); provider != "" {
			return provider, nil
		}
		return "", nil // Unknown key format and failed verification
	}
}

// VerifyKey tries to use the key with each provider to see if it works
func VerifyKey(apiKey string) string {
	providers := []string{"openai", "anthropic", "google", "groq"}
	for _, p := range providers {
		// We use FetchModels as a lightweight auth check
		if _, err := FetchModels(p, apiKey); err == nil {
			return p
		}
	}
	return ""
}

// FetchModels dispatches to the correct provider's fetch function
func FetchModels(providerName, apiKey string) ([]string, error) {
	switch providerName {
	case "openai":
		return FetchOpenAIModels(apiKey)
	case "anthropic":
		return FetchAnthropicModels(apiKey)
	case "google":
		return FetchGoogleModels(apiKey)
	case "groq":
		return FetchGroqModels(apiKey)
	default:
		return nil, nil
	}
}

func GetProviderDisplayName(provider string) string {
	switch provider {
	case "openai":
		return "OpenAI"
	case "anthropic":
		return "Anthropic"
	case "google":
		return "Google"
	case "groq":
		return "Groq"
	default:
		return provider
	}
}

func GetModels(provider string) []string {
	switch provider {
	case "openai":
		return []string{"gpt-4o", "gpt-4o-mini", "gpt-4-turbo", "gpt-3.5-turbo"}
	case "anthropic":
		return []string{"claude-sonnet-4-20250514", "claude-3-5-haiku-20241022", "claude-3-haiku-20240307"}
	case "google":
		return []string{"gemini-2.5-flash", "gemini-2.0-flash", "gemini-1.5-flash", "gemini-1.5-pro"}
	case "groq":
		return []string{"llama-3.3-70b-versatile", "llama-3.1-8b-instant", "mixtral-8x7b-32768", "gemma2-9b-it"}
	default:
		return []string{}
	}
}

type MultiClient struct {
	apiKey    string
	model     string
	primary   Provider
	fallbacks []Provider
}

func NewMultiClient(apiKey, model, primaryName string, fallbackNames []string) *MultiClient {
	m := &MultiClient{apiKey: apiKey, model: model}

	m.primary = createProvider(primaryName, apiKey, model)

	for _, name := range fallbackNames {
		if p := createProvider(name, apiKey, model); p != nil {
			m.fallbacks = append(m.fallbacks, p)
		}
	}

	return m
}

func createProvider(name, apiKey, model string) Provider {
	switch name {
	case "openai":
		return NewOpenAI(apiKey, model)
	case "anthropic":
		return NewAnthropic(apiKey, model)
	case "google":
		return NewGoogle(apiKey, model)
	case "groq":
		return NewGroq(apiKey, model)
	default:
		return nil
	}
}

func (m *MultiClient) GetCommand(userInput, cwd string, shellType shell.ShellType, hist *history.History) (string, error) {
	cmd, err := m.primary.GetCommand(userInput, cwd, shellType, hist)
	if err == nil {
		return cmd, nil
	}

	for _, fb := range m.fallbacks {
		cmd, err = fb.GetCommand(userInput, cwd, shellType, hist)
		if err == nil {
			return cmd, nil
		}
	}

	return "", err
}

func (m *MultiClient) PrimaryName() string {
	return m.primary.Name()
}

func (m *MultiClient) PrimaryModel() string {
	return m.primary.Model()
}
