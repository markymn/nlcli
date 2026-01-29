package shell

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type ShellType string

const (
	ShellPowerShell ShellType = "powershell"
	ShellBash       ShellType = "bash"
	ShellZsh        ShellType = "zsh"
	ShellCmd        ShellType = "cmd"
	ShellFish       ShellType = "fish"
)

func DetectShell() ShellType {
	if shell := os.Getenv("SHELL"); shell != "" {
		switch {
		case strings.Contains(shell, "zsh"):
			return ShellZsh
		case strings.Contains(shell, "bash"):
			return ShellBash
		case strings.Contains(shell, "fish"):
			return ShellFish
		}
	}

	if os.Getenv("PSModulePath") != "" {
		return ShellPowerShell
	}

	if runtime.GOOS == "windows" {
		if os.Getenv("PROMPT") != "" && os.Getenv("PSModulePath") == "" {
			return ShellCmd
		}
		return ShellPowerShell
	}

	return ShellBash
}

func GetShellName(st ShellType) string {
	switch st {
	case ShellPowerShell:
		return "PowerShell"
	case ShellBash:
		return "Bash"
	case ShellZsh:
		return "Zsh"
	case ShellCmd:
		return "Command Prompt"
	case ShellFish:
		return "Fish"
	default:
		return "Shell"
	}
}

func IsValidSyntax(st ShellType, input string) bool {
	if strings.TrimSpace(input) == "" {
		return false
	}

	switch st {
	case ShellBash:
		if !checkShellSyntax("bash", "-n", "-c", input) {
			return false
		}
		if strings.ContainsAny(input, ";|&{}$=") {
			return true
		}
		fields := strings.Fields(input)
		if len(fields) == 0 {
			return false
		}
		checkCmd := fmt.Sprintf("type -t %s >/dev/null 2>&1", fields[0])
		return checkShellSyntax("bash", "-c", checkCmd)

	case ShellZsh:
		if !checkShellSyntax("zsh", "-n", "-c", input) {
			return false
		}
		if strings.ContainsAny(input, ";|&{}$=") {
			return true
		}
		fields := strings.Fields(input)
		if len(fields) == 0 {
			return false
		}
		checkCmd := fmt.Sprintf("whence -w %s >/dev/null 2>&1", fields[0])
		return checkShellSyntax("zsh", "-c", checkCmd)
	case ShellFish:
		return checkShellSyntax("fish", "-n", "-c", input)
	case ShellPowerShell:
		script := "try { [ScriptBlock]::Create('" + strings.ReplaceAll(input, "'", "''") + "') | Out-Null } catch { exit 1 }"
		if !checkShellSyntax("powershell", "-NoProfile", "-NoLogo", "-Command", script) {
			return false
		}

		if strings.ContainsAny(input, ";|{}$=") {
			return true
		}

		fields := strings.Fields(input)
		if len(fields) == 0 {
			return false
		}
		cmdName := fields[0]

		checkCmd := fmt.Sprintf("if (Get-Command -Name '%s' -ErrorAction SilentlyContinue) { exit 0 } else { exit 1 }", strings.ReplaceAll(cmdName, "'", "''"))
		return checkShellSyntax("powershell", "-NoProfile", "-NoLogo", "-Command", checkCmd)

	case ShellCmd:
		return isValidCmdCommand(input)
	default:
		return false
	}
}

func checkShellSyntax(shell string, args ...string) bool {
	cmd := exec.Command(shell, args...)
	return cmd.Run() == nil
}

func isValidCmdCommand(input string) bool {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return false
	}
	first := strings.ToLower(fields[0])
	known := []string{"dir", "cd", "copy", "move", "del", "mkdir", "rmdir", "type",
		"echo", "set", "cls", "exit", "call", "start", "ren", "rename", "attrib",
		"find", "findstr", "more", "sort", "xcopy", "robocopy", "tasklist", "taskkill"}
	for _, c := range known {
		if first == c {
			return true
		}
	}
	return false
}
