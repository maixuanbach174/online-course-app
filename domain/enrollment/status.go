package enrollment

var (
	Enrolled   = Status{s: "enrolled"}
	Started    = Status{s: "started"}
	InProgress = Status{s: "in_progress"}
	Completed  = Status{s: "completed"}
)

type Status struct {
	s string
}

func (s Status) String() string {
	return s.s
}

func (s Status) NewStatusFromString(statusStr string) (Status, error) {
	panic("not implemented")
}
