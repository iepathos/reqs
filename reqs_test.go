package reqs

import (
	log "github.com/sirupsen/logrus"
	"os/exec"
	"testing"
)

func TestReqsUbuntu(t *testing.T) {
	out, err := exec.Command("/bin/sh", "-c", "vagrant up ubuntu").Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(string(out))

	out, err := exec.Command("/bin/sh", "-c", "vagrant destroy ubuntu -f").Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(string(out))
}

func TestReqsFedora(t *testing.T) {
	out, err := exec.Command("/bin/sh", "-c", "vagrant up fedora").Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(string(out))

	out, err := exec.Command("/bin/sh", "-c", "vagrant destroy fedora -f").Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(string(out))
}
