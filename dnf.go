package reqs

import (
	"os/exec"
	"strings"
)

func DnfListInstalled(withVersion bool) (reqs string) {
	out, err := exec.Command("dnf", "list", "installed").Output()
	FatalCheck(err)
	for _, line := range strings.Split(string(out), "\n")[1:] {
		lSplit := strings.Split(string(line), " ")
		req := lSplit[0]
		if withVersion {
			req = req + "=" + lSplit[1]
		}
		reqs = NewLineIfNotEmpty(reqs, req)
	}
	return strings.TrimSpace(reqs)
}
