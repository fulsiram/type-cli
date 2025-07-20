package exercise

import (
	"time"
)

type State int

const (
	Pending State = iota
	Running
	Finished
)

type Service struct {
	eg exerciseGenerator

	Words      []string
	TypedWords []string

	wordIdx int

	state      State
	startedAt  time.Time
	finishedAt time.Time

	typed     int
	correct   int
	incorrect int
}

type Result struct {
	CharsTyped     int
	CharsCorrect   int
	CharsIncorrect int

	StartedAt  time.Time
	FinishedAt time.Time
	Duration   time.Duration
}

func NewService(words []string) Service {
	service := Service{
		eg:    NewExerciseGenerator(words),
		state: State(Pending),
	}
	service.Reset()

	return service
}

func (s *Service) Reset() {
	s.state = State(Pending)
	s.Words = s.eg.Generate(50)
	s.TypedWords = make([]string, 50)
	s.wordIdx = 0
	s.correct = 0
	s.incorrect = 0
	s.typed = 0
}

func (s *Service) Start() {
	s.state = State(Running)
	s.startedAt = time.Now()
}

func (s *Service) Finish() {
	s.state = State(Finished)
	s.finishedAt = time.Now()
}

func (s Service) Result() Result {
	endTime := time.Now()
	if s.Finished() {
		endTime = s.finishedAt
	}

	return Result{
		CharsTyped:     s.typed,
		CharsCorrect:   s.correct,
		CharsIncorrect: s.incorrect,

		// StartedAt:  s.startedAt,
		// FinishedAt: s.finishedAt,
		Duration: endTime.Sub(s.startedAt),
	}
}

func (s Service) State() State {
	return s.state
}

func (s Service) Pending() bool {
	return s.state == Pending
}

func (s Service) Running() bool {
	return s.state == Running
}

func (s Service) Finished() bool {
	return s.state == Finished
}

func (s Service) CurrentWord() string {
	return s.Words[s.wordIdx]
}

func (s Service) CurrentTypedWord() string {
	return s.TypedWords[s.wordIdx]
}

func (s Service) NextLetter() string {
	return string(s.CurrentWord()[len(s.CurrentTypedWord())])
}

func (s Service) WordIdx() int {
	return s.wordIdx
}

func (s Service) Word(idx int) string {
	return s.Words[idx]
}

func (s Service) TypedWord(idx int) string {
	return s.TypedWords[idx]
}

func (s Service) IsCurrentWord(idx int) bool {
	return s.wordIdx == idx
}

func (s *Service) Space() {
	if s.state != Running {
		return
	}

	if len(s.CurrentTypedWord()) == 0 {
		return
	}

	if len(s.CurrentTypedWord()) < len(s.CurrentWord()) {
		s.incorrect++
	}

	s.wordIdx++
}

func (s *Service) BackSpace() {
	if s.state != Running {
		return
	}

	curWord := s.CurrentTypedWord()
	if s.wordIdx == 0 && len(curWord) == 0 {
		return
	}

	if len(curWord) > 0 {
		s.TypedWords[s.wordIdx] = curWord[:len(curWord)-1]
	} else {
		s.wordIdx--
	}
}

func (s *Service) TypeLetter(letter string) {
	if s.state != Running {
		return
	}

	curWord, curTypedWord := s.CurrentWord(), s.CurrentTypedWord()

	s.typed++

	if len(curTypedWord) > len(curWord)+15 {
		s.incorrect++
		// Don't add a letter if curTypedWord is much longer
		return
	} else if len(curTypedWord) >= len(curWord) {
		s.incorrect++
	} else if letter != s.NextLetter() {
		s.incorrect++
	} else {
		s.correct++
	}

	s.TypedWords[s.wordIdx] += letter
}
