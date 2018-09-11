package reqs

import (
	log "github.com/sirupsen/logrus"
	"os/exec"
	"testing"
)

func TestReqsUbuntu(t *testing.T) {
	system := "ubuntu"
	cmdStr := "vagrant up " + system
	log.Info(cmdStr)
	out, err := exec.Command("/bin/sh", "-c", cmdStr).Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(string(out))

	cmdStr = "vagrant ssh " + system + " -c \"reqs -r -spip -snpm\""
	log.Info(cmdStr)
	out, err = exec.Command("/bin/sh", "-c", cmdStr).Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(string(out))

	cmdStr = "vagrant destroy " + system + " -f"
	log.Info(cmdStr)
	out, err = exec.Command("/bin/sh", "-c", cmdStr).Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(string(out))
}

func TestReqsFedora(t *testing.T) {
	system := "fedora"
	cmdStr := "vagrant up " + system
	log.Info(cmdStr)
	out, err := exec.Command("/bin/sh", "-c", cmdStr).Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(string(out))

	cmdStr = "vagrant ssh " + system + " -c \"reqs -r -spip -snpm\""
	log.Info(cmdStr)
	out, err = exec.Command("/bin/sh", "-c", cmdStr).Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(string(out))

	cmdStr = "vagrant destroy " + system + " -f"
	log.Info(cmdStr)
	out, err = exec.Command("/bin/sh", "-c", cmdStr).Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(string(out))
}

func TestReqsDarwin(t *testing.T) {
	system := "osx"
	cmdStr := "vagrant up " + system
	log.Info(cmdStr)
	out, err := exec.Command("/bin/sh", "-c", cmdStr).Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(string(out))

	cmdStr = "vagrant ssh " + system + " -c \"reqs -r -spip -npm\""
	log.Info(cmdStr)
	out, err = exec.Command("/bin/sh", "-c", cmdStr).Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(string(out))

	cmdStr = "vagrant destroy " + system + " -f"
	log.Info(cmdStr)
	out, err = exec.Command("/bin/sh", "-c", cmdStr).Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(string(out))
}
