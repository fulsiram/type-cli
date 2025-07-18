package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
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
	cursor   cursor.Model
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

	testStarted         bool
	testStartedAt       time.Time
	charactersTyped     int
	correctCharsTyped   int
	incorrectCharsTyped int
}

func generateWords(wordlist []string, n int) []string {
	var words []string

	for range n {
		words = append(words, wordlist[rand.Intn(len(wordlist))])
	}

	return words
}

func initialModel(words []string) model {
	m := model{
		timer:  timer.NewWithInterval(30*time.Second, time.Second),
		cursor: cursor.New(),

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

		wordlist: words,
		words:    generateWords(words, 100),

		typedWords:  make([]string, 100),
		testStarted: false,
	}

	m.cursor.SetChar("a")
	m.cursor.SetMode(cursor.CursorStatic)
	m.cursor.Style = lipgloss.NewStyle().
		Background(lipgloss.Color("#000000")).
		Foreground(lipgloss.Color("#FFFFFF"))
	// Bold(true)

	// m.cursor.Style = lipgloss.NewStyle().
	// 	Border(lipgloss.BlockBorder())

	m.cursor.Focus()

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

			m.renderedWords[m.currentWord] = m.renderWord(m.typedWords[m.currentWord], m.words[m.currentWord])
			m.currentWord += 1
			m.currentCharId = 0
			m.charactersTyped++

		case key.Matches(msg, m.keymap.backSpace):
			if m.currentCharId == 0 {
				if m.currentWord > 0 {
					// m.typedWords = m.typedWords[:m.currentWord]
					m.renderedWords[m.currentWord] = m.renderWord("", m.words[m.currentWord])
					m.currentWord -= 1
					m.currentCharId = len(m.typedWords[m.currentWord])
				}
				return m, nil
			}

			m.currentCharId -= 1
			m.typedWords[m.currentWord] = m.typedWords[m.currentWord][:m.currentCharId]
			// m.renderedWords[m.currentWord] = m.renderWord(m.typedWords[m.currentWord], m.words[m.currentWord])
			m.renderedWords[m.currentWord] = m.renderCurrentWord()
		default:
			if !m.testStarted {
				m.testStarted = true
				m.testStartedAt = time.Now()
			}

			if len(m.typedWords[m.currentWord]) > len(m.words[m.currentWord])+15 {
				break
			}

			m.charactersTyped++
			if len(m.typedWords[m.currentWord]) >= len(m.words[m.currentWord]) {
				m.incorrectCharsTyped++
			} else if msg.String() != string(m.words[m.currentWord][m.currentCharId]) {
				m.incorrectCharsTyped++
			} else {
				m.correctCharsTyped++
			}

			m.currentCharId += 1
			m.typedWords[m.currentWord] += msg.String()
			// m.renderedWords[m.currentWord] = m.renderWord(m.typedWords[m.currentWord], m.words[m.currentWord])
			m.renderedWords[m.currentWord] = m.renderCurrentWord()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmd tea.Cmd
	m.cursor, cmd = m.cursor.Update(msg)

	m.renderedWords[m.currentWord] = m.renderCurrentWord()
	return m, cmd
}

func (m model) View() string {
	header := lipgloss.NewStyle().
		MarginTop(1).
		Width(m.width).
		Align(lipgloss.Center).
		Render(
			"Type CLI",
			fmt.Sprintf("%.2f wpm", float64(m.charactersTyped)/time.Since(m.testStartedAt).Minutes()/5),
			fmt.Sprintf("%.0f", float32(m.correctCharsTyped)/float32(max(m.correctCharsTyped+m.incorrectCharsTyped, 1))*100),
		)

	content := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-2).
		Align(lipgloss.Center, lipgloss.Center).
		Render(
			lipgloss.NewStyle().
				Width(60).
				Align(lipgloss.Left).
				Render(m.renderLines()),
		)

	return lipgloss.JoinVertical(lipgloss.Top, header, content)
}

func (m model) renderLines() string {
	var lines []string
	line := ""
	lineLenght := 0
	currentLine := 0
	for i, word := range m.words {
		// TODO: Account for spaces
		wordLength := max(len(word), len(m.typedWords[i]))

		if wordLength+lineLenght+1 > 60 {
			lines = append(lines, line)
			line = ""
			lineLenght = 0
		}

		lineLenght += wordLength + 1

		if i == m.currentWord {
			currentLine = len(lines)
			m.renderedWords[i] = m.renderCurrentWord()
		}

		line += m.renderedWords[i]

		if i == m.currentWord && len(m.typedWords[i]) >= len(word) {
			m.cursor.SetChar(" ")
			line += m.cursor.View()
		} else {
			line += " "
		}
	}

	var shownLines []string
	if currentLine == 0 {
		shownLines = lines[:3]
	} else if currentLine >= len(lines)-3 {
		shownLines = lines[len(lines)-3:]
	} else {
		shownLines = lines[currentLine-1 : currentLine+2]
	}

	return strings.Join(shownLines, "\n")
}

func (m model) renderCurrentWord() string {
	untypedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#646669"))

	correctLetterStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	incorrectLetterStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#DB4B4C"))

	currentWord := m.words[m.currentWord]
	typedWord := m.typedWords[m.currentWord]

	renderedWord := ""
	for i := range typedWord {
		if i >= len(currentWord) {
			renderedWord += incorrectLetterStyle.Render(string(typedWord[i]))
		} else if currentWord[i] == typedWord[i] {
			renderedWord += correctLetterStyle.Render(string(typedWord[i]))
		} else {
			renderedWord += incorrectLetterStyle.Render(string(typedWord[i]))
		}
	}

	if len(typedWord) < len(currentWord) {
		m.cursor.SetChar(string(currentWord[m.currentCharId]))
		renderedWord += untypedStyle.Render(m.cursor.View())
		renderedWord += untypedStyle.Render(currentWord[len(typedWord)+1:])
	}

	return renderedWord
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

func (m model) Init() tea.Cmd {
	return tea.Batch(m.timer.Init(), cursor.Blink)
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
