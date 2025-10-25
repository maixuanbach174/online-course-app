package exercise

import "github.com/pkg/errors"

type Exercise struct {
	id            string
	lessonID      string
	question      string
	answers       []string
	correctAnswer string
	order         int
}

func NewExercise(id string, lessonID string, question string, answers []string, correctAnswer string, order int) (*Exercise, error) {
	if id == "" {
		return nil, errors.New("exercise id is required")
	}
	if lessonID == "" {
		return nil, errors.New("lesson id is required")
	}
	if question == "" {
		return nil, errors.New("question is required")
	}
	if len(answers) < 2 {
		return nil, errors.New("at least 2 answers are required")
	}
	if correctAnswer == "" {
		return nil, errors.New("correct answer is required")
	}

	// Validate correct answer exists in answers
	found := false
	for _, ans := range answers {
		if ans == correctAnswer {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("correct answer must be one of the provided answers")
	}

	return &Exercise{
		id:            id,
		lessonID:      lessonID,
		question:      question,
		answers:       answers,
		correctAnswer: correctAnswer,
		order:         order,
	}, nil
}

// Getters (read-only access for serialization/display)
func (e *Exercise) ID() string        { return e.id }
func (e *Exercise) LessonID() string  { return e.lessonID }
func (e *Exercise) Question() string  { return e.question }
func (e *Exercise) Answers() []string { return e.answers }
func (e *Exercise) Order() int        { return e.order }

// Note: correctAnswer is NOT exposed via public getter for security
// Students should not be able to see the correct answer directly

// CorrectAnswerForStorage returns the correct answer for persistence layer only
// This should only be used by the repository layer, not exposed to API/service layers
func (e *Exercise) CorrectAnswerForStorage() string {
	return e.correctAnswer
}

// Behavior methods
func (e *Exercise) CheckAnswer(answer string) bool {
	return e.correctAnswer == answer
}

func (e *Exercise) UpdateQuestion(question string) error {
	if question == "" {
		return errors.New("question is required")
	}
	e.question = question
	return nil
}

func (e *Exercise) UpdateAnswers(answers []string, correctAnswer string) error {
	if len(answers) < 2 {
		return errors.New("at least 2 answers are required")
	}
	if correctAnswer == "" {
		return errors.New("correct answer is required")
	}

	// Validate correct answer exists in answers
	found := false
	for _, ans := range answers {
		if ans == correctAnswer {
			found = true
			break
		}
	}
	if !found {
		return errors.New("correct answer must be one of the provided answers")
	}

	e.answers = answers
	e.correctAnswer = correctAnswer
	return nil
}

func (e *Exercise) UpdateOrder(order int) error {
	if order < 0 {
		return errors.New("order cannot be negative")
	}
	e.order = order
	return nil
}
