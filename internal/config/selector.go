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
	title := fmt.Sprintf("%s models", providerName)
	return SelectGeneric(models, title)
}

func SelectGeneric(options []string, title string) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("no options available")
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return options[0], nil
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	selected := 0
	buf := make([]byte, 3)

	fmt.Printf("\r\n%s (use ↑↓ arrows, Enter to select):\r\n", title)
	printOptions(options, selected)

	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return options[selected], nil
		}

		if n == 1 {
			switch buf[0] {
			case keyEnter, keyLF:
				clearOptions(len(options) + 1)
				return options[selected], nil
			case 'q', 3:
				clearOptions(len(options) + 1)
				return options[0], nil
			}
		} else if n == 3 && buf[0] == keyESC && buf[1] == '[' {
			switch buf[2] {
			case keyUp:
				if selected > 0 {
					selected--
				}
			case keyDown:
				if selected < len(options)-1 {
					selected++
				}
			}
			clearOptions(len(options))
			printOptions(options, selected)
		}
	}
}

func printOptions(options []string, selected int) {
	for i, option := range options {
		if i == selected {
			fmt.Printf("\r  \033[36m▸ %s\033[0m\r\n", option)
		} else {
			fmt.Printf("\r    %s\r\n", option)
		}
	}
}

func clearOptions(lines int) {
	for i := 0; i < lines; i++ {
		fmt.Print("\033[A\033[2K")
	}
}
