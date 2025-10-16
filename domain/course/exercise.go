package course

type Exercise struct {
	id             string
	question       string
	answers        []string
	correct_answer string
	order          int
}

func NewExercise(id string, question string, answers []string, correct_answer string, order int) (*Exercise, error) {
	return &Exercise{
		id:             id,
		question:       question,
		answers:        answers,
		correct_answer: correct_answer,
		order:          order,
	}, nil
}
