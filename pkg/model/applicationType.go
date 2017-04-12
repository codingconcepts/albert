package model

// ApplicationType is an enum value which describes the
// various application types it's possible to kill.
type ApplicationType string

const (
	// UnknownApplicationType is the default value and results in
	// a validation error.
	UnknownApplicationType ApplicationType = ""

	// DummyApplicationType can be assigned to an application
	// to test its connectivity without actually killing anything.
	DummyApplicationType ApplicationType = "dummy"

	// ProcessApplicationType is used to kill individual
	// processes running on machines.
	ProcessApplicationType ApplicationType = "process"

	// MachineApplicationType is used to kill entire machines.
	MachineApplicationType ApplicationType = "machine"

	// DockerApplicationType is used to kill a docker image.
	DockerApplicationType ApplicationType = "docker"

	// NetworkApplicationType is used to kill a machine's
	// network interface, leaving all processes running, just
	// unable to talk to the outside world.
	NetworkApplicationType ApplicationType = "network"

	// CustomApplicationType can be used to provide custom
	// kill logic, configurable to kill any of the above types.
	CustomApplicationType ApplicationType = "custom"
)
