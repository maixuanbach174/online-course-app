package module

import (
	"testing"
)

func TestNewModule(t *testing.T) {
	t.Parallel()
	t.Run("successfully creates module with valid data", func(t *testing.T) {
		module, err := NewModule(
			"module-123",
			"course-456",
			"Introduction to Go",
			1,
		)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if module == nil {
			t.Fatal("expected module to be created, got nil")
		}
		if module.ID() != "module-123" {
			t.Errorf("expected ID 'module-123', got '%s'", module.ID())
		}
		if module.CourseID() != "course-456" {
			t.Errorf("expected CourseID 'course-456', got '%s'", module.CourseID())
		}
		if module.Title() != "Introduction to Go" {
			t.Errorf("expected Title 'Introduction to Go', got '%s'", module.Title())
		}
		if module.Order() != 1 {
			t.Errorf("expected Order 1, got %d", module.Order())
		}
	})

	t.Run("fails when id is empty", func(t *testing.T) {
		module, err := NewModule(
			"",
			"course-456",
			"Introduction to Go",
			1,
		)

		if err == nil {
			t.Fatal("expected error for empty id, got nil")
		}
		if module != nil {
			t.Error("expected nil module, got module instance")
		}
		if err.Error() != "module id is required" {
			t.Errorf("expected error 'module id is required', got '%s'", err.Error())
		}
	})

	t.Run("fails when courseID is empty", func(t *testing.T) {
		module, err := NewModule(
			"module-123",
			"",
			"Introduction to Go",
			1,
		)

		if err == nil {
			t.Fatal("expected error for empty courseID, got nil")
		}
		if module != nil {
			t.Error("expected nil module, got module instance")
		}
		if err.Error() != "course id is required" {
			t.Errorf("expected error 'course id is required', got '%s'", err.Error())
		}
	})

	t.Run("fails when title is empty", func(t *testing.T) {
		module, err := NewModule(
			"module-123",
			"course-456",
			"",
			1,
		)

		if err == nil {
			t.Fatal("expected error for empty title, got nil")
		}
		if module != nil {
			t.Error("expected nil module, got module instance")
		}
		if err.Error() != "module title is required" {
			t.Errorf("expected error 'module title is required', got '%s'", err.Error())
		}
	})

	t.Run("creates module with zero order", func(t *testing.T) {
		module, err := NewModule(
			"module-123",
			"course-456",
			"Introduction",
			0,
		)

		if err != nil {
			t.Fatalf("expected no error for zero order, got %v", err)
		}
		if module.Order() != 0 {
			t.Errorf("expected Order 0, got %d", module.Order())
		}
	})

	t.Run("creates module with negative order", func(t *testing.T) {
		module, err := NewModule(
			"module-123",
			"course-456",
			"Introduction",
			-1,
		)

		if err != nil {
			t.Fatalf("expected no error (validation only in UpdateOrder), got %v", err)
		}
		if module.Order() != -1 {
			t.Errorf("expected Order -1, got %d", module.Order())
		}
	})
}

func TestModule_UpdateOrder(t *testing.T) {
	t.Parallel()
	module, _ := NewModule(
		"module-123",
		"course-456",
		"Introduction to Go",
		1,
	)

	t.Run("successfully updates order", func(t *testing.T) {
		err := module.UpdateOrder(3)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if module.Order() != 3 {
			t.Errorf("expected order 3, got %d", module.Order())
		}
	})

	t.Run("allows zero order", func(t *testing.T) {
		err := module.UpdateOrder(0)

		if err != nil {
			t.Fatalf("expected no error for zero order, got %v", err)
		}
		if module.Order() != 0 {
			t.Errorf("expected order 0, got %d", module.Order())
		}
	})

	t.Run("fails when order is negative", func(t *testing.T) {
		originalOrder := module.Order()
		err := module.UpdateOrder(-1)

		if err == nil {
			t.Fatal("expected error for negative order, got nil")
		}
		if err.Error() != "order cannot be negative" {
			t.Errorf("expected error 'order cannot be negative', got '%s'", err.Error())
		}
		if module.Order() != originalOrder {
			t.Error("expected order to remain unchanged after failed update")
		}
	})

	t.Run("can update order multiple times", func(t *testing.T) {
		module.UpdateOrder(5)
		module.UpdateOrder(2)
		module.UpdateOrder(10)

		if module.Order() != 10 {
			t.Errorf("expected final order 10, got %d", module.Order())
		}
	})
}

func TestModule_UpdateTitle(t *testing.T) {
	t.Parallel()
	module, _ := NewModule(
		"module-123",
		"course-456",
		"Original Title",
		1,
	)

	t.Run("successfully updates title", func(t *testing.T) {
		err := module.UpdateTitle("Advanced Go Concepts")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if module.Title() != "Advanced Go Concepts" {
			t.Errorf("expected title 'Advanced Go Concepts', got '%s'", module.Title())
		}
	})

	t.Run("fails when title is empty", func(t *testing.T) {
		originalTitle := module.Title()
		err := module.UpdateTitle("")

		if err == nil {
			t.Fatal("expected error for empty title, got nil")
		}
		if err.Error() != "title is required" {
			t.Errorf("expected error 'title is required', got '%s'", err.Error())
		}
		if module.Title() != originalTitle {
			t.Error("expected title to remain unchanged after failed update")
		}
	})

	t.Run("can update title multiple times", func(t *testing.T) {
		module.UpdateTitle("First Update")
		module.UpdateTitle("Second Update")
		module.UpdateTitle("Final Title")

		if module.Title() != "Final Title" {
			t.Errorf("expected final title 'Final Title', got '%s'", module.Title())
		}
	})

	t.Run("allows very long title", func(t *testing.T) {
		longTitle := "This is a very long module title that contains many words and characters to test if the system can handle long titles properly without any issues"
		err := module.UpdateTitle(longTitle)

		if err != nil {
			t.Fatalf("expected no error for long title, got %v", err)
		}
		if module.Title() != longTitle {
			t.Error("expected long title to be set correctly")
		}
	})
}

func TestModule_Getters(t *testing.T) {
	t.Parallel()
	module, _ := NewModule(
		"module-123",
		"course-456",
		"Introduction to Go",
		5,
	)

	t.Run("ID getter returns correct value", func(t *testing.T) {
		if module.ID() != "module-123" {
			t.Errorf("expected ID 'module-123', got '%s'", module.ID())
		}
	})

	t.Run("CourseID getter returns correct value", func(t *testing.T) {
		if module.CourseID() != "course-456" {
			t.Errorf("expected CourseID 'course-456', got '%s'", module.CourseID())
		}
	})

	t.Run("Title getter returns correct value", func(t *testing.T) {
		if module.Title() != "Introduction to Go" {
			t.Errorf("expected Title 'Introduction to Go', got '%s'", module.Title())
		}
	})

	t.Run("Order getter returns correct value", func(t *testing.T) {
		if module.Order() != 5 {
			t.Errorf("expected Order 5, got %d", module.Order())
		}
	})
}
