package course

import "github.com/pkg/errors"

type Lesson struct {
	id        string
	title     string
	overview  string
	content   string
	videoID   string
	exercises []Exercise
	order     int
}

func NewLesson(id string, title string, overview string, content string, videoID string, exercises []Exercise, order int) (*Lesson, error) {
	return &Lesson{
		id:        id,
		title:     title,
		overview:  overview,
		content:   content,
		videoID:   videoID,
		exercises: exercises,
		order:     order,
	}, nil
}

// Getters (read-only access for serialization/display)
func (l *Lesson) ID() string            { return l.id }
func (l *Lesson) Title() string         { return l.title }
func (l *Lesson) Overview() string      { return l.overview }
func (l *Lesson) Content() string       { return l.content }
func (l *Lesson) VideoID() string       { return l.videoID }
func (l *Lesson) Exercises() []Exercise { return l.exercises }
func (l *Lesson) Order() int            { return l.order }

// Behavior methods
func (l *Lesson) AddExercise(exercise Exercise) error {
	l.exercises = append(l.exercises, exercise)
	return nil
}

func (l *Lesson) RemoveExercise(exerciseID string) error {
	for i, e := range l.exercises {
		if e.ID() == exerciseID {
			l.exercises = append(l.exercises[:i], l.exercises[i+1:]...)
			return nil
		}
	}
	return errors.New("exercise not found")
}

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

func (l *Lesson) UpdateOrder(order int) {
	l.order = order
}

func (l *Lesson) UpdateTitle(title string) error {
	if title == "" {
		return errors.New("title is required")
	}
	l.title = title
	return nil
}

func (l *Lesson) CalculateDuration() int {
	// For simplicity, assume each lesson is 10 minutes + 2 minutes per exercise
	// You can adjust this logic based on your needs
	baseDuration := 10
	exerciseDuration := len(l.exercises) * 2
	return baseDuration + exerciseDuration
}

func (l *Lesson) GetExerciseByID(exerciseID string) (*Exercise, error) {
	for _, e := range l.exercises {
		if e.ID() == exerciseID {
			return &e, nil
		}
	}
	return nil, errors.New("exercise not found")
}
