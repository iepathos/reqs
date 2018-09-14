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

// test basic apt
func TestReqsApt(t *testing.T) {
	err := testReqsOn("ubuntu", "-r")
	assert.Nil(t, err)
}

// also test pip3 with update and upgrade
func TestReqsUbuntuPip3(t *testing.T) {
	// test with pip and pip3 update and upgrade
	err := testReqsOn("ubuntu", "-d pip3-setup -up -spip3")
	assert.Nil(t, err)
}

func TestReqsUbuntuNpm(t *testing.T) {
	err := testReqsOn("ubuntu", "-d web-service -snpm")
	assert.Nil(t, err)
}

// test basic dnf
func TestReqsDnf(t *testing.T) {
	err := testReqsOn("fedora", "-r")
	assert.Nil(t, err)
}

// also tests update and upgrade
func TestReqsFedoraPip3(t *testing.T) {
	err := testReqsOn("fedora", "-d pip3-setup -up -spip3")
	assert.Nil(t, err)
}

// test basic brew
func TestReqsBrew(t *testing.T) {
	err := testReqsOn("osx", "-r")
	assert.Nil(t, err)
}

// test osx with npm
func TestReqsOsxNpm(t *testing.T) {
	err := testReqsOn("osx", "-d web-service -npm")
	assert.Nil(t, err)
}

// test osx with pip3
// func TestReqsOsxPip3(t *testing.T) {
// 	err := testReqsOn("osx", "-d pip3-service -pip3")
// 	assert.Nil(t, err)
// }
