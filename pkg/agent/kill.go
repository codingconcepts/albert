package agent

import (
	"os"
	"os/exec"

	"github.com/Sirupsen/logrus"
)

func (a *Agent) killSimulation() {
	a.Logger.WithFields(logrus.Fields{
		"application": a.Application,
	}).Info("kill simulation")
}

// http://stackoverflow.com/a/32074098/304957
func (a *Agent) killContainer(name string) (err error) {
	a.Logger.WithFields(logrus.Fields{
		"container": name,
	}).Info("killing container")

	cmd := exec.Command(`$(docker stop $(docker ps -a -q --filter ancestor=%s --format="{{.ID}}"))`, name)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (a *Agent) killCustom() (err error) {
	name := a.CustomInstructions[0]

	var args []string
	if len(a.CustomInstructions) > 1 {
		args = a.CustomInstructions[1:]
	}

	a.Logger.WithFields(logrus.Fields{
		"name": name,
		"args": args,
	}).Info("custom killing")

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
