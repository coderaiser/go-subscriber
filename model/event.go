package model

// Event represents a single go test -json event.
type Event struct {
	Action  string
	Package string
	Test    string
	Output  string
	Elapsed float64
}
