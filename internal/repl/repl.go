package repl

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/markymn/nlcli/internal/config"
	"github.com/markymn/nlcli/internal/history"
	"github.com/markymn/nlcli/internal/provider"
	"github.com/markymn/nlcli/internal/shell"
)

const (
	colorReset  = "\033[0m"
	colorPurple = "\033[35m" // Purple
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorRed    = "\033[31m"
	colorBold   = "\033[1m"
)

type REPL struct {
	client    *provider.MultiClient
	executor  *shell.Executor
	shellType shell.ShellType
	history   *history.History
	reader    *bufio.Reader
}

func New(client *provider.MultiClient, executor *shell.Executor, shellType shell.ShellType) *REPL {
	return &REPL{
		client:    client,
		executor:  executor,
		shellType: shellType,
		history:   history.New(),
		reader:    bufio.NewReader(os.Stdin),
	}
}

func (r *REPL) Start() {
	r.setupSignals()

	fmt.Printf("Shell:    %s%s%s\n", colorYellow, shell.GetShellName(r.shellType), colorReset)
	fmt.Printf("Provider: %s%s%s\n", colorYellow, r.client.PrimaryName(), colorReset)
	fmt.Printf("Model:    %s%s%s\n", colorYellow, r.client.PrimaryModel(), colorReset)
	fmt.Println()

	for {
		r.printPrompt()

		input, err := r.reader.ReadString('\n')
		if err != nil {
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		if r.handleSpecial(input) {
			continue
		}

		if strings.HasPrefix(input, "cd ") || input == "cd" {
			path := strings.TrimPrefix(input, "cd")
			if err := shell.ExecuteCD(path); err != nil {
				fmt.Printf("%s%s%s\n", colorRed, err, colorReset)
			}
			continue
		}

		if shell.IsValidSyntax(r.shellType, input) {
			r.runCommand(input)
			continue
		}

		r.translateAndRun(input)
	}
}

func (r *REPL) printPrompt() {
	cwd, _ := os.Getwd()
	fmt.Printf("%s%s%s > ", colorPurple, cwd, colorReset)
}

func (r *REPL) setupSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	go func() {
		for range c {
			fmt.Println()
			r.printPrompt()
		}
	}()
}

func (r *REPL) handleSpecial(input string) bool {
	switch strings.ToLower(input) {
	case ".exit":
		os.Exit(0)
	case ".help":
		r.showHelp()
		return true
	case ".api":
		r.changeAPI()
		return true
	case ".model":
		r.changeModel()
		return true
	case ".uninstall":
		r.uninstall()
		return true
	}
	return false
}

func (r *REPL) showHelp() {
	fmt.Printf("\n%s%snlcli - Natural Language CLI%s\n\n", colorBold, colorCyan, colorReset)
	fmt.Printf("Shell:    %s%s%s\n", colorYellow, shell.GetShellName(r.shellType), colorReset)
	fmt.Printf("Provider: %s%s%s\n", colorYellow, r.client.PrimaryName(), colorReset)
	fmt.Printf("Model:    %s%s%s\n\n", colorYellow, r.client.PrimaryModel(), colorReset)
	fmt.Println("Usage:")
	fmt.Println("  Type naturally   System translates to shell command")
	fmt.Println("  Type command     Runs directly (syntax validated)")
	fmt.Println("  cd <path>        Change directory")
	fmt.Println()
	fmt.Println("Special commands:")
	fmt.Println("  .help            Show this help")
	fmt.Println("  .api             Change API key and model")
	fmt.Println("  .model           Change model only")
	fmt.Println("  .uninstall       Remove nlcli")
	fmt.Println("  .exit            Exit nlcli")
	fmt.Println()
}

