package excercise

import (
	"time"
)

type Service struct {
	eg excerciseGenarator

	Words      []string
	TypedWords []string

	wordIdx int

	running    bool
	startedAt  time.Time
	finishedAt time.Time

	typed     int
	correct   int
	incorrect int
}

func NewService(words []string) Service {
	service := Service{
		eg:      excerciseGenarator{},
		running: false,
	}
	service.Reset()

	return service
}

func (s *Service) Reset() {
	s.running = false
	s.Words = s.eg.Generate(500)
	s.TypedWords = make([]string, 500)
}

func (s *Service) Start() {
	s.running = true
	s.startedAt = time.Now()
}

func (s *Service) Finish() {
	s.running = false
	s.finishedAt = time.Now()
}

func (s *Service) Running() bool {
	return s.running
}

func (s Service) CurrentWord() string {
	return s.Words[s.wordIdx]
}

func (s Service) CurrentTypedWord() string {
	return s.Words[s.wordIdx]
}

func (s Service) nextLetter() string {
	return string(s.CurrentWord()[len(s.CurrentTypedWord())])
}

func (s *Service) Space() {
	if !s.Running() {
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
	if !s.Running() {
		return
	}

	curWord := s.CurrentTypedWord()
	if s.wordIdx == 0 && len(curWord) == 0 {
		return
	}

	if len(curWord) > 0 {
		s.Words[s.wordIdx] = curWord[:len(curWord)-1]
	} else {
		s.wordIdx--
	}
}

func (s *Service) TypeLetter(letter string) {
	if !s.Running() {
		return
	}

	curWord, curTypedWord := s.CurrentWord(), s.CurrentTypedWord()

	s.typed++

	if len(curTypedWord) > len(curWord)+15 {
		s.incorrect++
		// Don't add a letter if curTypedWord is much longer
		return
	} else if len(curWord) >= len(curTypedWord) {
		s.incorrect++
	} else if letter != s.nextLetter() {
		s.incorrect++
	} else {
		s.correct++
	}

	s.TypedWords[s.wordIdx] += letter
}
