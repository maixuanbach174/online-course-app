package enrollment

type Progress struct {
	progress float64
	status   Status
}

type CourseProgress struct {
	progress Progress
}

type ModuleProgress struct {
	moduleID string
	progress Progress
}

type LessonProgress struct {
	lessonID      string
	progress      Progress
	exerciseScore float64
}

func NewProgress(progress float64, status Status) Progress {
	return Progress{
		progress: progress,
		status:   status,
	}
}

func NewCourseProgress() CourseProgress {
	return CourseProgress{
		progress: NewProgress(0, Started),
	}
}

func NewModuleProgress(moduleID string) ModuleProgress {
	return ModuleProgress{
		moduleID: moduleID,
		progress: NewProgress(0, Started),
	}
}

func NewLessonProgress(lessonID string) LessonProgress {
	return LessonProgress{
		lessonID:      lessonID,
		progress:      NewProgress(0, Started),
		exerciseScore: 0,
	}
}
