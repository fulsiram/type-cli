package main

import (
	"encoding/json"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fulsiram/type-cli/internal/app"
)

func main() {
	file, err := os.ReadFile("words.json")
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	var words []string
	json.Unmarshal(file, &words)

	p := tea.NewProgram(app.NewModel(words))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
