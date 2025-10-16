package course

import "github.com/pkg/errors"

// CourseLevel enum
var (
	Beginner     = CourseLevel{l: "beginner"}
	Intermediate = CourseLevel{l: "intermediate"}
	Advanced     = CourseLevel{l: "advanced"}
)

var courseLevelValues = []CourseLevel{
	Beginner,
	Intermediate,
	Advanced,
}

type CourseLevel struct {
	l string
}

func (l CourseLevel) String() string {
	return l.l
}

func NewCourseLevelFromString(levelStr string) (CourseLevel, error) {
	for _, level := range courseLevelValues {
		if level.String() == levelStr {
			return level, nil
		}
	}
	return CourseLevel{}, errors.Errorf("unknown '%s' course level", levelStr)
}

// Domain enum
var (
	DomainProgramming   = Domain{d: "programming"}
	DomainDesign        = Domain{d: "design"}
	DomainBusiness      = Domain{d: "business"}
	DomainMarketing     = Domain{d: "marketing"}
	DomainDataScience   = Domain{d: "data_science"}
	DomainPersonalDev   = Domain{d: "personal_development"}
	DomainPhotography   = Domain{d: "photography"}
	DomainMusic         = Domain{d: "music"}
	DomainHealthFitness = Domain{d: "health_fitness"}
	DomainLanguage      = Domain{d: "language"}
)

var domainValues = []Domain{
	DomainProgramming,
	DomainDesign,
	DomainBusiness,
	DomainMarketing,
	DomainDataScience,
	DomainPersonalDev,
	DomainPhotography,
	DomainMusic,
	DomainHealthFitness,
	DomainLanguage,
}

type Domain struct {
	d string
}

func (d Domain) String() string {
	return d.d
}

func NewDomainFromString(domainStr string) (Domain, error) {
	for _, domain := range domainValues {
		if domain.String() == domainStr {
			return domain, nil
		}
	}
	return Domain{}, errors.Errorf("unknown '%s' domain", domainStr)
}

// Tag enum
var (
	TagBackend     = Tag{t: "backend"}
	TagFrontend    = Tag{t: "frontend"}
	TagFullStack   = Tag{t: "fullstack"}
	TagMobile      = Tag{t: "mobile"}
	TagDevOps      = Tag{t: "devops"}
	TagDatabase    = Tag{t: "database"}
	TagSecurity    = Tag{t: "security"}
	TagTesting     = Tag{t: "testing"}
	TagAPI         = Tag{t: "api"}
	TagCloud       = Tag{t: "cloud"}
	TagAI          = Tag{t: "ai"}
	TagMachineLearning = Tag{t: "machine_learning"}
	TagWebDev      = Tag{t: "web_development"}
	TagGameDev     = Tag{t: "game_development"}
	TagBeginner    = Tag{t: "beginner_friendly"}
	TagAdvanced    = Tag{t: "advanced"}
	TagCertified   = Tag{t: "certified"}
	TagFree        = Tag{t: "free"}
	TagPaid        = Tag{t: "paid"}
)

var tagValues = []Tag{
	TagBackend,
	TagFrontend,
	TagFullStack,
	TagMobile,
	TagDevOps,
	TagDatabase,
	TagSecurity,
	TagTesting,
	TagAPI,
	TagCloud,
	TagAI,
	TagMachineLearning,
	TagWebDev,
	TagGameDev,
	TagBeginner,
	TagAdvanced,
	TagCertified,
	TagFree,
	TagPaid,
}

type Tag struct {
	t string
}

func (t Tag) String() string {
	return t.t
}

func NewTagFromString(tagStr string) (Tag, error) {
	for _, tag := range tagValues {
		if tag.String() == tagStr {
			return tag, nil
		}
	}
	return Tag{}, errors.Errorf("unknown '%s' tag", tagStr)
}
