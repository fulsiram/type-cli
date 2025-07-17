package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keymap struct {
	quit      key.Binding
	nextWord  key.Binding
	backSpace key.Binding
}

type model struct {
	timer    timer.Model
	keymap   keymap
	quitting bool

	width  int
	height int

	words           []string
	lines           [][]string
	typedLines      [][]string
	currentLine     int
	currentLineWord int
	currentCharId   int
}

func generateLines(words []string, n int) [][]string {
	var lines [][]string

	for range n {
		lines = append(lines, generateRandomLine(words))
	}

	return lines
}

func initialModel(words []string) model {
	return model{
		timer: timer.NewWithInterval(30*time.Second, time.Second),
		keymap: keymap{
			quit: key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
			nextWord: key.NewBinding(
				key.WithKeys(" "),
			),
			backSpace: key.NewBinding(
				key.WithKeys("backspace"),
			),
		},
		words: words,
		lines: generateLines(words, 10),
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.nextWord):
			if m.currentCharId == 0 {
				return m, nil
			}

		default:

		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m model) View() string {
	header := lipgloss.NewStyle().
		MarginTop(1).
		Width(m.width).
		Align(lipgloss.Center).
		Render("Type CLI")

	toTypeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#646669"))

	lineText := ""
	for _, line := range m.lines[:3] {
		lineText += strings.Join(line, " ") + "\n"
	}

	content := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-2).
		Align(lipgloss.Center, lipgloss.Center).
		Render(toTypeStyle.Render(lineText))

	return lipgloss.JoinVertical(lipgloss.Top, header, content)
}

func generateRandomLine(words []string) []string {
	totalChars := 0
	selectedWords := []string{}

	for totalChars < 80 {
		selectedWord := words[rand.Intn(len(words))]
		totalChars += len(selectedWord)
		selectedWords = append(selectedWords, selectedWord)
	}

	return selectedWords
}

func (m model) Init() tea.Cmd {
	return m.timer.Init()
}

func main() {
	file, err := os.ReadFile("words.json")
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	var words []string
	json.Unmarshal(file, &words)

	fmt.Printf("words: %v", words)

	p := tea.NewProgram(initialModel(words))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
