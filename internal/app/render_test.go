package app_test

import (
	"testing"

	"github.com/fulsiram/type-cli/internal/app"
)

func generateModel() app.Model {
	words := make([]string, 500)
	typedWords := make([]string, 500)
	for i := range words {
		words[i] = "benchmark"
		typedWords[i] = "aaaaaaaaaaaa"
	}

	m := app.NewModel(words)
	m.Exercise.TypedWords = typedWords
	return m
}

func BenchmarkRenderLines(b *testing.B) {
	m := generateModel()

	b.ResetTimer()
	for b.Loop() {
		_ = m.RenderLines()
	}
}

func BenchmarkRenderWord(b *testing.B) {
	m := generateModel()

	b.ResetTimer()
	for b.Loop() {
		_ = m.RenderWord(1)
	}
}

func BenchmarkView(b *testing.B) {
	m := generateModel()

	b.ResetTimer()
	for b.Loop() {
		m.View()
	}
}
