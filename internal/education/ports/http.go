package ports

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/maixuanbach174/online-course-app/internal/common/server/httperr"
	"github.com/maixuanbach174/online-course-app/internal/education/app"
	"github.com/maixuanbach174/online-course-app/internal/education/app/query/course_query"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(application app.Application) HttpServer {
	return HttpServer{
		app: application,
	}
}

func (h HttpServer) GetCourses(w http.ResponseWriter, r *http.Request, params GetCoursesParams) {
	context := r.Context()

	courses, err := h.app.Queries.GetAllCourses.Handle(context, course_query.GetAllCourses{
		Domain: string(*params.Domain),
		Level:  string(*params.Level),
		Tag:    string(*params.Tag),
	})

	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
	render.Respond(w, r, courses)
}

func (h HttpServer) CreateCourse(w http.ResponseWriter, r *http.Request) {

}

func (h HttpServer) DeleteCourse(w http.ResponseWriter, r *http.Request, courseId string) {

}

func (h HttpServer) GetCourseById(w http.ResponseWriter, r *http.Request, courseId string) {

}

func (h HttpServer) UpdateCourse(w http.ResponseWriter, r *http.Request, courseId string) {

}

func (h HttpServer) GetCoursesByTeacher(w http.ResponseWriter, r *http.Request, teacherId string) {

}
