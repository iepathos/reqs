package reqs

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

// responsible for interfacing with package tools
// deals with apt, brew, and dnf, pip

func runShell(code string) {
	log.Info(code)
	cmd := exec.Command("/bin/sh", "-c", code)
	err := cmd.Run()
	FatalCheck(err)
}

type PackageConfig struct {
	Tool          string
	Sudo, AutoYes string
	Reqs          string
	Quiet, Force  bool
}

func (pc PackageConfig) getForceArg() (forceArg string) {
	if pc.Force {
		if pc.Tool == "brew" {
			forceArg = "--force "
		} else {
			forceArg = "-f "
		}
	}
	return forceArg
}

func (pc PackageConfig) Install(upgrade bool) {
	log.Info("Installing system requirements with " + pc.Tool)
	forceArg := pc.getForceArg()
	envArg := ""
	if pc.Tool == "brew" {
		envArg += "HOMEBREW_NO_AUTO_UPDATE=1 "
	}
	upgradeArg := ""
	if upgrade {
		if pc.Tool == "apt" {
			upgradeArg = "--upgrade "
		}
	}
	cmdStr := envArg + pc.Sudo + pc.Tool + " install " + pc.AutoYes + forceArg + upgradeArg + pc.Reqs
	log.Info(cmdStr)
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if !pc.Quiet {
		fmt.Print(string(out.String()))
	}
	if err != nil {
		log.Fatal(stderr.String())
	}
}

func (pc PackageConfig) abstractUp(upArg string) {
	log.Info("Running " + pc.Tool + " packages " + upArg)
	forceArg := pc.getForceArg()
	if pc.Tool == "brew" {
		runShell(pc.Tool + " " + upArg + " " + forceArg)
	} else {
		runShell("sudo " + pc.Tool + " " + upArg + " " + forceArg + pc.AutoYes)
	}
}

func (pc PackageConfig) Update() {
	pc.abstractUp("update")
}

func (pc PackageConfig) Upgrade() {
	pc.abstractUp("upgrade")
}
