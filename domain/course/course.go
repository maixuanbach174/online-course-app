package course

import "github.com/pkg/errors"

type Course struct {
	id          string
	title       string
	description string
	thumbnail   string
	duration    int
	domain      Domain
	tags        []Tag
	rating      float64
	level       CourseLevel
	modules     []Module
}

func NewCourse(
	id string,
	title string,
	description string,
	thumbnail string,
	duration int,
	domain Domain,
	tags []Tag,
	rating float64,
	level CourseLevel,
	modules []Module,
) (*Course, error) {

	// a new course must have at least one module
	if len(modules) == 0 {
		return nil, errors.New("a course must have at least one module")
	}

	return &Course{
		id:          id,
		title:       title,
		description: description,
		thumbnail:   thumbnail,
		duration:    duration,
		domain:      domain,
		tags:        tags,
		rating:      rating,
		level:       level,
		modules:     modules,
	}, nil

}

func (c *Course) AddModule(module Module) error {
	panic("not implemented")
}

func (c *Course) RemoveModule(module Module) error {
	panic("not implemented")
}

func (c *Course) ReOrderModules() error {
	panic("not implemented")
}

func (c *Course) CanBePublished() bool {
	panic("not implemented")
}

func (c *Course) CalculateTotalDuration() int {
	panic("not implemented")
}
