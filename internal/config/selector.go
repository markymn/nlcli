package config

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

const (
	keyUp    = 'A'
	keyDown  = 'B'
	keyEnter = '\r'
	keyLF    = '\n'
	keyESC   = 27
)

func SelectModel(models []string, providerName string) (string, error) {
	if len(models) == 0 {
		return "", fmt.Errorf("no models available")
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return models[0], nil
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	selected := 0
	buf := make([]byte, 3)

	fmt.Printf("\r\n%s models (use ↑↓ arrows, Enter to select):\r\n", providerName)
	printOptions(models, selected)

	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return models[selected], nil
		}

		if n == 1 {
			switch buf[0] {
			case keyEnter, keyLF:
				clearOptions(len(models) + 1)
				return models[selected], nil
			case 'q', 3:
				clearOptions(len(models) + 1)
				return models[0], nil
			}
		} else if n == 3 && buf[0] == keyESC && buf[1] == '[' {
			switch buf[2] {
			case keyUp:
				if selected > 0 {
					selected--
				}
			case keyDown:
				if selected < len(models)-1 {
					selected++
				}
			}
			clearOptions(len(models))
			printOptions(models, selected)
		}
	}
}

func printOptions(models []string, selected int) {
	for i, model := range models {
		if i == selected {
			fmt.Printf("\r  \033[36m▸ %s\033[0m\r\n", model)
		} else {
			fmt.Printf("\r    %s\r\n", model)
		}
	}
}

func clearOptions(lines int) {
	for i := 0; i < lines; i++ {
		fmt.Print("\033[A\033[2K")
	}
}
