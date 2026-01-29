# nlcli

**nlcli** is a natural language interface for your terminal. It translates natural language requests into shell commands using powerful AI models from OpenAI, Anthropic, Google Gemini, and Groq.

## Features

- **Natural Language Translation**: Type what you want to do (e.g., "list all files larger than 10MB") and get the corresponding shell command.
- **Smart Execution**: Validates shell syntax and runs commands directly if they are already valid.
- **Multi-Provider Support**: Seamlessly switch between OpenAI, Anthropic, Google Gemini, and Groq.
- **Context-Aware**: Remembers previous commands to provide better translations.
- **Cross-Platform**: Designed for Windows (Powershell/Cmd) and Unix-like systems (Bash/Zsh/Fish).

## Quick Start

Install and add `nlcli` to your PATH instantly by running one of these commands:

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/markymn/nlcli/main/install/install-remote.ps1 | iex
```

### macOS / Linux (Bash/Zsh)
```bash
curl -sSL https://raw.githubusercontent.com/markymn/nlcli/main/install/install-remote.sh | bash
```

## Usage

Start the interactive session:
```bash
nlcli
```

Inside the session:
- Type naturally: `show me my current directory`
- Run commands directly: `ls -la` (Automatically validated)
- Special commands:
    - `.help`: Show help menu
    - `.safety`: Rotate through 4 safety levels
    - `.api`: Change API key and provider
    - `.model`: Change the AI model
    - `.uninstall`: Completely remove nlcli and clean up PATH
    - `.exit`: Quit the terminal

## Safety System

To keep your terminal safe, `nlcli` features a 4-stage permission system:

1. **Instant (Default)**: Executes all translated commands immediately.
2. **Lax**: Prompts for confirmation only on high-danger commands (e.g., `rm -rf`).
3. **Cautious**: Prompts for any command that modifies the file system (write operations).
4. **Strict**: Always prompts for confirmation before any execution.

Switch levels anytime using the `.safety` command.

## Supported Providers

- OpenAI
- Anthropic
- Google Gemini
- Groq

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
