package provider

import (
	"testing"
)

func TestDetectProvider(t *testing.T) {
	tests := []struct {
		name              string
		apiKey            string
		expectedPrimary   string
		expectedFallbacks []string
	}{
		{
			name:            "OpenAI Project Key",
			apiKey:          "sk-proj-12345",
			expectedPrimary: "openai",
		},
		{
			name:            "Anthropic Key",
			apiKey:          "sk-ant-12345",
			expectedPrimary: "anthropic",
		},
		{
			name:              "Legacy OpenAI Key",
			apiKey:            "sk-12345",
			expectedPrimary:   "openai",
			expectedFallbacks: []string{"anthropic"},
		},
		{
			name:            "Google Key",
			apiKey:          "AIza12345",
			expectedPrimary: "google",
		},
		{
			name:            "Groq Key",
			apiKey:          "gsk_12345",
			expectedPrimary: "groq",
		},
		{
			name:            "Unknown Prefix",
			apiKey:          "unknown_prefix_12345",
			expectedPrimary: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			primary, fallbacks := DetectProvider(tt.apiKey)
			if primary != tt.expectedPrimary {
				t.Errorf("DetectProvider() primary = %v, want %v", primary, tt.expectedPrimary)
			}
			if len(fallbacks) != len(tt.expectedFallbacks) {
				t.Errorf("DetectProvider() fallbacks length = %v, want %v", len(fallbacks), len(tt.expectedFallbacks))
			}
			for i, fb := range fallbacks {
				if fb != tt.expectedFallbacks[i] {
					t.Errorf("DetectProvider() fallback[%d] = %v, want %v", i, fb, tt.expectedFallbacks[i])
				}
			}
		})
	}
}
