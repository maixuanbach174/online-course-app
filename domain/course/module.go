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

// Getters (read-only access for serialization/display)
func (m *Module) ID() string        { return m.id }
func (m *Module) Title() string     { return m.title }
func (m *Module) Lessons() []Lesson { return m.lessons }
func (m *Module) Order() int        { return m.order }

// Behavior methods
func (m *Module) AddLesson(lesson Lesson) error {
	m.lessons = append(m.lessons, lesson)
	return nil
}

func (m *Module) RemoveLesson(lessonID string) error {
	for i, l := range m.lessons {
		if l.ID() == lessonID {
			m.lessons = append(m.lessons[:i], m.lessons[i+1:]...)
			return nil
		}
	}
	return errors.Errorf("lesson with id '%s' not found", lessonID)
}

func (m *Module) ReOrderLessons(orderedLessonIDs []string) error {
	if len(orderedLessonIDs) != len(m.lessons) {
		return errors.New("lesson count mismatch")
	}

	reordered := make([]Lesson, 0, len(orderedLessonIDs))
	for order, lessonID := range orderedLessonIDs {
		found := false
		for _, l := range m.lessons {
			if l.ID() == lessonID {
				l.UpdateOrder(order)
				reordered = append(reordered, l)
				found = true
				break
			}
		}
		if !found {
			return errors.Errorf("lesson with id '%s' not found", lessonID)
		}
	}

	m.lessons = reordered
	return nil
}

func (m *Module) UpdateOrder(order int) {
	m.order = order
}

func (m *Module) UpdateTitle(title string) error {
	if title == "" {
		return errors.New("title is required")
	}
	m.title = title
	return nil
}

func (m *Module) CalculateDuration() int {
	total := 0
	for _, lesson := range m.lessons {
		total += lesson.CalculateDuration()
	}
	return total
}

func (m *Module) GetLessonByID(lessonID string) (*Lesson, error) {
	for _, l := range m.lessons {
		if l.ID() == lessonID {
			return &l, nil
		}
	}
	return nil, errors.Errorf("lesson with id '%s' not found", lessonID)
}
