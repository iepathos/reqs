package reqs

import (
	"testing"
)

func testReqsOn(arch, reqsArgsStr string) {
	vm := VagrantSystem{
		Arch:   arch,
		Status: "down",
	}
	vm.Run("reqs " + reqsArgsStr)
	vm.Down()
}

func TestReqsUbuntu(t *testing.T) {
	testReqsOn("ubuntu", "-r")
}

func TestReqsFedora(t *testing.T) {
	testReqsOn("fedora", "-r")
}

func TestReqsOsx(t *testing.T) {
	testReqsOn("osx", "-r")
}
