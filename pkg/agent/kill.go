package agent

import (
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
	return cmd.Run()
}
