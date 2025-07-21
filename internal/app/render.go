package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) RenderLines() string {
	var lines [][]int
	var line []int
	lineLength := 0
	currentLine := -1
	for i, word := range m.Exercise.Words {
		wordLength := max(len(word), len(m.Exercise.TypedWord(i)))

		if wordLength+lineLength+1 > 60 {
			lines = append(lines, line)
			line = make([]int, 0)
			lineLength = 0

			if currentLine >= 0 && len(lines)-3 >= currentLine {
				break
			}
		}

		lineLength += wordLength + 1

		if m.Exercise.IsCurrentWord(i) {
			currentLine = len(lines)
		}

		line = append(line, i)
	}

	if len(line) > 0 {
		lines = append(lines, line)
	}

	var shownLines [][]int
	if len(lines) <= 3 {
		shownLines = lines
	} else if currentLine == 0 {
		shownLines = lines[:3]
	} else if currentLine >= len(lines)-2 {
		shownLines = lines[len(lines)-3:]
	} else {
		shownLines = lines[currentLine-1 : currentLine+2]
	}

	var renderedLines []string
	for _, line := range shownLines {
		renderedLine := ""
		for _, wordIdx := range line {
			renderedLine += m.RenderWord(wordIdx)

			if !m.Exercise.IsCurrentWord(wordIdx) {
				renderedLine += " "
				continue
			}

			curWord := m.Exercise.CurrentWord()
			curTypedWord := m.Exercise.CurrentTypedWord()

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

func (m Model) RenderWord(idx int) string {
	untypedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#646669"))

	correctLetterStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	incorrectLetterStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#DB4B4C"))

	currentWord := m.Exercise.Word(idx)
	typedWord := m.Exercise.TypedWord(idx)

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
		if m.Exercise.IsCurrentWord(idx) {
			m.cursor.SetChar(m.Exercise.NextLetter())
			renderedWord += m.cursor.View()
			renderedWord += untypedStyle.Render(currentWord[len(typedWord)+1:])
		} else {
			renderedWord += untypedStyle.Render(currentWord[len(typedWord):])
		}
	}

	return renderedWord
}
