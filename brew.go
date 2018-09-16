package reqs

import (
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

func BrewListInstalled() string {
	out, err := exec.Command("brew", "list").Output()
	FatalCheck(err)
	return strings.TrimSpace(string(out))
}

func InstallHomebrew() {
	log.Info("Installing homebrew")
	cmd := exec.Command("/usr/bin/ruby",
		"-e",
		"\"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)\"")
	err := cmd.Run()
	FatalCheck(err)
}

func GetBrewTaps() string {
	out, err := exec.Command("brew", "tap").Output()
	FatalCheck(err)
	return strings.TrimSpace(string(out))
}
