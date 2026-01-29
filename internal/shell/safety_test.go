package shell

import (
	"testing"
)

func TestIsDangerous(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		level    SafetyLevel
		expected bool
	}{
		{name: "Instant - rm", cmd: "rm -rf /", level: SafetyInstant, expected: false},
		{name: "Lax - ls", cmd: "ls", level: SafetyLax, expected: false},
		{name: "Lax - rm file", cmd: "rm file.txt", level: SafetyLax, expected: true},
		{name: "Lax - rm recursive", cmd: "rm -r folder", level: SafetyLax, expected: true},
		{name: "Lax - mkdir", cmd: "mkdir test", level: SafetyLax, expected: false},
		{name: "Cautious - ls", cmd: "ls", level: SafetyCautious, expected: false},
		{name: "Cautious - mkdir", cmd: "mkdir test", level: SafetyCautious, expected: true},
		{name: "Cautious - touch", cmd: "touch file", level: SafetyCautious, expected: true},
		{name: "Cautious - cat", cmd: "cat file", level: SafetyCautious, expected: false},
		{name: "Cautious - echo redirect", cmd: "echo hello > file", level: SafetyCautious, expected: true},

		{name: "Strict - ls", cmd: "ls", level: SafetyStrict, expected: true},
		{name: "Strict - pwd", cmd: "pwd", level: SafetyStrict, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDangerous(tt.cmd, tt.level); got != tt.expected {
				t.Errorf("IsDangerous(%q, %v) = %v, want %v", tt.cmd, tt.level, got, tt.expected)
			}
		})
	}
}
