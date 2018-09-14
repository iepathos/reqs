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

func getAptSources() (out string) {
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
			joinArg := ""
			if !strings.HasSuffix(dirPath, "/") {
				joinArg = "/"
			}
			fileNames = append(fileNames, dirPath+joinArg+f.Name())
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

func getPip3Requirements(dirPath string, recurse bool) (text string) {
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

func getPipRequirementsMultipleDirs(dirPaths []string, recurse bool) (reqs string) {
	for _, dirPath := range dirPaths {
		reqs = NewLineIfNotEmpty(reqs, getPipRequirements(dirPath, recurse))
	}
	return reqs
}

func getPip3RequirementsMultipleDirs(dirPaths []string, recurse bool) (reqs string) {
	for _, dirPath := range dirPaths {
		reqs = NewLineIfNotEmpty(reqs, getPip3Requirements(dirPath, recurse))
	}
	return reqs
}

func findNpmPackageDirs(dir string, recurse bool) (packageDirs []string) {
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

func getNpmRequirements(dir string, recurse bool) (text string) {
	const reqsYml = "reqs.yml"
	const npmRequirements = "npm-requirements.txt"
	fileNames := getRequirementFilenames(dir, recurse)
	// TODO:
	// npmMap := make(map[string]string)
	// dir: packages
	// if just dir: "" then just run npm install there for the packages.json
	// if global: packages then install as global
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

func getNpmRequirementsMultipleDirs(dirPaths []string, recurse bool) (reqs string) {
	for _, dirPath := range dirPaths {
		reqs = NewLineIfNotEmpty(reqs, getNpmRequirements(dirPath, recurse))
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

func brewListInstalled() string {
	out, err := exec.Command("brew", "list").Output()
	FatalCheck(err)
	return strings.TrimSpace(string(out))
}

func dnfListInstalled(withVersion bool) (reqs string) {
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

type RequirementsParser struct {
	Dir, File           string
	UseStdout, UseStdin bool
	WithVersion         bool
	Recurse             bool
	Sources             bool
}

func (rp RequirementsParser) FindNpmPackageDirs() (packageDirs []string) {
	dirArg := "."
	if rp.Dir != "" {
		dirArg = rp.Dir
	}
	packageDirs = findNpmPackageDirs(dirArg, rp.Recurse)
	return packageDirs
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
		if packageTool == "" {
			log.Fatal("Failed to find any support package management tools")
			os.Exit(1)
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

func (rp RequirementsParser) ParsePip3() (reqs string) {
	if rp.Dir != "" {
		// search directory for requirements
		if strings.Contains(rp.Dir, ",") {
			reqs = getPip3RequirementsMultipleDirs(strings.Split(rp.Dir, ","), rp.Recurse)
		} else {
			reqs = getPip3Requirements(rp.Dir, rp.Recurse)
		}
	} else if rp.File != "" {
		// read specified file for requirements
		b, err := ioutil.ReadFile(rp.File)
		FatalCheck(err)
		reqs = string(b)
	} else {
		// parse the current directory
		reqs = getPip3Requirements(".", rp.Recurse)
	}
	return reqs
}

func (rp RequirementsParser) ParseNpm() (reqs string) {
	if rp.Dir != "" {
		// search directory for requirements
		if strings.Contains(rp.Dir, ",") {
			reqs = getNpmRequirementsMultipleDirs(strings.Split(rp.Dir, ","), rp.Recurse)
		} else {
			reqs = getNpmRequirements(rp.Dir, rp.Recurse)
		}
	} else if rp.File != "" {
		// read specified file for requirements
		b, err := ioutil.ReadFile(rp.File)
		FatalCheck(err)
		reqs = string(b)
	} else {
		// parse the current directory
		reqs = getNpmRequirements(".", rp.Recurse)
	}
	return reqs
}

// TODO: check the exist reqs.yml in current directory if one exists
// and merge the results together, removing duplicate entries
// check the currently installed packages for system and/or pip deps
// and return the string for a reqs.yml
func (rp RequirementsParser) GenerateReqsYml() map[string][]string {
	yml := make(map[string][]string)
	_, packageTool, _ := rp.parseTooling()
	installed := ""
	if packageTool == "apt" {
		installed = aptListInstalled(rp.WithVersion)
	} else if packageTool == "brew" {
		installed = brewListInstalled()
	} else if packageTool == "dnf" {
		installed = dnfListInstalled(rp.WithVersion)
	}

	yml[packageTool] = strings.Split(installed, " ")
	return yml
}

func StdoutReqsYml(yml map[string][]string) {
	for tool, packages := range yml {
		fmt.Println(tool + ":")
		for _, p := range packages {
			for _, line := range strings.Split(p, "\n") {
				fmt.Println("  - " + line)
			}
		}
	}
}
