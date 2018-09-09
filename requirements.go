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

// requirements for parsing requirements files
// and for determining currently installed requirements

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

func recurseForFiles(dir string, fnames []string) (filePaths []string) {
	filepathList := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		filepathList = append(filepathList, path)
		return nil
	})
	FatalCheck(err)

	for _, path := range filepathList {
		if StringContainedInSlice(path, fnames) {
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
			text = AppendNewLinesOnly(text, string(b))
		} else if strings.Contains(fname, reqsYml) {
			log.Info("Found " + fname)
			conf := ymlToMap(fname)
			for tool, packages := range conf {
				if tool == "common" || tool == packageTool {
					for _, p := range packages {
						text = AppendNewLinesOnly(text, string(p))
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

func getPipRequirements(dirPath string, recurse bool) (text string) {
	const reqsYml = "reqs.yml"
	const pipRequirements = "requirements.txt"
	// possible variation
	const pipRequirementsDarwin = "requirements-osx.txt"
	fileNames := getRequirementFilenames(dirPath, recurse)

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

func getPipRequirementsMultipleDirs(dirPaths []string, recurse bool) (reqs string) {
	for _, dirPath := range dirPaths {
		reqs = NewLineIfNotEmpty(reqs, getPipRequirements(dirPath, recurse))
	}
	return reqs
}

func getSysRequirementsMultipleDirs(dirPaths []string, packageTool string, recurse bool) (reqs string) {
	for _, dirPath := range dirPaths {
		reqs = NewLineIfNotEmpty(reqs, getSysRequirements(dirPath, packageTool, recurse))
	}
	return reqs
}

func aptListInstalled(withVersion bool) (reqs string) {
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
			reqs = NewLineIfNotEmpty(reqs, req)
		}
	}
	return reqs
}

func brewListInstalled() string {
	out, err := exec.Command("brew", "list").Output()
	FatalCheck(err)
	return strings.TrimSpace(string(out))
}

func dnfListInstalled(withVersion bool) (reqs string) {
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
			reqs = NewLineIfNotEmpty(reqs, req)
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

func (rp RequirementsParser) ListInstalled(packageTool string) (requirements string) {
	switch packageTool {
	case "apt":
		requirements = aptListInstalled(rp.WithVersion)
	case "brew":
		requirements = brewListInstalled()
	case "dnf":
		requirements = dnfListInstalled(rp.WithVersion)
	}
	return requirements
}

func amIRoot() bool {
	if os.Geteuid() == 0 {
		return true
	}
	return false
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
			if IsCommandAvailable(tool) {
				packageTool = tool
				break
			}
		}
		if !amIRoot() {
			sudo = "sudo "
		}
		autoYes = "-y "
	case "darwin":
		if !rp.UseStdout {
			log.Info("Darwin system detected")
		}
		if !IsCommandAvailable("brew") {
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
		reqs = rp.ListInstalled(packageTool)
		fmt.Print(reqs)
		os.Exit(0)
	} else {
		// parse the current directory
		reqs = getSysRequirements(".", packageTool, rp.Recurse)
	}
	reqs = strings.TrimSpace(strings.Replace(reqs, "\n", " ", -1))
	return sudo, packageTool, autoYes, reqs
}

func (rp RequirementsParser) ParsePip() (reqs string) {

	if rp.Dir != "" {
		// search directory for requirements
		if strings.Contains(rp.Dir, ",") {
			reqs = getPipRequirementsMultipleDirs(strings.Split(rp.Dir, ","), rp.Recurse)
		} else {
			reqs = getPipRequirements(rp.Dir, rp.Recurse)
		}
	} else if rp.File != "" {
		// read specified file for requirements
		b, err := ioutil.ReadFile(rp.File)
		FatalCheck(err)
		reqs = string(b)
	} else {
		// parse the current directory
		reqs = getPipRequirements(".", rp.Recurse)
	}
	return reqs
}
