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

	wordlist      []string
	words         []string
	typedWords    []string
	currentWord   int
	currentCharId int

	renderedWords []string
}

// func generateLines(words []string, n int) [][]string {
// 	var lines [][]string
//
// 	for range n {
// 		lines = append(lines, generateRandomLine(words))
// 	}
//
// 	return lines
// }

func generateWords(wordlist []string, n int) []string {
	var words []string

	for range n {
		words = append(words, wordlist[rand.Intn(len(wordlist))])
	}

	return words
}

func initialModel(words []string) model {
	m := model{
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
		wordlist:   words,
		words:      generateWords(words, 100),
		typedWords: []string{""},
	}

	m.renderedWords = m.renderWords()

	return m
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

			m.typedWords = append(m.typedWords, "")
			m.currentWord += 1
			m.currentCharId = 0

		case key.Matches(msg, m.keymap.backSpace):
			if m.currentCharId == 0 {
				if m.currentWord > 0 {
					m.typedWords = m.typedWords[:m.currentWord]
					m.currentWord -= 1
					m.currentCharId = len(m.typedWords[m.currentWord])
				}
				return m, nil
			}

			m.currentCharId -= 1
			m.typedWords[m.currentWord] = m.typedWords[m.currentWord][:m.currentCharId]
			m.renderedWords[m.currentWord] = m.renderWord(m.typedWords[m.currentWord], m.words[m.currentWord])

		default:
			m.currentCharId += 1
			m.typedWords[m.currentWord] += msg.String()
			m.renderedWords[m.currentWord] = m.renderWord(m.typedWords[m.currentWord], m.words[m.currentWord])
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

	content := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-2).
		Align(lipgloss.Center, lipgloss.Center).
		Render(m.renderLines())

	return lipgloss.JoinVertical(lipgloss.Top, header, content)
}

func (m model) renderLines() string {
	// untypedStyle := lipgloss.NewStyle().
	// 	Foreground(lipgloss.Color("#646669"))
	//
	lines := m.getRenderedLines()
	currentLine := m.getCurrentLine()

	// fmt.Printf("%v\n%v", lines, currentLine)

	var shownLines [][]string
	if currentLine == 0 {
		shownLines = lines[:3]
	} else if currentLine >= len(lines)-3 {
		shownLines = lines[len(lines)-3:]
	} else {
		shownLines = lines[currentLine-1 : currentLine+2]
	}

	var renderedLines []string
	for _, line := range shownLines {
		renderedLines = append(renderedLines, strings.Join(line, " "))
	}

	return strings.Join(renderedLines, "\n")
}

func (m model) renderWord(typedWord string, fullWord string) string {
	untypedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#646669"))

	correctLetterStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	incorrectLetterStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#DB4B4C"))

	renderedWord := ""
	for i := range typedWord {
		if i >= len(fullWord) {
			renderedWord += incorrectLetterStyle.Render(string(typedWord[i]))
		} else if fullWord[i] == typedWord[i] {
			renderedWord += correctLetterStyle.Render(string(typedWord[i]))
		} else {
			renderedWord += incorrectLetterStyle.Render(string(typedWord[i]))
		}
	}

	if len(typedWord) < len(fullWord) {
		renderedWord += untypedStyle.Render(fullWord[len(typedWord):])
	}

	return renderedWord
}

func (m model) renderWords() []string {
	var renderedWords []string

	for i, word := range m.words {
		if i < len(m.typedWords) {
			typedWord := m.typedWords[i]

			renderedWords = append(renderedWords, m.renderWord(typedWord, word))
		} else {
			renderedWords = append(renderedWords, m.renderWord("", word))
		}
	}

	return renderedWords
}

func (m model) getRenderedLines() [][]string {
	charCount := 0
	var lines [][]string
	var words []string

	for i, word := range m.words {
		// fmt.Printf("%s %d %d %d\n\n", word, len(word), charCount, len(words))
		if charCount+len(word) > 60 {
			lines = append(lines, words)
			charCount = 0
			words = []string{}
		}

		charCount += len(word)
		words = append(words, m.renderedWords[i])
	}

	if len(words) > 0 {
		lines = append(lines, words)
	}

	return lines
}

func (m model) getCurrentLine() int {
	charCount := 0

	for i, word := range m.words {
		charCount += len(word)

		if i == m.currentWord {
			return charCount / 60
		}
	}
	panic("shouldn't be reachable")
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

	// fmt.Printf("words: %v", words)

	p := tea.NewProgram(initialModel(words))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
