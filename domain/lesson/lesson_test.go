package lesson

import (
	"testing"
)

func TestNewLesson(t *testing.T) {
	t.Parallel()
	t.Run("successfully creates lesson with all fields", func(t *testing.T) {
		lesson, err := NewLesson(
			"lesson-123",
			"module-456",
			"Variables and Data Types",
			"Learn about Go variables",
			"Detailed content about variables...",
			"video-789",
			900,
			1,
		)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if lesson == nil {
			t.Fatal("expected lesson to be created, got nil")
		}
		if lesson.ID() != "lesson-123" {
			t.Errorf("expected ID 'lesson-123', got '%s'", lesson.ID())
		}
		if lesson.ModuleID() != "module-456" {
			t.Errorf("expected ModuleID 'module-456', got '%s'", lesson.ModuleID())
		}
		if lesson.Title() != "Variables and Data Types" {
			t.Errorf("expected Title 'Variables and Data Types', got '%s'", lesson.Title())
		}
		if lesson.Overview() != "Learn about Go variables" {
			t.Errorf("expected Overview 'Learn about Go variables', got '%s'", lesson.Overview())
		}
		if lesson.Content() != "Detailed content about variables..." {
			t.Errorf("expected specific content, got '%s'", lesson.Content())
		}
		if lesson.VideoID() != "video-789" {
			t.Errorf("expected VideoID 'video-789', got '%s'", lesson.VideoID())
		}
		if lesson.Duration() != 900 {
			t.Errorf("expected Duration 900, got %d", lesson.Duration())
		}
		if lesson.Order() != 1 {
			t.Errorf("expected Order 1, got %d", lesson.Order())
		}
	})

	t.Run("fails when id is empty", func(t *testing.T) {
		lesson, err := NewLesson(
			"",
			"module-456",
			"Variables and Data Types",
			"Overview",
			"Content",
			"video-789",
			900,
			1,
		)

		if err == nil {
			t.Fatal("expected error for empty id, got nil")
		}
		if lesson != nil {
			t.Error("expected nil lesson, got lesson instance")
		}
		if err.Error() != "lesson id is required" {
			t.Errorf("expected error 'lesson id is required', got '%s'", err.Error())
		}
	})

	t.Run("fails when moduleID is empty", func(t *testing.T) {
		lesson, err := NewLesson(
			"lesson-123",
			"",
			"Variables and Data Types",
			"Overview",
			"Content",
			"video-789",
			900,
			1,
		)

		if err == nil {
			t.Fatal("expected error for empty moduleID, got nil")
		}
		if lesson != nil {
			t.Error("expected nil lesson, got lesson instance")
		}
		if err.Error() != "module id is required" {
			t.Errorf("expected error 'module id is required', got '%s'", err.Error())
		}
	})

	t.Run("fails when title is empty", func(t *testing.T) {
		lesson, err := NewLesson(
			"lesson-123",
			"module-456",
			"",
			"Overview",
			"Content",
			"video-789",
			900,
			1,
		)

		if err == nil {
			t.Fatal("expected error for empty title, got nil")
		}
		if lesson != nil {
			t.Error("expected nil lesson, got lesson instance")
		}
		if err.Error() != "lesson title is required" {
			t.Errorf("expected error 'lesson title is required', got '%s'", err.Error())
		}
	})

	t.Run("creates lesson with minimal required fields", func(t *testing.T) {
		lesson, err := NewLesson(
			"lesson-123",
			"module-456",
			"Basic Lesson",
			"",
			"",
			"",
			0,
			0,
		)

		if err != nil {
			t.Fatalf("expected no error for minimal valid data, got %v", err)
		}
		if lesson == nil {
			t.Fatal("expected lesson to be created, got nil")
		}
		if lesson.Overview() != "" {
			t.Error("expected empty overview")
		}
		if lesson.Content() != "" {
			t.Error("expected empty content")
		}
		if lesson.VideoID() != "" {
			t.Error("expected empty videoID")
		}
	})
}

func TestLesson_HasVideo(t *testing.T) {
	t.Parallel()
	t.Run("returns true when video exists", func(t *testing.T) {
		lesson, _ := NewLesson(
			"lesson-123",
			"module-456",
			"Video Lesson",
			"",
			"",
			"video-789",
			900,
			1,
		)

		if !lesson.HasVideo() {
			t.Error("expected lesson to have video")
		}
	})

	t.Run("returns false when videoID is empty", func(t *testing.T) {
		lesson, _ := NewLesson(
			"lesson-123",
			"module-456",
			"Text Lesson",
			"",
			"",
			"",
			0,
			1,
		)

		if lesson.HasVideo() {
			t.Error("expected lesson not to have video")
		}
	})
}

