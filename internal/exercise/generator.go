package exercise

import "math/rand"

type exerciseGenerator struct {
	wordlist []string
}

func NewExerciseGenerator(words []string) exerciseGenerator {
	return exerciseGenerator{
		wordlist: words,
	}
}

func (eg exerciseGenerator) Generate(length int) []string {
	var words []string

	for range length {
		words = append(words, eg.wordlist[rand.Intn(len(eg.wordlist))])
	}

	return words
}
