package app

import (
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/lipgloss"
	"github.com/fulsiram/type-cli/internal/exercise"
	"github.com/fulsiram/type-cli/internal/stats"
)

type keymap struct {
	quit      key.Binding
	nextWord  key.Binding
	backSpace key.Binding
	restart   key.Binding
}

type Model struct {
	timer    timer.Model
	cursor   cursor.Model
	keymap   keymap
	quitting bool

	width  int
	height int

	ExerciseService exercise.Service
	statsCalc       stats.Calculator
}

func NewModel(words []string) Model {
	m := Model{
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
			restart: key.NewBinding(key.WithKeys("enter")),
			backSpace: key.NewBinding(
				key.WithKeys("backspace"),
			),
		},

		ExerciseService: exercise.NewService(words),
		statsCalc:       stats.NewCalculator(),
	}

	m.cursor.SetMode(cursor.CursorStatic)
	m.cursor.Style = lipgloss.NewStyle().
		Background(lipgloss.Color("#000000")).
		Foreground(lipgloss.Color("#FFFFFF"))

	m.cursor.Focus()

	return m
}
