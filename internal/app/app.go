package app

import (
	"fmt"
	"time"

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

		if m.testStarted && m.timer.Timedout() {
			m.testFinishedAt = time.Now()
			m.testStarted = false
		}
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
			if m.timer.Timedout() {
				break
			}

			if !m.testStarted {
				m.testStarted = true
				m.testStartedAt = time.Now()
				m.charactersTyped = 0
				m.incorrectCharsTyped = 0
				m.correctCharsTyped = 0
				cmds = append(cmds, m.timer.Start())
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

	m.cursor, cmd = m.cursor.Update(msg)
	cmds = append(cmds, cmd)

	m.renderedWords[m.currentWord] = m.renderCurrentWord()
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
			fmt.Sprintf("%.2f wpm", float64(m.charactersTyped)/time.Since(m.testStartedAt).Minutes()/5),
			fmt.Sprintf("%.0f", float32(m.correctCharsTyped)/float32(max(m.correctCharsTyped+m.incorrectCharsTyped, 1))*100),
		)

	contentStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-2).
		Align(lipgloss.Center, lipgloss.Center)

	content := ""
	if !m.timer.Timedout() {
		content = contentStyle.Render(
			lipgloss.NewStyle().
				Width(60).
				Align(lipgloss.Left).
				Render(m.renderLines()),
		)
	} else {
		content = contentStyle.Render(
			fmt.Sprintf("%.2f wpm\n", float64(m.charactersTyped)/m.testFinishedAt.Sub(m.testStartedAt).Minutes()/5),
			fmt.Sprintf("%.0f%% accuracy\n", float32(m.correctCharsTyped)/float32(max(m.correctCharsTyped+m.incorrectCharsTyped, 1))*100),
			fmt.Sprintf("%.0f sec", m.testFinishedAt.Sub(m.testStartedAt).Seconds()),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Top, header, content)
}

func (m model) Init() tea.Cmd {
	return tea.Batch(cursor.Blink)
}