func TestLesson_UpdateContent(t *testing.T) {
	t.Parallel()
	lesson, _ := NewLesson(
		"lesson-123",
		"module-456",
		"Lesson Title",
		"",
		"Original content",
		"",
		0,
		1,
	)

	t.Run("successfully updates content", func(t *testing.T) {
		err := lesson.UpdateContent("New detailed content")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if lesson.Content() != "New detailed content" {
			t.Errorf("expected content 'New detailed content', got '%s'", lesson.Content())
		}
	})

	t.Run("allows empty content", func(t *testing.T) {
		err := lesson.UpdateContent("")

		if err != nil {
			t.Fatalf("expected no error for empty content, got %v", err)
		}
		if lesson.Content() != "" {
			t.Error("expected empty content")
		}
	})

	t.Run("allows very long content", func(t *testing.T) {
		longContent := "This is a very long content " + "repeated text " // simulate long content
		err := lesson.UpdateContent(longContent)

		if err != nil {
			t.Fatalf("expected no error for long content, got %v", err)
		}
		if lesson.Content() != longContent {
			t.Error("expected long content to be set correctly")
		}
	})
}

func TestLesson_UpdateOverview(t *testing.T) {
	t.Parallel()
	lesson, _ := NewLesson(
		"lesson-123",
		"module-456",
		"Lesson Title",
		"Original overview",
		"",
		"",
		0,
		1,
	)

	t.Run("successfully updates overview", func(t *testing.T) {
		err := lesson.UpdateOverview("New overview text")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if lesson.Overview() != "New overview text" {
			t.Errorf("expected overview 'New overview text', got '%s'", lesson.Overview())
		}
	})

	t.Run("allows empty overview", func(t *testing.T) {
		err := lesson.UpdateOverview("")

		if err != nil {
			t.Fatalf("expected no error for empty overview, got %v", err)
		}
		if lesson.Overview() != "" {
			t.Error("expected empty overview")
		}
	})
}

func TestLesson_UpdateVideoID(t *testing.T) {
	t.Parallel()
	lesson, _ := NewLesson(
		"lesson-123",
		"module-456",
		"Lesson Title",
		"",
		"",
		"video-123",
		900,
		1,
	)

	t.Run("successfully updates videoID", func(t *testing.T) {
		err := lesson.UpdateVideoID("video-456")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if lesson.VideoID() != "video-456" {
			t.Errorf("expected videoID 'video-456', got '%s'", lesson.VideoID())
		}
		if !lesson.HasVideo() {
			t.Error("expected lesson to have video after update")
		}
	})

	t.Run("allows empty videoID", func(t *testing.T) {
		err := lesson.UpdateVideoID("")

		if err != nil {
			t.Fatalf("expected no error for empty videoID, got %v", err)
		}
		if lesson.VideoID() != "" {
			t.Error("expected empty videoID")
		}
		if lesson.HasVideo() {
			t.Error("expected lesson not to have video after clearing videoID")
		}
	})
}

func TestLesson_UpdateDuration(t *testing.T) {
	t.Parallel()
	lesson, _ := NewLesson(
		"lesson-123",
		"module-456",
		"Lesson Title",
		"",
		"",
		"",
		900,
		1,
	)

	t.Run("successfully updates duration", func(t *testing.T) {
		err := lesson.UpdateDuration(1800)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if lesson.Duration() != 1800 {
			t.Errorf("expected duration 1800, got %d", lesson.Duration())
		}
	})

	t.Run("allows zero duration", func(t *testing.T) {
		err := lesson.UpdateDuration(0)

		if err != nil {
			t.Fatalf("expected no error for zero duration, got %v", err)
		}
		if lesson.Duration() != 0 {
			t.Errorf("expected duration 0, got %d", lesson.Duration())
		}
	})

	t.Run("fails when duration is negative", func(t *testing.T) {
		originalDuration := lesson.Duration()
		err := lesson.UpdateDuration(-100)

		if err == nil {
			t.Fatal("expected error for negative duration, got nil")
		}
		if err.Error() != "duration cannot be negative" {
			t.Errorf("expected error 'duration cannot be negative', got '%s'", err.Error())
		}
		if lesson.Duration() != originalDuration {
			t.Error("expected duration to remain unchanged after failed update")
		}
	})
}

