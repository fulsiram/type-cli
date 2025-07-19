package stats

import "github.com/fulsiram/type-cli/internal/exercise"

var CHARS_IN_WORD = 5

type Calculator struct {
}

func NewCalculator() Calculator {
	return Calculator{}
}

func (c Calculator) RawWpm(result exercise.Result) float64 {
	duration := result.Duration.Seconds()
	return float64(result.CharsTyped) / float64(CHARS_IN_WORD) / duration * 60
}

func (c Calculator) Accuracy(result exercise.Result) float64 {
	total := result.CharsCorrect + result.CharsIncorrect
	return float64(result.CharsCorrect) / float64(total)
}
