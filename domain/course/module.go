package course

import "github.com/pkg/errors"

type Module struct {
	id      string
	title   string
	lessons []Lesson
	order   int
}

func NewModule(id string, title string, lessons []Lesson, order int) (*Module, error) {

	// a module must have at least one lesson
	if len(lessons) == 0 {
		return nil, errors.New("a module must have at least one lesson")
	}

	return &Module{
		id:      id,
		title:   title,
		lessons: lessons,
		order:   order,
	}, nil
}

func (m *Module) AddLesson(lesson Lesson) error {
	panic("not implemented")
}

func (m *Module) RemoveLesson(lesson Lesson) error {
	panic("not implemented")
}

func (m *Module) ReOrderLessons() error {
	panic("not implemented")
}
