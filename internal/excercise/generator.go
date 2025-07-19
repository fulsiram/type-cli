package excercise

import "math/rand"

type excerciseGenarator struct {
	wordlist []string
}

func NewExcerciseGenerator(words []string) excerciseGenarator {
	return excerciseGenarator{
		wordlist: words,
	}
}

func (eg excerciseGenarator) Generate(length int) []string {
	var words []string

	for range length {
		words = append(words, eg.wordlist[rand.Intn(len(eg.wordlist))])
	}

	return words
}
