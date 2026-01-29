# nlcli 🚀

**nlcli** is a natural language interface for your terminal. It translates natural language requests into shell commands using powerful AI models from OpenAI, Anthropic, Google Gemini, and Groq.

## Features

- **Natural Language Translation**: Type what you want to do (e.g., "list all files larger than 10MB") and get the corresponding shell command.
- **Smart Execution**: Validates shell syntax and runs commands directly if they are already valid.
- **Multi-Provider Support**: Seamlessly switch between OpenAI, Anthropic, Google Gemini, and Groq.
- **Context-Aware**: Remembers previous commands to provide better translations.
- **Cross-Platform**: Designed for Windows (Powershell/Cmd) and Unix-like systems (Bash/Zsh/Fish).

## Installation

### Quick Install (Windows - PowerShell)
```powershell
irm https://raw.githubusercontent.com/markymn/nlcli/main/install/install-remote.ps1 | iex
```

### Quick Install (Unix - Bash/Zsh)
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
- Run commands directly: `ls -la`
- Special commands:
    - `.help`: Show help menu
    - `.api`: Change API key and provider
    - `.model`: Change the AI model
    - `.exit`: Quit the terminal

## Supported Providers

- **OpenAI**: Requires `sk-...`
- **Anthropic**: Requires `sk-ant-...`
- **Google Gemini**: Requires `AIza...`
- **Groq**: Requires `gsk_...`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
