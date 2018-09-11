package reqs

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

type VagrantSystem struct {
	Arch, Status string
}

func (vm VagrantSystem) Up() {
	cmdStr := "vagrant up " + vm.Arch
	log.Info(cmdStr)
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			log.Info(scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
	vm.Status = "up"
}

func (vm VagrantSystem) Down() {
	cmdStr := "vagrant destroy " + vm.Arch + " -f"
	log.Info(cmdStr)
	out, err := exec.Command("/bin/sh", "-c", cmdStr).Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(string(out))
	vm.Status = "down"
}

func (vm VagrantSystem) Run(cmdStr string) {
	if vm.Status != "up" {
		vm.Up()
	}
	cmdStr = "vagrant ssh " + vm.Arch + " -c \"" + cmdStr + "\""
	log.Info(cmdStr)
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			log.Info(scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
