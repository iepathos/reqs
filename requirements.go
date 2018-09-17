package reqs

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// requirements for parsing requirements files
// and for determining currently installed requirements

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

func GetRequirementFilenames(dirPath string, recurse bool) (fileNames []string) {
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
	fileNames := GetRequirementFilenames(dirPath, recurse)
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
		log.Warn("No system requirements files found")
	}
	return strings.TrimSpace(text)
}

func getSysRequirementsMultipleDirs(dirPaths []string, packageTool string, recurse bool) (reqs string) {
	for _, dirPath := range dirPaths {
		reqs = NewLineIfNotEmpty(reqs, getSysRequirements(dirPath, packageTool, recurse))
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

func (rp RequirementsParser) FindNpmPackageDirs() (packageDirs []string) {
	dirArg := "."
	if rp.Dir != "" {
		dirArg = rp.Dir
	}
	packageDirs = FindNpmPackageDirs(dirArg, rp.Recurse)
	return packageDirs
}

func (rp RequirementsParser) ListInstalled(packageTool string) (requirements string) {
	switch packageTool {
	case "apt":
		requirements = AptListInstalled(rp.WithVersion)
	case "brew":
		requirements = BrewListInstalled()
	case "dnf":
		requirements = DnfListInstalled(rp.WithVersion)
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
			"yum",
		}
		for _, tool := range linuxTools {
			if IsCommandAvailable(tool) {
				packageTool = tool
				break
			}
		}
		if packageTool == "" {
			log.Fatal("Failed to find supported package management tools")
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
			InstallHomebrew()
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
			fmt.Print(GetAptSources())
		case "brew":
			fmt.Print(GetBrewTaps())
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
			reqs = GetPipRequirementsMultipleDirs(strings.Split(rp.Dir, ","), rp.Recurse)
		} else {
			reqs = GetPipRequirements(rp.Dir, rp.Recurse)
		}
	} else if rp.File != "" {
		// read specified file for requirements
		b, err := ioutil.ReadFile(rp.File)
		FatalCheck(err)
		reqs = string(b)
	} else {
		// parse the current directory
		reqs = GetPipRequirements(".", rp.Recurse)
	}
	return reqs
}

func (rp RequirementsParser) ParsePip3() (reqs string) {
	if rp.Dir != "" {
		// search directory for requirements
		if strings.Contains(rp.Dir, ",") {
			reqs = GetPip3RequirementsMultipleDirs(strings.Split(rp.Dir, ","), rp.Recurse)
		} else {
			reqs = GetPip3Requirements(rp.Dir, rp.Recurse)
		}
	} else if rp.File != "" {
		// read specified file for requirements
		b, err := ioutil.ReadFile(rp.File)
		FatalCheck(err)
		reqs = string(b)
	} else {
		// parse the current directory
		reqs = GetPip3Requirements(".", rp.Recurse)
	}
	return reqs
}

func (rp RequirementsParser) ParseNpm() (reqs string) {
	if rp.Dir != "" {
		// search directory for requirements
		if strings.Contains(rp.Dir, ",") {
			reqs = GetNpmRequirementsMultipleDirs(strings.Split(rp.Dir, ","), rp.Recurse)
		} else {
			reqs = GetNpmRequirements(rp.Dir, rp.Recurse)
		}
	} else if rp.File != "" {
		// read specified file for requirements
		b, err := ioutil.ReadFile(rp.File)
		FatalCheck(err)
		reqs = string(b)
	} else {
		// parse the current directory
		reqs = GetNpmRequirements(".", rp.Recurse)
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
		installed = AptListInstalled(rp.WithVersion)
	} else if packageTool == "brew" {
		installed = BrewListInstalled()
	} else if packageTool == "dnf" {
		installed = DnfListInstalled(rp.WithVersion)
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
