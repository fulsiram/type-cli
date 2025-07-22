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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case timer.TickMsg, timer.StartStopMsg:
		m.timer, cmd = m.timer.Update(msg)
		cmds = append(cmds, cmd)

		if m.Exercise.Running() && m.timer.Timedout() {
			m.Exercise.Finish()
			cmds = append(cmds, m.timer.Stop())
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.nextWord):
			m.Exercise.Space()

		case key.Matches(msg, m.keymap.restart):
			cmds = append(cmds, m.timer.Stop())
			m.timer = timer.New(m.duration)
			m.Exercise.Reset()

		case key.Matches(msg, m.keymap.backSpace):
			m.Exercise.BackSpace()

		default:
			if m.Exercise.Pending() {
				cmds = append(cmds, m.timer.Start())
				m.Exercise.Start()
			}

			if len(msg.Runes) == 0 {
				break
			}

			pKey := msg.Runes[0]

			if !unicode.IsLetter(pKey) && !unicode.IsNumber(pKey) &&
				!unicode.IsPunct(pKey) && !unicode.IsSymbol(pKey) {
				break
			}

			m.Exercise.TypeLetter(msg.String())
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.cursor, cmd = m.cursor.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	result := m.Exercise.Result()

	header := lipgloss.NewStyle().
		MarginTop(1).
		Width(m.width).
		Align(lipgloss.Center).
		Render(
			"Type CLI\n",
			fmt.Sprintf("%s", m.timer.View()),
			fmt.Sprintf("%.2f wpm", m.statsCalc.RawWpm(result)),
			fmt.Sprintf("%.0f%%", m.statsCalc.Accuracy(result)*100),
		)

	contentStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-5).
		Align(lipgloss.Center, lipgloss.Center)

	content := ""
	if !m.Exercise.Finished() {
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

	footer := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		MarginBottom(1).
		Render(m.help.ShortHelpView([]key.Binding{
			m.keymap.restart,
		}))

	// return lipgloss.JoinVertical(lipgloss.Top, header, content)
	return lipgloss.JoinVertical(lipgloss.Top, header, content, footer)
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(cursor.Blink)
}
