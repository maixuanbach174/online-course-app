package course

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

func (l *Lesson) AddExercise(exercise Exercise) error {
	panic("not implemented")
}

func (l *Lesson) RemoveExercise(exercise Exercise) error {
	panic("not implemented")
}

func (l *Lesson) HasVideo() bool {
	return l.videoID != ""
}

func (l *Lesson) UpdateContent(content string) error {
	panic("not implemented")
}

func (l *Lesson) UpdateOverview(videoID string) error {
	panic("not implemented")
}