func TestLesson_UpdateOrder(t *testing.T) {
	t.Parallel()
	lesson, _ := NewLesson(
		"lesson-123",
		"module-456",
		"Lesson Title",
		"",
		"",
		"",
		0,
		1,
	)

	t.Run("successfully updates order", func(t *testing.T) {
		err := lesson.UpdateOrder(5)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if lesson.Order() != 5 {
			t.Errorf("expected order 5, got %d", lesson.Order())
		}
	})

	t.Run("allows zero order", func(t *testing.T) {
		err := lesson.UpdateOrder(0)

		if err != nil {
			t.Fatalf("expected no error for zero order, got %v", err)
		}
		if lesson.Order() != 0 {
			t.Errorf("expected order 0, got %d", lesson.Order())
		}
	})

	t.Run("fails when order is negative", func(t *testing.T) {
		originalOrder := lesson.Order()
		err := lesson.UpdateOrder(-1)

		if err == nil {
			t.Fatal("expected error for negative order, got nil")
		}
		if err.Error() != "order cannot be negative" {
			t.Errorf("expected error 'order cannot be negative', got '%s'", err.Error())
		}
		if lesson.Order() != originalOrder {
			t.Error("expected order to remain unchanged after failed update")
		}
	})
}

func TestLesson_UpdateTitle(t *testing.T) {
	t.Parallel()
	lesson, _ := NewLesson(
		"lesson-123",
		"module-456",
		"Original Title",
		"",
		"",
		"",
		0,
		1,
	)

	t.Run("successfully updates title", func(t *testing.T) {
		err := lesson.UpdateTitle("New Lesson Title")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if lesson.Title() != "New Lesson Title" {
			t.Errorf("expected title 'New Lesson Title', got '%s'", lesson.Title())
		}
	})

	t.Run("fails when title is empty", func(t *testing.T) {
		originalTitle := lesson.Title()
		err := lesson.UpdateTitle("")

		if err == nil {
			t.Fatal("expected error for empty title, got nil")
		}
		if err.Error() != "title is required" {
			t.Errorf("expected error 'title is required', got '%s'", err.Error())
		}
		if lesson.Title() != originalTitle {
			t.Error("expected title to remain unchanged after failed update")
		}
	})
}

func TestLesson_Getters(t *testing.T) {
	t.Parallel()
	lesson, _ := NewLesson(
		"lesson-123",
		"module-456",
		"Go Fundamentals",
		"Quick overview",
		"Detailed lesson content",
		"video-789",
		1200,
		3,
	)

	t.Run("ID getter returns correct value", func(t *testing.T) {
		if lesson.ID() != "lesson-123" {
			t.Errorf("expected ID 'lesson-123', got '%s'", lesson.ID())
		}
	})

	t.Run("ModuleID getter returns correct value", func(t *testing.T) {
		if lesson.ModuleID() != "module-456" {
			t.Errorf("expected ModuleID 'module-456', got '%s'", lesson.ModuleID())
		}
	})

	t.Run("Title getter returns correct value", func(t *testing.T) {
		if lesson.Title() != "Go Fundamentals" {
			t.Errorf("expected Title 'Go Fundamentals', got '%s'", lesson.Title())
		}
	})

	t.Run("Overview getter returns correct value", func(t *testing.T) {
		if lesson.Overview() != "Quick overview" {
			t.Errorf("expected Overview 'Quick overview', got '%s'", lesson.Overview())
		}
	})

	t.Run("Content getter returns correct value", func(t *testing.T) {
		if lesson.Content() != "Detailed lesson content" {
			t.Errorf("expected Content 'Detailed lesson content', got '%s'", lesson.Content())
		}
	})

	t.Run("VideoID getter returns correct value", func(t *testing.T) {
		if lesson.VideoID() != "video-789" {
			t.Errorf("expected VideoID 'video-789', got '%s'", lesson.VideoID())
		}
	})

	t.Run("Duration getter returns correct value", func(t *testing.T) {
		if lesson.Duration() != 1200 {
			t.Errorf("expected Duration 1200, got %d", lesson.Duration())
		}
	})

	t.Run("Order getter returns correct value", func(t *testing.T) {
		if lesson.Order() != 3 {
			t.Errorf("expected Order 3, got %d", lesson.Order())
		}
	})
}
