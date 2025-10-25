package app

import (
	"github.com/maixuanbach174/online-course-app/internal/education/app/command"
	"github.com/maixuanbach174/online-course-app/internal/education/app/command/course_command"
	"github.com/maixuanbach174/online-course-app/internal/education/app/query"
	"github.com/maixuanbach174/online-course-app/internal/education/app/query/course_query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	RegisterUser   command.RegisterUserHandler
	EnrollInCourse command.EnrollInCourseHandler
	CompleteLesson command.CompleteLessonHandler
	CreateCourse   course_command.CreateCourseHandler
	DeleteCourse   course_command.DeleteCourseHandler
	UpdateCourse   course_command.UpdateCourseHandler
}

type Queries struct {
	GetAllCourses    course_query.GetAllCoursesHandler
	GetCourseDetails course_query.GetCourseDetailsHandler
	CoursesByTeacher course_query.CourseByTeacherHandler
	GetMyEnrollments query.GetMyEnrollmentsHandler
}
