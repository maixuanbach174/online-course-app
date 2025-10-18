package command

import (
	"context"

	"github.com/maixuanbach174/online-course-app/internal/common/decorator"
	"github.com/maixuanbach174/online-course-app/internal/core/domain/course"
	"github.com/maixuanbach174/online-course-app/internal/core/domain/user"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type CreateCourse struct {
	CourseID    string
	TeacherID   string
	Title       string
	Description string
	Thumbnail   string
	Duration    int
	Domain      string
	Tags        []string
	Level       string
	Modules     []ModuleInput
}

type ModuleInput struct {
	ModuleID string
	Title    string
	Lessons  []LessonInput
	Order    int
}

type LessonInput struct {
	LessonID  string
	Title     string
	Overview  string
	Content   string
	VideoID   string
	Exercises []ExerciseInput
	Order     int
}

type ExerciseInput struct {
	ExerciseID    string
	Question      string
	Answers       []string
	CorrectAnswer string
	Order         int
}

type CreateCourseHandler decorator.CommandHandler[CreateCourse]

type createCourseHandler struct {
	courseRepository course.CourseRepository
	userRepository   user.UserRepository
}

func NewCreateCourseHandler(
	courseRepository course.CourseRepository,
	userRepository user.UserRepository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) CreateCourseHandler {
	if courseRepository == nil {
		panic("course repository is required")
	}
	if userRepository == nil {
		panic("user repository is required")
	}

	return decorator.ApplyCommandDecorators(
		createCourseHandler{
			courseRepository: courseRepository,
			userRepository:   userRepository,
		},
		logger,
		metricsClient,
	)
}

func (h createCourseHandler) Handle(ctx context.Context, cmd CreateCourse) error {
	// Validate input
	if cmd.CourseID == "" {
		return errors.New("course ID is required")
	}
	if cmd.TeacherID == "" {
		return errors.New("teacher ID is required")
	}
	if cmd.Title == "" {
		return errors.New("course title is required")
	}

	// Verify teacher exists and can teach
	teacher, err := h.userRepository.Get(ctx, cmd.TeacherID)
	if err != nil {
		return errors.Wrap(err, "teacher not found")
	}
	if !teacher.CanTeach() {
		return errors.New("user is not authorized to create courses")
	}

	// Parse domain
	domain, err := course.NewDomainFromString(cmd.Domain)
	if err != nil {
		return errors.Wrap(err, "invalid domain")
	}

	// Parse tags
	tags := make([]course.Tag, 0, len(cmd.Tags))
	for _, tagStr := range cmd.Tags {
		tag, err := course.NewTagFromString(tagStr)
		if err != nil {
			return errors.Wrapf(err, "invalid tag: %s", tagStr)
		}
		tags = append(tags, tag)
	}

	// Parse level
	level, err := course.NewCourseLevelFromString(cmd.Level)
	if err != nil {
		return errors.Wrap(err, "invalid course level")
	}

	// Build modules
	modules := make([]course.Module, 0, len(cmd.Modules))
	for _, modInput := range cmd.Modules {
		lessons := make([]course.Lesson, 0, len(modInput.Lessons))
		for _, lesInput := range modInput.Lessons {
			exercises := make([]course.Exercise, 0, len(lesInput.Exercises))
			for _, exInput := range lesInput.Exercises {
				exercise, err := course.NewExercise(
					exInput.ExerciseID,
					exInput.Question,
					exInput.Answers,
					exInput.CorrectAnswer,
					exInput.Order,
				)
				if err != nil {
					return errors.Wrap(err, "failed to create exercise")
				}
				exercises = append(exercises, *exercise)
			}

			lesson, err := course.NewLesson(
				lesInput.LessonID,
				lesInput.Title,
				lesInput.Overview,
				lesInput.Content,
				lesInput.VideoID,
				exercises,
				lesInput.Order,
			)
			if err != nil {
				return errors.Wrap(err, "failed to create lesson")
			}
			lessons = append(lessons, *lesson)
		}

		module, err := course.NewModule(modInput.ModuleID, modInput.Title, lessons, modInput.Order)
		if err != nil {
			return errors.Wrap(err, "failed to create module")
		}
		modules = append(modules, *module)
	}

	// Create course entity
	newCourse, err := course.NewCourse(
		cmd.CourseID,
		cmd.TeacherID,
		cmd.Title,
		cmd.Description,
		cmd.Thumbnail,
		cmd.Duration,
		domain,
		tags,
		0.0, // initial rating
		level,
		modules,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create course")
	}

	// Persist to repository
	if err := h.courseRepository.Create(ctx, newCourse); err != nil {
		return errors.Wrap(err, "failed to save course")
	}

	return nil
}
