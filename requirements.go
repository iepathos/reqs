package reqs

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// requirements parsing and generating functions
func isCommandAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func newLineIfNotEmpty(text, newText string) string {
	if text == "" {
		text = newText
	} else {
		text += "\n" + newText
	}
	return text
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func appendNewLinesOnly(text, newText string) string {
	textSplit := strings.Split(text, "\n")
	newTextSplit := strings.Split(newText, "\n")
	returnText := text
	for _, line := range newTextSplit {
		if !stringInSlice(line, textSplit) {
			returnText += "\n" + line
		}
	}
	return returnText
}

func installHomebrew() {
	log.Info("Installing homebrew")
	cmd := exec.Command("/usr/bin/ruby",
		"-e",
		"\"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)\"")
	err := cmd.Run()
	FatalCheck(err)
}

func getAptSources() (sources string) {
	b, err := ioutil.ReadFile("/etc/apt/sources.list")
	FatalCheck(err)
	sources += "\n" + string(b)
	return sources
}

func getBrewTaps() string {
	out, err := exec.Command("brew", "tap").Output()
	FatalCheck(err)
	return strings.TrimSpace(string(out))
}

func stringContainedInSlice(s string, arr []string) bool {
	for _, v := range arr {
		if strings.Contains(s, v) {
			return true
		}
	}
	return false
}

func recurseForFiles(dir string, fnames []string) (filePaths []string) {
	filepathList := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		filepathList = append(filepathList, path)
		return nil
	})
	FatalCheck(err)

	for _, path := range filepathList {
		if stringContainedInSlice(path, fnames) {
			filePaths = append(filePaths, path)
		}
	}
	return filePaths
}

func getRequirementFilenames(dirPath string, recurse bool) (fileNames []string) {
	requirementFilenames := []string{
		"requirements.txt",
		"reqs.yml",
	}
	if recurse {
		fileNames = recurseForFiles(dirPath, requirementFilenames)
	} else {
		files, err := ioutil.ReadDir(dirPath)
		FatalCheck(err)
		for _, f := range files {
			fileNames = append(fileNames, dirPath+"/"+f.Name())
		}
	}
	return fileNames
}

func ymlToMap(ymlPath string) (conf map[string][]string) {
	b, err := ioutil.ReadFile(ymlPath)
	FatalCheck(err)
	m := make(map[string][]string)
	err = yaml.Unmarshal(b, &m)
	FatalCheck(err)
	return m
}

// find tool-requirements.txt, common-requirements.txt and/or reqs.yml
// in the specified directory, can recurse down the directory
func getSysRequirements(dirPath, packageTool string, recurse bool) (text string) {
	fileNames := getRequirementFilenames(dirPath, recurse)
	toolRequirements := packageTool + "-requirements.txt"
	const commonRequirements = "common-requirements.txt"
	const reqsYml = "reqs.yml"

	for _, fname := range fileNames {
		if strings.Contains(fname, commonRequirements) || strings.Contains(fname, toolRequirements) {
			log.Info("Found " + fname)
			b, err := ioutil.ReadFile(fname)
			FatalCheck(err)
			text = appendNewLinesOnly(text, string(b))
		} else if strings.Contains(fname, reqsYml) {
			log.Info("Found " + fname)
			conf := ymlToMap(fname)
			for tool, packages := range conf {
				if tool == "common" || tool == packageTool {
					for _, p := range packages {
						text = appendNewLinesOnly(text, string(p))
					}
				}
			}
		}
	}
	if len(text) == 0 {
		log.Fatal("No requirements files found")
	}
	return strings.TrimSpace(text)
}

func getSysRequirementsMultipleDirs(dirPaths []string, packageTool string, recurse bool) (reqs string) {
	for _, dirPath := range dirPaths {
		reqs = newLineIfNotEmpty(reqs, getSysRequirements(dirPath, packageTool, recurse))
	}
	return reqs
}

