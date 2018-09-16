package reqs

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
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
