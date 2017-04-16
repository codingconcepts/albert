package orchestrator

import "errors"

var (
	// ErrMissingApplications is returned when no applications
	// have been configured for the orchestrator
	ErrMissingApplications = errors.New("missing applications")

	// ErrMissingGatherTimeout is returned when no scatter-
	// gather timeout has been configured.
	ErrMissingGatherTimeout = errors.New("missing gather timeout")

	// ErrInvalidGatherChanSize is returned when no scatter-
	// gather channel size has been configured.
	ErrInvalidGatherChanSize = errors.New("gather channel size must be greater than zero")
)
