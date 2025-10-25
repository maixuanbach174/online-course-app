package course

import "github.com/pkg/errors"

type Course struct {
	id          string
	teacherID   string
	title       string
	description string
	thumbnail   string
	duration    int
	domain      Domain
	tags        []Tag
	rating      float64
	level       CourseLevel
}

func NewCourse(
	id string,
	teacherID string,
	title string,
	description string,
	thumbnail string,
	duration int,
	domain Domain,
	tags []Tag,
	rating float64,
	level CourseLevel,
) (*Course, error) {

	// Validate required fields
	if id == "" {
		return nil, errors.New("course id is required")
	}
	if teacherID == "" {
		return nil, errors.New("teacher id is required")
	}
	if title == "" {
		return nil, errors.New("course title is required")
	}

	return &Course{
		id:          id,
		teacherID:   teacherID,
		title:       title,
		description: description,
		thumbnail:   thumbnail,
		duration:    duration,
		domain:      domain,
		tags:        tags,
		rating:      rating,
		level:       level,
	}, nil

}

// Getters (read-only access for serialization/display)
func (c *Course) ID() string          { return c.id }
func (c *Course) TeacherID() string   { return c.teacherID }
func (c *Course) Title() string       { return c.title }
func (c *Course) Description() string { return c.description }
func (c *Course) Thumbnail() string   { return c.thumbnail }
func (c *Course) Duration() int       { return c.duration }
func (c *Course) Domain() Domain      { return c.domain }
func (c *Course) Tags() []Tag         { return c.tags }
func (c *Course) Rating() float64     { return c.rating }
func (c *Course) Level() CourseLevel  { return c.level }

// Behavior methods
func (c *Course) IsOwnedBy(teacherID string) bool {
	return c.teacherID == teacherID
}

func (c *Course) HasTag(tag Tag) bool {
	for _, t := range c.tags {
		if t.String() == tag.String() {
			return true
		}
	}
	return false
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

func (c *Course) UpdateDuration(duration int) error {
	if duration < 0 {
		return errors.New("duration cannot be negative")
	}
	c.duration = duration
	return nil
}

func (c *Course) UpdateRating(rating float64) error {
	if rating < 0 || rating > 5 {
		return errors.New("rating must be between 0 and 5")
	}
	c.rating = rating
	return nil
}

func (c *Course) AddTag(tag Tag) error {
	if c.HasTag(tag) {
		return errors.New("tag already exists")
	}
	c.tags = append(c.tags, tag)
	return nil
}

func (c *Course) RemoveTag(tag Tag) error {
	for i, t := range c.tags {
		if t.String() == tag.String() {
			c.tags = append(c.tags[:i], c.tags[i+1:]...)
			return nil
		}
	}
	return errors.New("tag not found")
}
