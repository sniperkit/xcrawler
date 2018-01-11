package model

// Course store openedx course data
type Course struct {
	CourseID  string
	Run       string
	Name      string
	Number    string
	StartDate *time.Time
	EndDate   *time.Time
	URL       string
}
