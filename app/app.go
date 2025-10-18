package app

import (
	"github.com/maixuanbach174/online-course-app/internal/core/app/command"
	"github.com/maixuanbach174/online-course-app/internal/core/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	RegisterUser   command.RegisterUserHandler
	CreateCourse   command.CreateCourseHandler
	EnrollInCourse command.EnrollInCourseHandler
	CompleteLesson command.CompleteLessonHandler
}

type Queries struct {
	GetAllCourses    query.GetAllCoursesHandler
	GetCourseDetails query.GetCourseDetailsHandler
	CoursesByTeacher query.CourseByTeacherHandler
	GetMyEnrollments query.GetMyEnrollmentsHandler
}
