package app

import (
	"fmt"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case timer.TickMsg, timer.StartStopMsg:
		m.timer, cmd = m.timer.Update(msg)
		cmds = append(cmds, cmd)

		if m.ExerciseService.Running() && m.timer.Timedout() {
			m.ExerciseService.Finish()
			cmds = append(cmds, m.timer.Stop())
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.nextWord):
			m.ExerciseService.Space()

		case key.Matches(msg, m.keymap.restart):
			cmds = append(cmds, m.timer.Stop())
			m.timer.Timeout = 10 * time.Second
			m.ExerciseService.Reset()

		case key.Matches(msg, m.keymap.backSpace):
			m.ExerciseService.BackSpace()

		default:
			if m.ExerciseService.Pending() {
				cmds = append(cmds, m.timer.Start())
				m.ExerciseService.Start()
			}

			if len(msg.Runes) == 0 {
				break
			}

			pKey := msg.Runes[0]

			if !unicode.IsLetter(pKey) && !unicode.IsNumber(pKey) &&
				!unicode.IsPunct(pKey) && !unicode.IsSymbol(pKey) {
				break
			}

			m.ExerciseService.TypeLetter(msg.String())
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.cursor, cmd = m.cursor.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	result := m.ExerciseService.Result()

	header := lipgloss.NewStyle().
		MarginTop(1).
		Width(m.width).
		Align(lipgloss.Center).
		Render(
			"Type CLI",
			fmt.Sprintf("%s", m.timer.View()),
			fmt.Sprintf("%.2f wpm", m.statsCalc.RawWpm(result)),
			fmt.Sprintf("%.0f", m.statsCalc.Accuracy(result)*100),
		)

	contentStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-2).
		Align(lipgloss.Center, lipgloss.Center)

	content := ""
	if !m.ExerciseService.Finished() {
		content = contentStyle.Render(
			lipgloss.NewStyle().
				Width(60).
				Align(lipgloss.Left).
				Render(m.RenderLines()),
		)
	} else {
		content = contentStyle.Render(
			"stats",
			fmt.Sprintf("%.2f wpm\n", m.statsCalc.RawWpm(result)),
			fmt.Sprintf("%.0f%% accuracy\n", m.statsCalc.Accuracy(result)*100),
			fmt.Sprintf("%.0f sec", result.Duration.Seconds()),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Top, header, content)
}

func (m model) Init() tea.Cmd {
	return tea.Batch(cursor.Blink)
}
