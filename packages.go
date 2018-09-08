package reqs

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

// responsible for dealing interfacing with package tools
// deal with apt, brew, and dnf

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

func (pc PackageConfig) Install() {
	log.Info("Installing system requirements with " + pc.Tool)
	forceArg := pc.getForceArg()
	cmdStr := pc.Sudo + pc.Tool + " install " + pc.AutoYes + forceArg + pc.Reqs
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
	log.Info("Updating " + pc.Tool + " packages")
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
