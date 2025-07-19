package app

import (
	"math/rand"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func generateWords(wordlist []string, n int) []string {
	var words []string

	for range n {
		words = append(words, wordlist[rand.Intn(len(wordlist))])
	}

	return words
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
