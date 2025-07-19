package app

import (
	"fmt"
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

		if m.exerciseService.Running() && m.timer.Timedout() {
			m.exerciseService.Finish()
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.nextWord):
			m.exerciseService.Space()

		case key.Matches(msg, m.keymap.backSpace):
			m.exerciseService.BackSpace()

		default:
			if m.exerciseService.Pending() {
				cmds = append(cmds, m.timer.Start())
				m.exerciseService.Start()
			}

			if len(msg.Runes) == 0 {
				break
			}

			pKey := msg.Runes[0]

			if !unicode.IsLetter(pKey) && !unicode.IsNumber(pKey) &&
				!unicode.IsPunct(pKey) && !unicode.IsSymbol(pKey) {
				break
			}

			m.exerciseService.TypeLetter(msg.String())
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
	header := lipgloss.NewStyle().
		MarginTop(1).
		Width(m.width).
		Align(lipgloss.Center).
		Render(
			"Type CLI",
			fmt.Sprintf("%s", m.timer.View()),
			// fmt.Sprintf("%.2f wpm", float64(m.charactersTyped)/time.Since(m.testStartedAt).Minutes()/5),
			// fmt.Sprintf("%.0f", float32(m.correctCharsTyped)/float32(max(m.correctCharsTyped+m.incorrectCharsTyped, 1))*100),
		)

	contentStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-2).
		Align(lipgloss.Center, lipgloss.Center)

	content := ""
	if !m.exerciseService.Finished() {
		content = contentStyle.Render(
			lipgloss.NewStyle().
				Width(60).
				Align(lipgloss.Left).
				Render(m.renderLines()),
		)
	} else {
		content = contentStyle.Render(
			"stats",
			// fmt.Sprintf("%.2f wpm\n", float64(m.charactersTyped)/m.testFinishedAt.Sub(m.testStartedAt).Minutes()/5),
			// fmt.Sprintf("%.0f%% accuracy\n", float32(m.correctCharsTyped)/float32(max(m.correctCharsTyped+m.incorrectCharsTyped, 1))*100),
			// fmt.Sprintf("%.0f sec", m.testFinishedAt.Sub(m.testStartedAt).Seconds()),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Top, header, content)
}

func (m model) Init() tea.Cmd {
	return tea.Batch(cursor.Blink)
}