func getInstalledAptRequirements(withVersion bool) (reqs string) {
	out, err := exec.Command("sudo", "apt", "list", "--installed").Output()
	FatalCheck(err)
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "/") {
			lSplit := strings.Split(string(line), "/")
			req := lSplit[0]
			if withVersion {
				version := strings.Split(lSplit[1], " ")[1]
				req = req + "=" + version
			}
			reqs = newLineIfNotEmpty(reqs, req)
		}
	}
	return reqs
}

func getInstalledBrewRequirements() string {
	out, err := exec.Command("brew", "list").Output()
	FatalCheck(err)
	return strings.TrimSpace(string(out))
}

func getInstalledDnfRequirements(withVersion bool) (reqs string) {
	out, err := exec.Command("sudo", "dnf", "list", "installed").Output()
	FatalCheck(err)
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "@System") {
			lSplit := strings.Split(string(line), " ")
			req := lSplit[0]
			if withVersion {
				version := strings.Split(lSplit[1], " ")[0]
				req = req + "=" + version
			}
			reqs = newLineIfNotEmpty(reqs, req)
		}
	}
	return reqs
}

type RequirementsParser struct {
	Dir, File           string
	UseStdout, UseStdin bool
	WithVersion         bool
	Recurse             bool
	Sources             bool
}

// determine the package tool, sudo and autoYes based on the current system
func (rp RequirementsParser) parseTooling() (sudo, packageTool, autoYes string) {
	switch runtime.GOOS {
	case "linux":
		if !rp.UseStdout {
			log.Info("Linux system detected")
		}
		linuxTools := []string{
			"apt",
			"dnf",
		}
		for _, tool := range linuxTools {
			if isCommandAvailable(tool) {
				packageTool = tool
				break
			}
		}
		sudo = "sudo "
		autoYes = "-y "
	case "darwin":
		if !rp.UseStdout {
			log.Info("Darwin system detected")
		}
		if !isCommandAvailable("brew") {
			installHomebrew()
		}
		packageTool = "brew"
	case "windows":
		log.Fatal("Windows system detected, abandon all hope")
		os.Exit(1)
	}

	return sudo, packageTool, autoYes
}

// determine package tool and args on this system
func (rp RequirementsParser) Parse() (sudo, packageTool, autoYes, reqs string) {
	sudo, packageTool, autoYes = rp.parseTooling()
	// output sources for apt, taps for brew
	if rp.Sources {
		switch packageTool {
		case "apt":
			fmt.Print(getAptSources())
		case "brew":
			fmt.Print(getBrewTaps())
		}
		os.Exit(0)
	}

	if rp.Dir != "" {
		// search directory for requirements
		if strings.Contains(rp.Dir, ",") {
			reqs = getSysRequirementsMultipleDirs(strings.Split(rp.Dir, ","), packageTool, rp.Recurse)
		} else {
			reqs = getSysRequirements(rp.Dir, packageTool, rp.Recurse)
		}
	} else if rp.File != "" {
		// read specified file for requirements
		b, err := ioutil.ReadFile(rp.File)
		FatalCheck(err)
		reqs = string(b)
	} else if rp.UseStdin {
		// read stdin for requirements
		reader := bufio.NewReader(os.Stdin)
		reqs, _ = reader.ReadString('\n')
	} else if rp.UseStdout {
		// output requirements to stdout
		switch packageTool {
		case "apt":
			reqs = getInstalledAptRequirements(rp.WithVersion)
		case "brew":
			reqs = getInstalledBrewRequirements()
		case "dnf":
			reqs = getInstalledDnfRequirements(rp.WithVersion)
		}
		fmt.Print(reqs)
		os.Exit(0)
	} else {
		// parse the current directory
		reqs = getSysRequirements(".", packageTool, rp.Recurse)
	}
	reqs = strings.TrimSpace(strings.Replace(reqs, "\n", " ", -1))
	return sudo, packageTool, autoYes, reqs
}
