package reqs

import (
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

func FatalCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func IsCommandAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func NewLineIfNotEmpty(text, newText string) string {
	if text == "" {
		text = newText
	} else {
		text += "\n" + newText
	}
	return text
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func StringContainedInSlice(s string, arr []string) bool {
	for _, v := range arr {
		if strings.Contains(s, v) {
			return true
		}
	}
	return false
}

func AppendNewLinesOnly(text, newText string) (returnText string) {
	textSplit := strings.Split(text, "\n")
	newTextSplit := strings.Split(newText, "\n")
	returnText = text
	for _, line := range newTextSplit {
		if !StringInSlice(strings.TrimSpace(line), textSplit) {
			returnText += "\n" + strings.TrimSpace(line)
		}
	}
	return strings.TrimSpace(returnText)
}
