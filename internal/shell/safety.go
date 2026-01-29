package shell

import (
	"strings"
)

type SafetyLevel int

const (
	SafetyInstant  SafetyLevel = 1
	SafetyLax      SafetyLevel = 2
	SafetyCautious SafetyLevel = 3
	SafetyStrict   SafetyLevel = 4
)

func (s SafetyLevel) String() string {
	switch s {
	case SafetyInstant:
		return "Instant"
	case SafetyLax:
		return "Lax"
	case SafetyCautious:
		return "Cautious"
	case SafetyStrict:
		return "Strict"
	default:
		return "Unknown"
	}
}

func IsDangerous(cmd string, level SafetyLevel) bool {
	if level == SafetyInstant {
		return false
	}
	if level == SafetyStrict {
		return true
	}

	cmd = strings.ToLower(strings.TrimSpace(cmd))
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return false
	}

	baseCmd := parts[0]

	destructive := []string{
		"rm", "del", "rd", "rmdir",
		"format", "mkfs", "fdisk", "shred",
		"reset", "reboot", "shutdown",
	}

	for _, d := range destructive {
		if baseCmd == d {
			if baseCmd == "rm" && (strings.Contains(cmd, "-r") || strings.Contains(cmd, "-f")) {
				return true
			}
			if level >= SafetyLax {
				return true
			}
		}
	}

	if level == SafetyLax {
		if baseCmd == "rm" && (strings.Contains(cmd, "-rf") || strings.Contains(cmd, "/*")) {
			return true
		}
		return false
	}

	if level == SafetyCautious {
		writeOps := []string{
			"mkdir", "touch", "cp", "mv", "ln", "git push",
			"dd", "wget", "curl", "pip install", "npm install",
			"apt", "yum", "dnf", "brew", "pacman",
			"chmod", "chown",
		}

		for _, w := range writeOps {
			if strings.Contains(cmd, w) {
				return true
			}
		}

		if strings.Contains(cmd, ">") || strings.Contains(cmd, ">>") || strings.Contains(cmd, "|") {
			return true
		}
	}

	return false
}
