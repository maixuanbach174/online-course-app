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

// Getters (read-only access for serialization/display)
func (c *Course) ID() string          { return c.id }
func (c *Course) Title() string       { return c.title }
func (c *Course) Description() string { return c.description }
func (c *Course) Thumbnail() string   { return c.thumbnail }
func (c *Course) Duration() int       { return c.duration }
func (c *Course) Domain() Domain      { return c.domain }
func (c *Course) Tags() []Tag         { return c.tags }
func (c *Course) Rating() float64     { return c.rating }
func (c *Course) Level() CourseLevel  { return c.level }
func (c *Course) Modules() []Module   { return c.modules }

// Behavior methods
func (c *Course) AddModule(module Module) error {
	c.modules = append(c.modules, module)
	return nil
}

func (c *Course) RemoveModule(moduleID string) error {
	for i, m := range c.modules {
		if m.ID() == moduleID {
			c.modules = append(c.modules[:i], c.modules[i+1:]...)
			return nil
		}
	}
	return errors.Errorf("module with id '%s' not found", moduleID)
}

func (c *Course) ReOrderModules(orderedModuleIDs []string) error {
	if len(orderedModuleIDs) != len(c.modules) {
		return errors.New("module count mismatch")
	}

	reordered := make([]Module, 0, len(orderedModuleIDs))
	for order, moduleID := range orderedModuleIDs {
		found := false
		for _, m := range c.modules {
			if m.ID() == moduleID {
				m.UpdateOrder(order)
				reordered = append(reordered, m)
				found = true
				break
			}
		}
		if !found {
			return errors.Errorf("module with id '%s' not found", moduleID)
		}
	}

	c.modules = reordered
	return nil
}

func (c *Course) CanBePublished() bool {
	// A course can be published if it has basic info and at least one module
	return c.title != "" &&
		c.description != "" &&
		len(c.modules) > 0
}

func (c *Course) CalculateTotalDuration() int {
	total := 0
	for _, module := range c.modules {
		total += module.CalculateDuration()
	}
	return total
}

// Additional behavior methods
func (c *Course) HasTag(tag Tag) bool {
	for _, t := range c.tags {
		if t.String() == tag.String() {
			return true
		}
	}
	return false
}

func (c *Course) GetModuleByID(moduleID string) (*Module, error) {
	for _, m := range c.modules {
		if m.ID() == moduleID {
			return &m, nil
		}
	}
	return nil, errors.Errorf("module with id '%s' not found", moduleID)
}

func (c *Course) UpdateBasicInfo(title, description, thumbnail string) error {
	if title == "" {
		return errors.New("title is required")
	}
	c.title = title
	c.description = description
	c.thumbnail = thumbnail
	return nil
}
