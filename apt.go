package reqs

import (
	"io/ioutil"
	"os/exec"
	"strings"
)

func GetAptSources() (out string) {
	b, err := ioutil.ReadFile("/etc/apt/sources.list")
	FatalCheck(err)
	// clean empty lines and comments out
	for _, line := range strings.Split(string(b), "\n") {
		if !strings.HasPrefix(line, "#") && strings.TrimSpace(line) != "" {
			if out == "" {
				out = strings.TrimSpace(line)
			} else {
				out += "\n" + strings.TrimSpace(line)
			}
		}
	}
	return out
}

func AptListInstalled(withVersion bool) (reqs string) {
	out, err := exec.Command("apt", "list", "--installed").Output()
	FatalCheck(err)
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "/") {
			lSplit := strings.Split(string(line), "/")
			req := lSplit[0]
			if withVersion {
				version := strings.Split(lSplit[1], " ")[1]
				req = req + "=" + version
			}
			reqs = NewLineIfNotEmpty(reqs, req)
		}
	}
	return strings.TrimSpace(reqs)
}
