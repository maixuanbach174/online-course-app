package course

import (
	"testing"
)

func TestNewCourse(t *testing.T) {
	t.Run("successfully creates course with valid data", func(t *testing.T) {
		course, err := NewCourse(
			"course-123",
			"teacher-456",
			"Learn Go Programming",
			"A comprehensive Go course",
			"thumbnail.jpg",
			3600,
			DomainProgramming,
			[]Tag{TagBackend, TagAPI},
			4.5,
			Beginner,
		)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if course == nil {
			t.Fatal("expected course to be created, got nil")
		}
		if course.ID() != "course-123" {
			t.Errorf("expected ID 'course-123', got '%s'", course.ID())
		}
		if course.TeacherID() != "teacher-456" {
			t.Errorf("expected TeacherID 'teacher-456', got '%s'", course.TeacherID())
		}
		if course.Title() != "Learn Go Programming" {
			t.Errorf("expected Title 'Learn Go Programming', got '%s'", course.Title())
		}
		if course.Description() != "A comprehensive Go course" {
			t.Errorf("expected Description 'A comprehensive Go course', got '%s'", course.Description())
		}
		if course.Thumbnail() != "thumbnail.jpg" {
			t.Errorf("expected Thumbnail 'thumbnail.jpg', got '%s'", course.Thumbnail())
		}
		if course.Duration() != 3600 {
			t.Errorf("expected Duration 3600, got %d", course.Duration())
		}
		if course.Domain().String() != "programming" {
			t.Errorf("expected Domain 'programming', got '%s'", course.Domain().String())
		}
		if len(course.Tags()) != 2 {
			t.Errorf("expected 2 tags, got %d", len(course.Tags()))
		}
		if course.Rating() != 4.5 {
			t.Errorf("expected Rating 4.5, got %f", course.Rating())
		}
		if course.Level().String() != "beginner" {
			t.Errorf("expected Level 'beginner', got '%s'", course.Level().String())
		}
	})

	t.Run("fails when id is empty", func(t *testing.T) {
		course, err := NewCourse(
			"",
			"teacher-456",
			"Learn Go Programming",
			"A comprehensive Go course",
			"thumbnail.jpg",
			3600,
			DomainProgramming,
			[]Tag{TagBackend},
			4.5,
			Beginner,
		)

		if err == nil {
			t.Fatal("expected error for empty id, got nil")
		}
		if course != nil {
			t.Error("expected nil course, got course instance")
		}
		if err.Error() != "course id is required" {
			t.Errorf("expected error 'course id is required', got '%s'", err.Error())
		}
	})

	t.Run("fails when teacherID is empty", func(t *testing.T) {
		course, err := NewCourse(
			"course-123",
			"",
			"Learn Go Programming",
			"A comprehensive Go course",
			"thumbnail.jpg",
			3600,
			DomainProgramming,
			[]Tag{TagBackend},
			4.5,
			Beginner,
		)

		if err == nil {
			t.Fatal("expected error for empty teacherID, got nil")
		}
		if course != nil {
			t.Error("expected nil course, got course instance")
		}
		if err.Error() != "teacher id is required" {
			t.Errorf("expected error 'teacher id is required', got '%s'", err.Error())
		}
	})

	t.Run("fails when title is empty", func(t *testing.T) {
		course, err := NewCourse(
			"course-123",
			"teacher-456",
			"",
			"A comprehensive Go course",
			"thumbnail.jpg",
			3600,
			DomainProgramming,
			[]Tag{TagBackend},
			4.5,
			Beginner,
		)

		if err == nil {
			t.Fatal("expected error for empty title, got nil")
		}
		if course != nil {
			t.Error("expected nil course, got course instance")
		}
		if err.Error() != "course title is required" {
			t.Errorf("expected error 'course title is required', got '%s'", err.Error())
		}
	})

	t.Run("creates course with minimal required fields", func(t *testing.T) {
		course, err := NewCourse(
			"course-123",
			"teacher-456",
			"Minimal Course",
			"",
			"",
			0,
			Domain{},
			nil,
			0,
			CourseLevel{},
		)

		if err != nil {
			t.Fatalf("expected no error for minimal valid data, got %v", err)
		}
		if course == nil {
			t.Fatal("expected course to be created, got nil")
		}
	})
}

func TestCourse_IsOwnedBy(t *testing.T) {
	course, _ := NewCourse(
		"course-123",
		"teacher-456",
		"Learn Go",
		"",
		"",
		0,
		DomainProgramming,
		nil,
		0,
		Beginner,
	)

	t.Run("returns true for correct teacher", func(t *testing.T) {
		if !course.IsOwnedBy("teacher-456") {
			t.Error("expected course to be owned by 'teacher-456'")
		}
	})

	t.Run("returns false for different teacher", func(t *testing.T) {
		if course.IsOwnedBy("teacher-999") {
			t.Error("expected course not to be owned by 'teacher-999'")
		}
	})
}

