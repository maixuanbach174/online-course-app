package course

import (
	"testing"
)

func TestCourseLevel_String(t *testing.T) {
	t.Run("Beginner level returns correct string", func(t *testing.T) {
		if Beginner.String() != "beginner" {
			t.Errorf("expected 'beginner', got '%s'", Beginner.String())
		}
	})

	t.Run("Intermediate level returns correct string", func(t *testing.T) {
		if Intermediate.String() != "intermediate" {
			t.Errorf("expected 'intermediate', got '%s'", Intermediate.String())
		}
	})

	t.Run("Advanced level returns correct string", func(t *testing.T) {
		if Advanced.String() != "advanced" {
			t.Errorf("expected 'advanced', got '%s'", Advanced.String())
		}
	})
}

func TestNewCourseLevelFromString(t *testing.T) {
	t.Run("successfully creates Beginner from string", func(t *testing.T) {
		level, err := NewCourseLevelFromString("beginner")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if level.String() != "beginner" {
			t.Errorf("expected 'beginner', got '%s'", level.String())
		}
	})

	t.Run("successfully creates Intermediate from string", func(t *testing.T) {
		level, err := NewCourseLevelFromString("intermediate")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if level.String() != "intermediate" {
			t.Errorf("expected 'intermediate', got '%s'", level.String())
		}
	})

	t.Run("successfully creates Advanced from string", func(t *testing.T) {
		level, err := NewCourseLevelFromString("advanced")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if level.String() != "advanced" {
			t.Errorf("expected 'advanced', got '%s'", level.String())
		}
	})

	t.Run("fails for invalid level string", func(t *testing.T) {
		level, err := NewCourseLevelFromString("expert")

		if err == nil {
			t.Fatal("expected error for invalid level, got nil")
		}
		if level.String() != "" {
			t.Error("expected empty level for failed creation")
		}
		expectedError := "unknown 'expert' course level"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("fails for empty string", func(t *testing.T) {
		_, err := NewCourseLevelFromString("")

		if err == nil {
			t.Fatal("expected error for empty string, got nil")
		}
	})

	t.Run("is case sensitive", func(t *testing.T) {
		_, err := NewCourseLevelFromString("Beginner")

		if err == nil {
			t.Fatal("expected error for uppercase 'Beginner', got nil")
		}
	})
}

func TestDomain_String(t *testing.T) {
	domains := []struct {
		domain   Domain
		expected string
	}{
		{DomainProgramming, "programming"},
		{DomainDesign, "design"},
		{DomainBusiness, "business"},
		{DomainMarketing, "marketing"},
		{DomainDataScience, "data_science"},
		{DomainPersonalDev, "personal_development"},
		{DomainPhotography, "photography"},
		{DomainMusic, "music"},
		{DomainHealthFitness, "health_fitness"},
		{DomainLanguage, "language"},
	}

	for _, tc := range domains {
		t.Run("returns correct string for "+tc.expected, func(t *testing.T) {
			if tc.domain.String() != tc.expected {
				t.Errorf("expected '%s', got '%s'", tc.expected, tc.domain.String())
			}
		})
	}
}

func TestNewDomainFromString(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"programming", "programming"},
		{"design", "design"},
		{"business", "business"},
		{"marketing", "marketing"},
		{"data_science", "data_science"},
		{"personal_development", "personal_development"},
		{"photography", "photography"},
		{"music", "music"},
		{"health_fitness", "health_fitness"},
		{"language", "language"},
	}

	for _, tc := range testCases {
		t.Run("successfully creates domain from "+tc.input, func(t *testing.T) {
			domain, err := NewDomainFromString(tc.input)

			if err != nil {
				t.Fatalf("expected no error for '%s', got %v", tc.input, err)
			}
			if domain.String() != tc.expected {
				t.Errorf("expected '%s', got '%s'", tc.expected, domain.String())
			}
		})
	}

	t.Run("fails for invalid domain string", func(t *testing.T) {
		domain, err := NewDomainFromString("invalid_domain")

		if err == nil {
			t.Fatal("expected error for invalid domain, got nil")
		}
		if domain.String() != "" {
			t.Error("expected empty domain for failed creation")
		}
		expectedError := "unknown 'invalid_domain' domain"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("fails for empty string", func(t *testing.T) {
		_, err := NewDomainFromString("")

		if err == nil {
			t.Fatal("expected error for empty string, got nil")
		}
	})

	t.Run("is case sensitive", func(t *testing.T) {
		_, err := NewDomainFromString("Programming")

		if err == nil {
			t.Fatal("expected error for uppercase 'Programming', got nil")
		}
	})
}

func TestTag_String(t *testing.T) {
	tags := []struct {
		tag      Tag
		expected string
	}{
		{TagBackend, "backend"},
		{TagFrontend, "frontend"},
		{TagFullStack, "fullstack"},
		{TagMobile, "mobile"},
		{TagDevOps, "devops"},
		{TagDatabase, "database"},
		{TagSecurity, "security"},
		{TagTesting, "testing"},
		{TagAPI, "api"},
		{TagCloud, "cloud"},
		{TagAI, "ai"},
		{TagMachineLearning, "machine_learning"},
		{TagWebDev, "web_development"},
		{TagGameDev, "game_development"},
		{TagBeginner, "beginner_friendly"},
		{TagAdvanced, "advanced"},
		{TagCertified, "certified"},
		{TagFree, "free"},
		{TagPaid, "paid"},
	}

	for _, tc := range tags {
		t.Run("returns correct string for "+tc.expected, func(t *testing.T) {
			if tc.tag.String() != tc.expected {
				t.Errorf("expected '%s', got '%s'", tc.expected, tc.tag.String())
			}
		})
	}
}

func TestNewTagFromString(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"backend", "backend"},
		{"frontend", "frontend"},
		{"fullstack", "fullstack"},
		{"mobile", "mobile"},
		{"devops", "devops"},
		{"database", "database"},
		{"security", "security"},
		{"testing", "testing"},
		{"api", "api"},
		{"cloud", "cloud"},
		{"ai", "ai"},
		{"machine_learning", "machine_learning"},
		{"web_development", "web_development"},
		{"game_development", "game_development"},
		{"beginner_friendly", "beginner_friendly"},
		{"advanced", "advanced"},
		{"certified", "certified"},
		{"free", "free"},
		{"paid", "paid"},
	}

	for _, tc := range testCases {
		t.Run("successfully creates tag from "+tc.input, func(t *testing.T) {
			tag, err := NewTagFromString(tc.input)

			if err != nil {
				t.Fatalf("expected no error for '%s', got %v", tc.input, err)
			}
			if tag.String() != tc.expected {
				t.Errorf("expected '%s', got '%s'", tc.expected, tag.String())
			}
		})
	}

	t.Run("fails for invalid tag string", func(t *testing.T) {
		tag, err := NewTagFromString("invalid_tag")

		if err == nil {
			t.Fatal("expected error for invalid tag, got nil")
		}
		if tag.String() != "" {
			t.Error("expected empty tag for failed creation")
		}
		expectedError := "unknown 'invalid_tag' tag"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("fails for empty string", func(t *testing.T) {
		_, err := NewTagFromString("")

		if err == nil {
			t.Fatal("expected error for empty string, got nil")
		}
	})

	t.Run("is case sensitive", func(t *testing.T) {
		_, err := NewTagFromString("Backend")

		if err == nil {
			t.Fatal("expected error for uppercase 'Backend', got nil")
		}
	})
}
