package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var untypedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#646669"))

var correctLetterStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFFFF"))

var incorrectLetterStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#DB4B4C"))

func (m Model) RenderLines() string {
	var lines [][]int
	var line []int
	lineLength := 0
	currentLine := 0

	for i, word := range m.ExerciseService.Words {
		wordLength := max(len(word), len(m.ExerciseService.TypedWord(i)))

		if wordLength+lineLength+1 > 60 {
			lines = append(lines, line)
			line = make([]int, 0)
			lineLength = 0

			if len(line) > currentLine+3 {
				break
			}
		}

		lineLength += wordLength + 1

		if m.ExerciseService.IsCurrentWord(i) {
			currentLine = len(lines)
		}

		line = append(line, i)
	}

	if lineLength > 0 {
		lines = append(lines, line)
	}

	start := max(0, min(currentLine-1, len(lines)-3))
	end := min(len(lines), start+3)
	shownLines := lines[start:end]

	var renderedLines []string
	for _, line := range shownLines {
		var sb strings.Builder
		for _, wordIdx := range line {
			sb.WriteString(m.RenderWord(wordIdx))

			if !m.ExerciseService.IsCurrentWord(wordIdx) {
				sb.WriteString(" ")
				continue
			}

			curWord := m.ExerciseService.CurrentWord()
			curTypedWord := m.ExerciseService.CurrentTypedWord()

			if len(curTypedWord) >= len(curWord) {
				m.cursor.SetChar(" ")
				sb.WriteString(m.cursor.View())
			} else {
				sb.WriteString(" ")
			}
		}
		renderedLines = append(renderedLines, sb.String())
	}

	return strings.Join(renderedLines, "\n")

}

func (m Model) RenderWord(idx int) string {
	currentWord := m.ExerciseService.Word(idx)
	typedWord := m.ExerciseService.TypedWord(idx)

	var sb strings.Builder

	lastCharCorrect := false
	var renderBuf strings.Builder
	for i := range typedWord {
		if i >= len(currentWord) || currentWord[i] != typedWord[i] {
			if renderBuf.Len() > 0 && lastCharCorrect {
				sb.WriteString(correctLetterStyle.Render(renderBuf.String()))
				renderBuf.Reset()
			}
			lastCharCorrect = false
		} else if currentWord[i] == typedWord[i] {
			if renderBuf.Len() > 0 && !lastCharCorrect {
				sb.WriteString(incorrectLetterStyle.Render(renderBuf.String()))
				renderBuf.Reset()
			}
			lastCharCorrect = true
		}
		renderBuf.WriteByte(typedWord[i])
	}

	if lastCharCorrect {
		sb.WriteString(correctLetterStyle.Render(renderBuf.String()))
	} else {
		sb.WriteString(incorrectLetterStyle.Render(renderBuf.String()))
	}

	if len(typedWord) < len(currentWord) {
		if m.ExerciseService.IsCurrentWord(idx) {
			m.cursor.SetChar(m.ExerciseService.NextLetter())
			sb.WriteString(m.cursor.View())
			sb.WriteString(untypedStyle.Render(currentWord[len(typedWord)+1:]))
		} else {
			sb.WriteString(untypedStyle.Render(currentWord[len(typedWord):]))
		}
	}

	return sb.String()
}
