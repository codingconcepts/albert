// +build windows

package agent

import (
	"os/exec"

	"github.com/Sirupsen/logrus"
)

func (a *Agent) killMachine() (err error) {
	a.Logger.Info("killing machine")

	cmd := exec.Command("shutdown", "-t", "0", "-r", "-f")
	return cmd.Run()
}

func (a *Agent) killProcess(name string) (err error) {
	a.Logger.WithFields(logrus.Fields{
		"process": name,
	}).Info("killing process")

	cmd := exec.Command("taskkill", "/f", "/t", "/im", name)
	return cmd.Run()
}
