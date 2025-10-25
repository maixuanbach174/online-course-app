package lesson

import "github.com/pkg/errors"

type Lesson struct {
	id       string
	moduleID string
	title    string
	overview string
	content  string
	videoID  string
	duration int
	order    int
}

func NewLesson(id string, moduleID string, title string, overview string, content string, videoID string, duration int, order int) (*Lesson, error) {
	if id == "" {
		return nil, errors.New("lesson id is required")
	}
	if moduleID == "" {
		return nil, errors.New("module id is required")
	}
	if title == "" {
		return nil, errors.New("lesson title is required")
	}

	return &Lesson{
		id:       id,
		moduleID: moduleID,
		title:    title,
		overview: overview,
		content:  content,
		videoID:  videoID,
		duration: duration,
		order:    order,
	}, nil
}

// Getters (read-only access for serialization/display)
func (l *Lesson) ID() string       { return l.id }
func (l *Lesson) ModuleID() string { return l.moduleID }
func (l *Lesson) Title() string    { return l.title }
func (l *Lesson) Overview() string { return l.overview }
func (l *Lesson) Content() string  { return l.content }
func (l *Lesson) VideoID() string  { return l.videoID }
func (l *Lesson) Duration() int    { return l.duration }
func (l *Lesson) Order() int       { return l.order }

// Behavior methods
func (l *Lesson) HasVideo() bool {
	return l.videoID != ""
}

func (l *Lesson) UpdateContent(content string) error {
	l.content = content
	return nil
}

func (l *Lesson) UpdateOverview(overview string) error {
	l.overview = overview
	return nil
}

func (l *Lesson) UpdateVideoID(videoID string) error {
	l.videoID = videoID
	return nil
}

func (l *Lesson) UpdateDuration(duration int) error {
	if duration < 0 {
		return errors.New("duration cannot be negative")
	}
	l.duration = duration
	return nil
}

func (l *Lesson) UpdateOrder(order int) error {
	if order < 0 {
		return errors.New("order cannot be negative")
	}
	l.order = order
	return nil
}

func (l *Lesson) UpdateTitle(title string) error {
	if title == "" {
		return errors.New("title is required")
	}
	l.title = title
	return nil
}
