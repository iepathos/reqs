package reqs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testReqsOn(arch, reqsArgsStr string) error {
	vm := VagrantSystem{
		Arch:   arch,
		Status: "",
	}
	vm.Down()
	err := vm.Run("reqs " + reqsArgsStr)
	vm.Down()
	return err
}

func TestReqsUbuntu(t *testing.T) {
	err := testReqsOn("ubuntu", "-r")
	assert.Nil(t, err)

	// test with update and upgrade
	err = testReqsOn("ubuntu", "-r -up")
	assert.Nil(t, err)
}

func TestReqsFedora(t *testing.T) {
	err := testReqsOn("fedora", "-r")
	assert.Nil(t, err)

	// test with update and upgrade
	err = testReqsOn("fedora", "-r -up")
	assert.Nil(t, err)
}

func TestReqsOsx(t *testing.T) {
	err := testReqsOn("osx", "-r")
	assert.Nil(t, err)

	// test with update and upgrade
	err = testReqsOn("osx", "-r -up")
	assert.Nil(t, err)
}