func TestCourse_HasTag(t *testing.T) {
	course, _ := NewCourse(
		"course-123",
		"teacher-456",
		"Learn Go",
		"",
		"",
		0,
		DomainProgramming,
		[]Tag{TagBackend, TagAPI, TagCloud},
		0,
		Beginner,
	)

	t.Run("returns true when tag exists", func(t *testing.T) {
		if !course.HasTag(TagBackend) {
			t.Error("expected course to have TagBackend")
		}
		if !course.HasTag(TagAPI) {
			t.Error("expected course to have TagAPI")
		}
	})

	t.Run("returns false when tag does not exist", func(t *testing.T) {
		if course.HasTag(TagFrontend) {
			t.Error("expected course not to have TagFrontend")
		}
	})
}

func TestCourse_UpdateBasicInfo(t *testing.T) {
	course, _ := NewCourse(
		"course-123",
		"teacher-456",
		"Original Title",
		"Original Description",
		"original.jpg",
		3600,
		DomainProgramming,
		nil,
		0,
		Beginner,
	)

	t.Run("successfully updates all basic info", func(t *testing.T) {
		err := course.UpdateBasicInfo(
			"New Title",
			"New Description",
			"new.jpg",
		)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if course.Title() != "New Title" {
			t.Errorf("expected title 'New Title', got '%s'", course.Title())
		}
		if course.Description() != "New Description" {
			t.Errorf("expected description 'New Description', got '%s'", course.Description())
		}
		if course.Thumbnail() != "new.jpg" {
			t.Errorf("expected thumbnail 'new.jpg', got '%s'", course.Thumbnail())
		}
	})

	t.Run("fails when title is empty", func(t *testing.T) {
		originalTitle := course.Title()
		err := course.UpdateBasicInfo("", "Some description", "thumb.jpg")

		if err == nil {
			t.Fatal("expected error for empty title, got nil")
		}
		if err.Error() != "title is required" {
			t.Errorf("expected error 'title is required', got '%s'", err.Error())
		}
		if course.Title() != originalTitle {
			t.Error("expected title to remain unchanged after failed update")
		}
	})

	t.Run("allows empty description and thumbnail", func(t *testing.T) {
		err := course.UpdateBasicInfo("Valid Title", "", "")

		if err != nil {
			t.Fatalf("expected no error for empty description/thumbnail, got %v", err)
		}
		if course.Description() != "" {
			t.Error("expected empty description")
		}
		if course.Thumbnail() != "" {
			t.Error("expected empty thumbnail")
		}
	})
}

func TestCourse_UpdateDuration(t *testing.T) {
	course, _ := NewCourse(
		"course-123",
		"teacher-456",
		"Learn Go",
		"",
		"",
		3600,
		DomainProgramming,
		nil,
		0,
		Beginner,
	)

	t.Run("successfully updates duration", func(t *testing.T) {
		err := course.UpdateDuration(7200)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if course.Duration() != 7200 {
			t.Errorf("expected duration 7200, got %d", course.Duration())
		}
	})

	t.Run("allows zero duration", func(t *testing.T) {
		err := course.UpdateDuration(0)

		if err != nil {
			t.Fatalf("expected no error for zero duration, got %v", err)
		}
		if course.Duration() != 0 {
			t.Errorf("expected duration 0, got %d", course.Duration())
		}
	})

	t.Run("fails when duration is negative", func(t *testing.T) {
		originalDuration := course.Duration()
		err := course.UpdateDuration(-100)

		if err == nil {
			t.Fatal("expected error for negative duration, got nil")
		}
		if err.Error() != "duration cannot be negative" {
			t.Errorf("expected error 'duration cannot be negative', got '%s'", err.Error())
		}
		if course.Duration() != originalDuration {
			t.Error("expected duration to remain unchanged after failed update")
		}
	})
}

