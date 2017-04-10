package orchestrator

// Application represents an application the Orchestrator can kill.
type Application struct {
	Name       string  `json:"name"`
	Schedule   string  `json:"schedule"`
	Percentage float64 `json:"percentage"`
}

// Applications is a slice of Application structs.
type Applications []Application
