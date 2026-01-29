package history

import (
	"fmt"
	"strings"
)

type Entry struct {
	Command string
	Output  string
}

type History struct {
	entries    []Entry
	maxEntries int
	maxChars   int
}

func New() *History {
	return &History{
		entries:    make([]Entry, 0, 10),
		maxEntries: 10,
		maxChars:   4000,
	}
}

func (h *History) Add(command, output string) {
	if len(output) > 500 {
		output = output[:500] + "..."
	}

	h.entries = append(h.entries, Entry{Command: command, Output: output})

	if len(h.entries) > h.maxEntries {
		h.entries = h.entries[1:]
	}
}

func (h *History) Format() string {
	if len(h.entries) == 0 {
		return "(no history)"
	}

	var sb strings.Builder
	total := 0
	start := len(h.entries) - 5
	if start < 0 {
		start = 0
	}

	for i := start; i < len(h.entries); i++ {
		e := h.entries[i]
		line := fmt.Sprintf("$ %s\n%s\n", e.Command, e.Output)
		if total+len(line) > h.maxChars {
			break
		}
		sb.WriteString(line)
		total += len(line)
	}

	return sb.String()
}

func (h *History) Last() *Entry {
	if len(h.entries) == 0 {
		return nil
	}
	return &h.entries[len(h.entries)-1]
}