func TestCourse_UpdateRating(t *testing.T) {
	course, _ := NewCourse(
		"course-123",
		"teacher-456",
		"Learn Go",
		"",
		"",
		0,
		DomainProgramming,
		nil,
		4.0,
		Beginner,
	)

	t.Run("successfully updates rating", func(t *testing.T) {
		err := course.UpdateRating(4.8)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if course.Rating() != 4.8 {
			t.Errorf("expected rating 4.8, got %f", course.Rating())
		}
	})

	t.Run("allows rating of 0", func(t *testing.T) {
		err := course.UpdateRating(0)

		if err != nil {
			t.Fatalf("expected no error for rating 0, got %v", err)
		}
		if course.Rating() != 0 {
			t.Errorf("expected rating 0, got %f", course.Rating())
		}
	})

	t.Run("allows rating of 5", func(t *testing.T) {
		err := course.UpdateRating(5)

		if err != nil {
			t.Fatalf("expected no error for rating 5, got %v", err)
		}
		if course.Rating() != 5 {
			t.Errorf("expected rating 5, got %f", course.Rating())
		}
	})

	t.Run("fails when rating is negative", func(t *testing.T) {
		originalRating := course.Rating()
		err := course.UpdateRating(-0.1)

		if err == nil {
			t.Fatal("expected error for negative rating, got nil")
		}
		if err.Error() != "rating must be between 0 and 5" {
			t.Errorf("expected error 'rating must be between 0 and 5', got '%s'", err.Error())
		}
		if course.Rating() != originalRating {
			t.Error("expected rating to remain unchanged after failed update")
		}
	})

	t.Run("fails when rating exceeds 5", func(t *testing.T) {
		originalRating := course.Rating()
		err := course.UpdateRating(5.1)

		if err == nil {
			t.Fatal("expected error for rating > 5, got nil")
		}
		if err.Error() != "rating must be between 0 and 5" {
			t.Errorf("expected error 'rating must be between 0 and 5', got '%s'", err.Error())
		}
		if course.Rating() != originalRating {
			t.Error("expected rating to remain unchanged after failed update")
		}
	})
}

func TestCourse_AddTag(t *testing.T) {
	course, _ := NewCourse(
		"course-123",
		"teacher-456",
		"Learn Go",
		"",
		"",
		0,
		DomainProgramming,
		[]Tag{TagBackend},
		0,
		Beginner,
	)

	t.Run("successfully adds new tag", func(t *testing.T) {
		err := course.AddTag(TagAPI)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !course.HasTag(TagAPI) {
			t.Error("expected course to have TagAPI after adding")
		}
		if len(course.Tags()) != 2 {
			t.Errorf("expected 2 tags, got %d", len(course.Tags()))
		}
	})

	t.Run("fails when adding duplicate tag", func(t *testing.T) {
		originalTagCount := len(course.Tags())
		err := course.AddTag(TagBackend)

		if err == nil {
			t.Fatal("expected error for duplicate tag, got nil")
		}
		if err.Error() != "tag already exists" {
			t.Errorf("expected error 'tag already exists', got '%s'", err.Error())
		}
		if len(course.Tags()) != originalTagCount {
			t.Error("expected tag count to remain unchanged after failed add")
		}
	})

	t.Run("can add multiple different tags", func(t *testing.T) {
		err1 := course.AddTag(TagCloud)
		err2 := course.AddTag(TagDatabase)

		if err1 != nil || err2 != nil {
			t.Fatal("expected no errors when adding multiple different tags")
		}
		if len(course.Tags()) != 4 {
			t.Errorf("expected 4 tags, got %d", len(course.Tags()))
		}
	})
}

func TestCourse_RemoveTag(t *testing.T) {
	course, _ := NewCourse(
		"course-123",
		"teacher-456",
		"Learn Go",
		"",
		"",
		0,
		DomainProgramming,
		[]Tag{TagBackend, TagAPI, TagCloud},
		0,
		Beginner,
	)

	t.Run("successfully removes existing tag", func(t *testing.T) {
		err := course.RemoveTag(TagAPI)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if course.HasTag(TagAPI) {
			t.Error("expected course not to have TagAPI after removal")
		}
		if len(course.Tags()) != 2 {
			t.Errorf("expected 2 tags, got %d", len(course.Tags()))
		}
	})

	t.Run("fails when removing non-existent tag", func(t *testing.T) {
		originalTagCount := len(course.Tags())
		err := course.RemoveTag(TagFrontend)

		if err == nil {
			t.Fatal("expected error for non-existent tag, got nil")
		}
		if err.Error() != "tag not found" {
			t.Errorf("expected error 'tag not found', got '%s'", err.Error())
		}
		if len(course.Tags()) != originalTagCount {
			t.Error("expected tag count to remain unchanged after failed remove")
		}
	})

	t.Run("can remove all tags", func(t *testing.T) {
		course.RemoveTag(TagBackend)
		course.RemoveTag(TagCloud)

		if len(course.Tags()) != 0 {
			t.Errorf("expected 0 tags, got %d", len(course.Tags()))
		}
	})
}
