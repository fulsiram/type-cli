package app

import (
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
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
	testFinishedAt      time.Time
	charactersTyped     int
	correctCharsTyped   int
	incorrectCharsTyped int
}

func NewModel(words []string) model {
	m := model{
		timer:  timer.NewWithInterval(10*time.Second, time.Second),
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
