package ports

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/maixuanbach174/online-course-app/internal/common/server/httperr"
	"github.com/maixuanbach174/online-course-app/internal/education/app"
	"github.com/maixuanbach174/online-course-app/internal/education/app/command/course_command"
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

	// Handle optional parameters safely
	var domain, level, tag string
	if params.Domain != nil {
		domain = string(*params.Domain)
	}
	if params.Level != nil {
		level = string(*params.Level)
	}
	if params.Tag != nil {
		tag = string(*params.Tag)
	}

	courses, err := h.app.Queries.GetAllCourses.Handle(context, course_query.GetAllCourses{
		Domain: domain,
		Level:  level,
		Tag:    tag,
	})

	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
	render.Respond(w, r, courses)
}

func (h HttpServer) CreateCourse(w http.ResponseWriter, r *http.Request) {
	var req CreateCourseRequest
	if err := render.Decode(r, &req); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}

	var tags []string
	if req.Tags != nil {
		for _, tag := range *req.Tags {
			tags = append(tags, string(tag))
		}
	}

	var rating float64
	if req.Rating != nil {
		rating = float64(*req.Rating)
	}

	err := h.app.Commands.CreateCourse.Handle(r.Context(), course_command.CreateCourse{
		TeacherID:   req.TeacherId,
		Title:       req.Title,
		Description: getStringValue(req.Description),
		Thumbnail:   getStringValue(req.Thumbnail),
		Duration:    req.Duration,
		Domain:      string(req.Domain),
		Tags:        tags,
		Rating:      rating,
		Level:       string(req.Level),
	})

	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h HttpServer) DeleteCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	err := h.app.Commands.DeleteCourse.Handle(r.Context(), course_command.DeleteCourse{
		CourseID: courseId,
	})

	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) GetCourseById(w http.ResponseWriter, r *http.Request, courseId string) {
	course, err := h.app.Queries.GetCourseDetails.Handle(r.Context(), course_query.GetCourseDetails{
		CourseID: courseId,
	})

	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	render.Respond(w, r, course)
}

func (h HttpServer) UpdateCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	var req UpdateCourseRequest
	if err := render.Decode(r, &req); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}

	// Get the existing course first to populate fields that aren't being updated
	existingCourse, err := h.app.Queries.GetCourseDetails.Handle(r.Context(), course_query.GetCourseDetails{
		CourseID: courseId,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	// Build the update command, using existing values for fields not provided in the request
	updateCmd := course_command.UpdateCourse{
		CourseID:    courseId,
		TeacherID:   existingCourse.TeacherID(),
		Title:       getStringValueWithDefault(req.Title, existingCourse.Title()),
		Description: getStringValueWithDefault(req.Description, existingCourse.Description()),
		Thumbnail:   getStringValueWithDefault(req.Thumbnail, existingCourse.Thumbnail()),
		Duration:    getIntValueWithDefault(req.Duration, existingCourse.Duration()),
		Domain:      getStringValueWithDefault((*string)(req.Domain), existingCourse.Domain().String()),
		Level:       getStringValueWithDefault((*string)(req.Level), existingCourse.Level().String()),
		Rating:      getFloat64ValueWithDefault(req.Rating, existingCourse.Rating()),
	}

	// Handle tags
	if req.Tags != nil {
		var tags []string
		for _, tag := range *req.Tags {
			tags = append(tags, string(tag))
		}
		updateCmd.Tags = tags
	} else {
		// Keep existing tags
		var tags []string
		for _, tag := range existingCourse.Tags() {
			tags = append(tags, tag.String())
		}
		updateCmd.Tags = tags
	}

	err = h.app.Commands.UpdateCourse.Handle(r.Context(), updateCmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) GetCoursesByTeacher(w http.ResponseWriter, r *http.Request, teacherId string) {
	courses, err := h.app.Queries.CoursesByTeacher.Handle(r.Context(), course_query.CourseByTeacherQuery{
		TeacherID: teacherId,
	})

	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	render.Respond(w, r, courses)
}

// Helper function to safely get string value from pointer
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Helper function to get string value from pointer with default
func getStringValueWithDefault(s *string, defaultValue string) string {
	if s == nil {
		return defaultValue
	}
	return *s
}

// Helper function to get int value from pointer with default
func getIntValueWithDefault(i *int, defaultValue int) int {
	if i == nil {
		return defaultValue
	}
	return *i
}

// Helper function to get float64 value from pointer with default
func getFloat64ValueWithDefault(f *float32, defaultValue float64) float64 {
	if f == nil {
		return defaultValue
	}
	return float64(*f)
}
