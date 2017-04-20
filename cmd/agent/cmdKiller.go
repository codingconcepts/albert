package main

import (
	"os"
	"os/exec"
)

type cmdKiller struct{}

func newCmdKiller() (k *cmdKiller) {
	return new(cmdKiller)
}

// Kill performs a kill operation on the command line.
func (k *cmdKiller) Kill(instructions []string) (err error) {
	name := instructions[0]

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
