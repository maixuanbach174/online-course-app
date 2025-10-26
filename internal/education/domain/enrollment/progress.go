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
		progress: NewProgress(0, Enrolled),
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

// Progress getters
func (p Progress) ProgressPercentage() float64 { return p.progress }
func (p Progress) Status() Status              { return p.status }

// LessonProgress methods
func (lp LessonProgress) LessonID() string         { return lp.lessonID }
func (lp LessonProgress) Progress() Progress       { return lp.progress }
func (lp LessonProgress) ExerciseScore() float64   { return lp.exerciseScore }

func (lp *LessonProgress) MarkCompleted() {
	lp.progress.progress = 100.0
	lp.progress.status = Completed
}

func (lp *LessonProgress) UpdateExerciseScore(score float64) {
	lp.exerciseScore = score
}

// ModuleProgress methods
func (mp ModuleProgress) ModuleID() string   { return mp.moduleID }
func (mp ModuleProgress) Progress() Progress { return mp.progress }

// CourseProgress methods
func (cp CourseProgress) Progress() Progress { return cp.progress }
