package agent

import "errors"

var (
	// ErrMissingApplication is returned when an application has not been
	// provided.
	ErrMissingApplication = errors.New("missing application")

	// ErrMissingApplicationType is returned when an application type has
	// not been provided.
	ErrMissingApplicationType = errors.New("missing application type")

	// ErrMissingIdentifier is returned when an identifier has not been
	// provided and the application type is not custom.
	ErrMissingIdentifier = errors.New("missing identifier")

	// ErrMissingInstructions is returned if no instructions have been
	// configured for the agent.
	ErrMissingInstructions = errors.New("missing instructions")
)
