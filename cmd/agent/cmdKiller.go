package main

import (
	"os"
	"os/exec"
)

// cmdKiller is an implementation of the Killer interface, which
// operates using the command line.
type cmdKiller struct{}

// newCmdKiller returns a pointer to a new instance of cmdKiller.
func newCmdKiller() (k *cmdKiller) {
	return new(cmdKiller)
}

// Kill performs a kill operation on the command line, given a
// set of instructions.
func (k *cmdKiller) Kill(instructions []string) (err error) {
	name := instructions[0]

	// allow for single commands and commands with arguments to
	// be executed (1: will fail for single commands).
	var args []string
	if len(instructions) > 1 {
		args = instructions[1:]
	}

	// #nosec - the purpose of this killer is to execute
	// subprocesses with variables...
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