func (r *REPL) changeAPI() {
	var key string
	var providerName string
	var err error

	for {
		key, err = config.SetupAPIKey()
		if err != nil {
			fmt.Printf("%sError: %s%s\n", colorRed, err, colorReset)
			return
		}

		providerName, _ = provider.DetectProvider(key)
		if providerName != "" {
			break
		}

		fmt.Printf("%sInvalid API key format.%s\n", colorRed, colorReset)
		fmt.Println("Supported providers:")
		fmt.Println("  OpenAI (sk-...)")
		fmt.Println("  Anthropic (sk-ant-...)")
		fmt.Println("  Google (AIza...)")
		fmt.Println("  Groq (gsk_...)")
		fmt.Println()
	}

	displayName := provider.GetProviderDisplayName(providerName)

	fmt.Printf("\033[32mDetected provider: %s\033[0m\n", displayName)
	fmt.Printf("Fetching available models...\n")
	models, fetchErr := provider.FetchModels(providerName, key)
	if fetchErr != nil || len(models) == 0 {
		models = provider.GetModels(providerName)
	}
	model, err := config.SelectModel(models, displayName)
	if err != nil {
		model = models[0]
	}

	if err := config.SaveConfig(key, model); err != nil {
		fmt.Printf("%sError saving config: %s%s\n", colorRed, err, colorReset)
		return
	}

	r.client = provider.NewMultiClient(key, model, providerName, nil)
	fmt.Printf("%sSwitched to %s (%s)%s\n", colorPurple, r.client.PrimaryName(), r.client.PrimaryModel(), colorReset)
}

func (r *REPL) changeModel() {
	key, err := config.LoadAPIKey()
	if err != nil || key == "" {
		fmt.Printf("%sError: No API key found. Please use .api first.%s\n", colorRed, colorReset)
		return
	}

	providerName, _ := provider.DetectProvider(key)
	displayName := provider.GetProviderDisplayName(providerName)

	fmt.Printf("\033[32mDetected provider: %s\033[0m\n", displayName)
	fmt.Printf("Fetching available models...\n")
	models, fetchErr := provider.FetchModels(providerName, key)
	if fetchErr != nil || len(models) == 0 {
		models = provider.GetModels(providerName)
	}

	model, err := config.SelectModel(models, displayName)
	if err != nil {
		return // cancelled or error
	}

	if err := config.SaveConfig(key, model); err != nil {
		fmt.Printf("%sError saving config: %s%s\n", colorRed, err, colorReset)
		return
	}

	r.client = provider.NewMultiClient(key, model, providerName, nil)
	fmt.Printf("%sSwitched to %s (%s)%s\n", colorPurple, r.client.PrimaryName(), r.client.PrimaryModel(), colorReset)
}

func (r *REPL) uninstall() {
	fmt.Print("Remove nlcli and all data? (y/N): ")
	input, _ := r.reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(input)) != "y" {
		fmt.Println("Cancelled.")
		return
	}

	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".nlcli")
	os.RemoveAll(configDir)
	fmt.Println("Removed ~/.nlcli")
	fmt.Println("To complete uninstall, remove the nlcli binary from your PATH.")
	os.Exit(0)
}

func (r *REPL) runCommand(cmd string) {
	cmd = strings.TrimSpace(cmd)
	if strings.HasPrefix(cmd, "cd ") || cmd == "cd" {
		path := strings.TrimPrefix(cmd, "cd")
		path = strings.TrimSpace(path)
		if err := shell.ExecuteCD(path); err != nil {
			fmt.Printf("%s%s%s\n", colorRed, err, colorReset)
		}
		// Don't add simple navigation to history context to avoid noise, or maybe do?
		// For now let's add it so AI knows where we are if it looks at history.
		r.history.Add(cmd, "")
		return
	}

	err := r.executor.ExecuteInteractive(cmd)
	output := ""
	if err != nil {
		output = err.Error()
	}
	r.history.Add(cmd, output)
}

func (r *REPL) translateAndRun(input string) {
	cwd, _ := os.Getwd()

	cmd, err := r.client.GetCommand(input, cwd, r.shellType, r.history)
	if err != nil {
		fmt.Printf("%sError: %s%s\n", colorRed, err, colorReset)
		return
	}

	fmt.Printf("  %s%s%s\n", colorYellow, cmd, colorReset)
	r.runCommand(cmd)
}
