package enrollment

import "time"

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
	return &Enrollment{
		id:             id,
		userID:         userID,
		courseID:       courseID,
		enrolledAt:     time.Now(),
		startedAt:      time.Now(),
		completedAt:    time.Now(),
		courseProgress: CourseProgress{},
		moduleProgress: []ModuleProgress{},
		lessonProgress: []LessonProgress{},
	}, nil
}
