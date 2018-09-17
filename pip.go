package reqs

import (
	"bufio"
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func GetPipRequirements(dirPath string, recurse bool) (text string) {
	const reqsYml = "reqs.yml"
	const pipRequirements = "requirements.txt"
	// possible variation
	const pipRequirementsDarwin = "requirements-osx.txt"
	fileNames := GetRequirementFilenames(dirPath, recurse)

	for _, fname := range fileNames {
		if runtime.GOOS == "darwin" {
			if strings.HasSuffix(fname, pipRequirementsDarwin) {
				log.Info("Found " + fname)
				b, err := ioutil.ReadFile(fname)
				FatalCheck(err)
				text = AppendNewLinesOnly(text, string(b))
			}
		}
		if strings.HasSuffix(fname, pipRequirements) && !strings.HasSuffix(fname, "-"+pipRequirements) {
			log.Info("Found " + fname)
			b, err := ioutil.ReadFile(fname)
			FatalCheck(err)
			text = AppendNewLinesOnly(text, string(b))
		} else if strings.Contains(fname, reqsYml) {
			log.Info("Found " + fname)
			conf := ymlToMap(fname)
			for tool, packages := range conf {
				if tool == "pip" {
					for _, p := range packages {
						text = AppendNewLinesOnly(text, string(p))
					}
				}
			}
		}
	}
	return strings.TrimSpace(strings.Replace(text, "\n", " ", -1))
}

func GetPip3Requirements(dirPath string, recurse bool) (text string) {
	const reqsYml = "reqs.yml"
	const pipRequirements = "requirements.txt"
	// possible variation
	const pipRequirementsDarwin = "requirements-osx.txt"
	fileNames := GetRequirementFilenames(dirPath, recurse)

	for _, fname := range fileNames {
		if runtime.GOOS == "darwin" {
			if strings.HasSuffix(fname, pipRequirementsDarwin) {
				log.Info("Found " + fname)
				b, err := ioutil.ReadFile(fname)
				FatalCheck(err)
				text = AppendNewLinesOnly(text, string(b))
			}
		}
		if strings.HasSuffix(fname, pipRequirements) && !strings.HasSuffix(fname, "-"+pipRequirements) {
			log.Info("Found " + fname)
			b, err := ioutil.ReadFile(fname)
			FatalCheck(err)
			text = AppendNewLinesOnly(text, string(b))
		} else if strings.Contains(fname, reqsYml) {
			log.Info("Found " + fname)
			conf := ymlToMap(fname)
			for tool, packages := range conf {
				if tool == "pip3" {
					for _, p := range packages {
						text = AppendNewLinesOnly(text, string(p))
					}
				}
			}
		}
	}
	return strings.TrimSpace(strings.Replace(text, "\n", " ", -1))
}

func GetPipRequirementsMultipleDirs(dirPaths []string, recurse bool) (reqs string) {
	for _, dirPath := range dirPaths {
		reqs = NewLineIfNotEmpty(reqs, GetPipRequirements(dirPath, recurse))
	}
	return reqs
}

func GetPip3RequirementsMultipleDirs(dirPaths []string, recurse bool) (reqs string) {
	for _, dirPath := range dirPaths {
		reqs = NewLineIfNotEmpty(reqs, GetPip3Requirements(dirPath, recurse))
	}
	return reqs
}

// pip install given requirements, optionally --upgrade as well
func PipInstall(requirements, pipPath string, sudo, upgrade, quiet bool) {
	// because pip requirements.txt files can be more complicated than the
	// cli accepts with args, we write out the requirements to a temporary
	// file and then pass the file with -r to pip to read
	log.Info("Installing " + pipPath + " requirements to currently active environment")
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

	tmpReqsFile, err := ioutil.TempFile("/tmp", "reqs-")
	FatalCheck(err)
	defer os.Remove(tmpReqsFile.Name())

	reqLines := strings.Split(requirements, " ")
	w := bufio.NewWriter(tmpReqsFile)

	for _, line := range reqLines {
		_, err := w.WriteString(line + "\n")
		FatalCheck(err)
	}
	w.Flush()

	cmdStr := sudoArg + pipPath + " install " + upgradeArg + quietArg + "-r " + tmpReqsFile.Name()
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
	err = cmd.Run()
	if !quiet {
		fmt.Print(string(out.String()))
	}
	if err != nil {
		log.Fatal(stderr.String())
	}
}
