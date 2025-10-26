package module

import "github.com/pkg/errors"

type Module struct {
	id       string
	courseID string
	title    string
	order    int
}

func NewModule(id string, courseID string, title string, order int) (*Module, error) {
	if id == "" {
		return nil, errors.New("module id is required")
	}
	if courseID == "" {
		return nil, errors.New("course id is required")
	}
	if title == "" {
		return nil, errors.New("module title is required")
	}

	return &Module{
		id:       id,
		courseID: courseID,
		title:    title,
		order:    order,
	}, nil
}

// Getters (read-only access for serialization/display)
func (m *Module) ID() string       { return m.id }
func (m *Module) CourseID() string { return m.courseID }
func (m *Module) Title() string    { return m.title }
func (m *Module) Order() int       { return m.order }

// Behavior methods
func (m *Module) UpdateOrder(order int) error {
	if order < 0 {
		return errors.New("order cannot be negative")
	}
	m.order = order
	return nil
}

func (m *Module) UpdateTitle(title string) error {
	if title == "" {
		return errors.New("title is required")
	}
	m.title = title
	return nil
}
