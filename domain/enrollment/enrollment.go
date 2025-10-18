package enrollment

import (
	"time"

	"github.com/pkg/errors"
)

type Enrollment struct {
	id             string
	userID         string
	courseID       string
	enrolledAt     time.Time
	startedAt      time.Time
	completedAt    time.Time
	courseProgress CourseProgress
	moduleProgress []ModuleProgress
	lessonProgress []LessonProgress
}

func NewEnrollment(id string, userID string, courseID string) (*Enrollment, error) {
	if id == "" {
		return nil, errors.New("enrollment id is required")
	}
	if userID == "" {
		return nil, errors.New("user id is required")
	}
	if courseID == "" {
		return nil, errors.New("course id is required")
	}

	return &Enrollment{
		id:             id,
		userID:         userID,
		courseID:       courseID,
		enrolledAt:     time.Now(),
		startedAt:      time.Time{}, // Not started yet
		completedAt:    time.Time{}, // Not completed yet
		courseProgress: NewCourseProgress(),
		moduleProgress: []ModuleProgress{},
		lessonProgress: []LessonProgress{},
	}, nil
}

// Getters (read-only access for serialization/display)
func (e *Enrollment) ID() string                       { return e.id }
func (e *Enrollment) UserID() string                   { return e.userID }
func (e *Enrollment) CourseID() string                 { return e.courseID }
func (e *Enrollment) EnrolledAt() time.Time            { return e.enrolledAt }
func (e *Enrollment) StartedAt() time.Time             { return e.startedAt }
func (e *Enrollment) CompletedAt() time.Time           { return e.completedAt }
func (e *Enrollment) CourseProgress() CourseProgress   { return e.courseProgress }
func (e *Enrollment) ModuleProgress() []ModuleProgress { return e.moduleProgress }
func (e *Enrollment) LessonProgress() []LessonProgress { return e.lessonProgress }

// !!! Change in the future
// Behavior methods
func (e *Enrollment) CompleteLesson(lessonID string) error {
	if lessonID == "" {
		return errors.New("lesson id is required")
	}

	// Find or create lesson progress
	found := false
	for i, lp := range e.lessonProgress {
		if lp.LessonID() == lessonID {
			e.lessonProgress[i].MarkCompleted()
			found = true
			break
		}
	}

	if !found {
		// Create new lesson progress
		newLessonProgress := NewLessonProgress(lessonID)
		newLessonProgress.MarkCompleted()
		e.lessonProgress = append(e.lessonProgress, newLessonProgress)
	}

	// Start the course if not started
	if e.startedAt.IsZero() {
		e.startedAt = time.Now()
	}

	return nil
}

func (e *Enrollment) GetLessonProgress(lessonID string) (*LessonProgress, error) {
	for _, lp := range e.lessonProgress {
		if lp.LessonID() == lessonID {
			return &lp, nil
		}
	}
	return nil, errors.Errorf("lesson progress for lesson '%s' not found", lessonID)
}

func (e *Enrollment) IsCompleted() bool {
	return !e.completedAt.IsZero()
}
