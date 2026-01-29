package shell

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Executor struct {
	shellType ShellType
	binary    string
	args      []string
}

func NewExecutor(st ShellType) *Executor {
	e := &Executor{shellType: st}

	switch st {
	case ShellPowerShell:
		if runtime.GOOS == "windows" {
			e.binary = "powershell.exe"
		} else {
			e.binary = "pwsh"
		}
		e.args = []string{"-NoProfile", "-NoLogo", "-Command"}
	case ShellBash:
		e.binary = "bash"
		e.args = []string{"-c"}
	case ShellZsh:
		e.binary = "zsh"
		e.args = []string{"-c"}
	case ShellFish:
		e.binary = "fish"
		e.args = []string{"-c"}
	case ShellCmd:
		e.binary = "cmd.exe"
		e.args = []string{"/C"}
	default:
		e.binary = "sh"
		e.args = []string{"-c"}
	}

	return e
}

func (e *Executor) Execute(command string) (string, string, error) {
	args := append(e.args, command)
	cmd := exec.Command(e.binary, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	if e.shellType == ShellBash && runtime.GOOS == "windows" {
		cmd.Env = append(cmd.Env, "MSYS_NO_PATHCONV=1")
	}

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func (e *Executor) ExecuteInteractive(command string) error {
	args := append(e.args, command)
	cmd := exec.Command(e.binary, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	if e.shellType == ShellBash && runtime.GOOS == "windows" {
		cmd.Env = append(cmd.Env, "MSYS_NO_PATHCONV=1")
	}
	return cmd.Run()
}

func ExecuteCD(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		path = home
	}
	return os.Chdir(path)
}
