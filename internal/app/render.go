package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) renderLines() string {
	var lines [][]int
	var line []int
	lineLength := 0
	currentLine := 0
	for i, word := range m.exerciseService.Words {
		wordLength := max(len(word), len(m.exerciseService.TypedWord(i)))

		if wordLength+lineLength+1 > 60 {
			lines = append(lines, line)
			line = make([]int, 0)
			lineLength = 0
		}

		lineLength += wordLength + 1

		if m.exerciseService.IsCurrentWord(i) {
			currentLine = len(lines)
		}

		line = append(line, i)
	}

	var shownLines [][]int
	if currentLine == 0 {
		shownLines = lines[:3]
	} else if currentLine >= len(lines)-3 {
		shownLines = lines[len(lines)-3:]
	} else {
		shownLines = lines[currentLine-1 : currentLine+2]
	}

	var renderedLines []string
	for _, line := range shownLines {
		renderedLine := ""
		for _, wordIdx := range line {
			renderedLine += m.renderWord(wordIdx)

			if !m.exerciseService.IsCurrentWord(wordIdx) {
				renderedLine += " "
				continue
			}

			curWord := m.exerciseService.CurrentWord()
			curTypedWord := m.exerciseService.CurrentTypedWord()

			if len(curTypedWord) >= len(curWord) {
				m.cursor.SetChar(" ")
				renderedLine += m.cursor.View()
			} else {
				renderedLine += " "
			}
		}
		renderedLines = append(renderedLines, renderedLine)
	}

	return strings.Join(renderedLines, "\n")
}

func (m model) renderWord(idx int) string {
	untypedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#646669"))

	correctLetterStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	incorrectLetterStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#DB4B4C"))

	currentWord := m.exerciseService.Word(idx)
	typedWord := m.exerciseService.TypedWord(idx)

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
		if m.exerciseService.IsCurrentWord(idx) {
			m.cursor.SetChar(m.exerciseService.NextLetter())
			renderedWord += m.cursor.View()
			renderedWord += untypedStyle.Render(currentWord[len(typedWord)+1:])
		} else {
			renderedWord += untypedStyle.Render(currentWord[len(typedWord):])
		}
	}

	return renderedWord
}
