package reqs

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func FindNpmPackageDirs(dir string, recurse bool) (packageDirs []string) {
	const packageJson = "package.json"
	if recurse {
		err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			if strings.Contains(path, packageJson) && !strings.Contains(path, "node_modules") && !strings.Contains(path, "bower_components") {
				d, _ := filepath.Split(path)
				log.Info("Found npm package directory " + d)
				packageDirs = append(packageDirs, d)
			}
			return nil
		})
		FatalCheck(err)
	} else {
		if _, err := os.Stat(dir + "/" + packageJson); !os.IsNotExist(err) {
			log.Info("Found npm package directory " + dir)
			packageDirs = append(packageDirs, dir)
		}
	}
	return packageDirs
}

func GetNpmRequirements(dir string, recurse bool) (text string) {
	const reqsYml = "reqs.yml"
	const npmRequirements = "npm-requirements.txt"
	fileNames := GetRequirementFilenames(dir, recurse)
	for _, fname := range fileNames {
		if strings.HasSuffix(fname, npmRequirements) {
			log.Info("Found " + fname)
			b, err := ioutil.ReadFile(fname)
			FatalCheck(err)
			text = AppendNewLinesOnly(text, string(b))
		} else if strings.Contains(fname, reqsYml) {
			log.Info("Found " + fname)
			conf := ymlToMap(fname)
			for tool, packages := range conf {
				if tool == "npm" {
					for _, p := range packages {
						text = AppendNewLinesOnly(text, string(p))
					}
				}
			}
		}
	}
	return strings.TrimSpace(strings.Replace(text, "\n", " ", -1))
}

func GetNpmRequirementsMultipleDirs(dirPaths []string, recurse bool) (reqs string) {
	for _, dirPath := range dirPaths {
		reqs = NewLineIfNotEmpty(reqs, GetNpmRequirements(dirPath, recurse))
	}
	return reqs
}

func NpmInstall(requirements, dir string, sudo, global, quiet bool) {
	if dir != "" {
		dir, _ = filepath.Abs(dir)
	}
	if global {
		log.Info("Installing npm global requirements")
	} else {
		logDir := ""
		if dir != "" {
			logDir = " in " + dir
		}
		log.Info("Running npm install" + logDir)
	}
	sudoArg := ""
	if sudo {
		sudoArg = "sudo "
	}
	globalArg := ""
	if global {
		globalArg = "-g "
	}
	cmdStr := sudoArg + "npm " + globalArg + "install " + requirements
	log.Info(cmdStr)
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Env = []string{
		"PATH=" + os.ExpandEnv("$PATH"),
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if !quiet {
		fmt.Print(string(out.String()))
	}
	if err != nil {
		log.Fatal(err)
		log.Fatal(stderr.String())
	}
}
