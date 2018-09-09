package reqs

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
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

// pip install given requirements, optionally --upgrade as well
func PipInstall(requirements, pipPath string, sudo, upgrade, quiet bool) {
	log.Info("Installing pip requirements to currently active environment")
	sudoArg := ""
	if sudo {
		sudoArg = "sudo "
	}
	upgradeArg := ""
	if upgrade {
		upgradeArg = "--upgrade "
	}
	quietArg := ""
	if quiet {
		quietArg = "-q "
	}
	cmdStr := sudoArg + pipPath + " install " + upgradeArg + quietArg + requirements
	if !quiet {
		log.Info(cmdStr)
	}
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Env = []string{
		"PATH=" + os.ExpandEnv("$PATH"),
		"PYTHONPATH=" + os.ExpandEnv("$PYTHONPATH"),
		"PYENV_VIRTUAL_ENV=" + os.ExpandEnv("$PYENV_VIRTUAL_ENV"),
		"PYENV_VERSION=" + os.ExpandEnv("$PYENV_VERSION"),
	}
	err := cmd.Run()
	if !quiet {
		fmt.Print(string(out.String()))
	}
	if err != nil {
		log.Fatal(stderr.String())
	}
}

func NpmInstall(requirements, dir string, sudo, global, quiet bool) {
	log.Info("Installing npm requirements")
	sudoArg := ""
	if sudo {
		sudoArg = "sudo "
	}
	globalArg := ""
	if global {
		globalArg = "-g "
	}
	cmdStr := sudoArg + "npm " + globalArg + "install " + requirements
	log.Info(cmdStr)
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	if dir != "" {
		cmd.Dir = dir
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Env = []string{
		"PATH=" + os.ExpandEnv("$PATH"),
	}
	err := cmd.Run()
	if !quiet {
		fmt.Print(string(out.String()))
	}
	if err != nil {
		log.Fatal(stderr.String())
	}
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
	envArg := ""
	if pc.Tool == "brew" {
		envArg += "HOMEBREW_NO_AUTO_UPDATE=1 "
	}
	cmdStr := envArg + pc.Sudo + pc.Tool + " install " + pc.AutoYes + forceArg + pc.Reqs
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
