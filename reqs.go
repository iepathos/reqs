package main

import (
    "bufio"
    "bytes"
    "flag"
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

func fatalCheck(err error) {
    if err != nil {
        log.Fatal(err)
    }
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

func isCommandAvailable(name string) bool {
    cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
    if err := cmd.Run(); err != nil {
        return false
    }
    return true
}

func installHomebrew() {
    log.Info("Installing homebrew")
    cmd := exec.Command("/usr/bin/ruby",
        "-e",
        "\"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)\"")
    err := cmd.Run()
    fatalCheck(err)
}

func runShell(code string) {
    log.Info(code)
    cmd := exec.Command("/bin/sh", "-c", code)
    err := cmd.Run()
    fatalCheck(err)
}

func getAptSources() (sources string) {
    b, err := ioutil.ReadFile("/etc/apt/sources.list")
    fatalCheck(err)
    sources += "\n" + string(b)
    return sources
}

func getBrewTaps() string {
    out, err := exec.Command("brew", "tap").Output()
    fatalCheck(err)
    return strings.TrimSpace(string(out))
}

func recurseForRequirementsFiles(searchPath string) []string {
    filepathList := []string{}
    err := filepath.Walk(searchPath, func(path string, f os.FileInfo, err error) error {
        filepathList = append(filepathList, path)
        return nil
    })
    fatalCheck(err)

    requirementsFilePaths := []string{}
    for _, path := range filepathList {
        if strings.Contains(path, "requirements.txt") || strings.Contains(path, "reqs.yml") {
            requirementsFilePaths = append(requirementsFilePaths, path)
        }
    }
    return requirementsFilePaths
}

func getRequirementFilenames(dirPath string, recurse bool) (fileNames []string) {
    if recurse {
        fileNames = recurseForRequirementsFiles(dirPath)
    } else {
        files, err := ioutil.ReadDir(dirPath)
        fatalCheck(err)
        for _, f := range files {
            fileNames = append(fileNames, dirPath+"/"+f.Name())
        }
    }
    return fileNames
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
            fatalCheck(err)
            text = appendNewLinesOnly(text, string(b))
        } else if strings.Contains(fname, reqsYml) {
            log.Info("Found " + fname)
            b, err := ioutil.ReadFile(fname)
            fatalCheck(err)
            m := make(map[string][]string)
            err = yaml.Unmarshal(b, &m)
            fatalCheck(err)

            for tool, packages := range m {
                if tool == "common" || tool == packageTool {
                    // add the list to text
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
    fatalCheck(err)
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
    fatalCheck(err)
    return strings.TrimSpace(string(out))
}

func getInstalledDnfRequirements(withVersion bool) (reqs string) {
    out, err := exec.Command("sudo", "dnf", "list", "installed").Output()
    fatalCheck(err)
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
func (rp RequirementsParser) parseTooling() (sudo, autoYes, packageTool string) {
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

    return sudo, autoYes, packageTool
}

// determine package tool and args on this system
func (rp RequirementsParser) Parse() (sudo, packageTool, autoYes, reqs string) {
    sudo, autoYes, packageTool = rp.parseTooling()
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
        fatalCheck(err)
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

type PackageConfig struct {
    Tool          string
    Sudo, AutoYes string
    Reqs          string
    Quiet, Force  bool
}

func (pc PackageConfig) getForceArg() (forceArg string) {
    if pc.Force {
        if pc.Tool == "brew" {
            forceArg = "--force "
        } else {
            forceArg = "-f "
        }
    }
    return forceArg
}

func (pc PackageConfig) Install() {
    log.Info("Installing system requirements with " + pc.Tool)
    forceArg := pc.getForceArg()
    cmdStr := pc.Sudo + pc.Tool + " install " + pc.AutoYes + forceArg + pc.Reqs
    log.Info(cmdStr)
    cmd := exec.Command("/bin/sh", "-c", cmdStr)
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    err := cmd.Run()
    if !pc.Quiet {
        fmt.Print(string(out.String()))
    }
    if err != nil {
        log.Fatal(stderr.String())
    }
}

func (pc PackageConfig) abstractUp(upArg string) {
    log.Info("Updating " + pc.Tool + " packages")
    forceArg := pc.getForceArg()
    if pc.Tool == "brew" {
        runShell(pc.Tool + " " + upArg + " " + forceArg)
    } else {
        runShell("sudo " + pc.Tool + " " + upArg + " " + forceArg + pc.AutoYes)
    }
}

func (pc PackageConfig) Update() {
    pc.abstractUp("update")
}

func (pc PackageConfig) Upgrade() {
    pc.abstractUp("upgrade")
}

func main() {
    // if arg -d then check the directory for <sys>-requirements.txt files and use them
    // if arg -f then use the specified file for requirements
    // if no args check the current directory

    dirPtr := flag.String("d", "", "directory or comma separated directories with requirements files")
    filePtr := flag.String("f", "", "specific requirements file to read from")
    useStdoutPtr := flag.Bool("o", false, "stdout the currently installed requirements for apt or brew")
    useStdinPtr := flag.Bool("i", false, "use stdin for requirements")
    withVersionPtr := flag.Bool("ov", false, "stdout the currently installed apt packages with version info")
    quietPtr := flag.Bool("q", false, "silence logging to error level")
    recursePtr := flag.Bool("r", false, "recurse down directories to find requirements")
    updatePtr := flag.Bool("u", false, "update packages before install")
    forcePtr := flag.Bool("force", false, "force reinstall packages")
    upgradePtr := flag.Bool("up", false, "update and upgrade packages before install")
    sourcesPtr := flag.Bool("so", false, "stdout package tool sources")
    flag.Parse()

    if *sourcesPtr || *useStdoutPtr {
        log.SetLevel(log.ErrorLevel)
    } else if !*quietPtr {
        log.SetLevel(log.DebugLevel)
    } else {
        log.SetLevel(log.ErrorLevel)
    }

    rp := RequirementsParser{
        Dir:         *dirPtr,
        File:        *filePtr,
        UseStdout:   *useStdoutPtr,
        UseStdin:    *useStdinPtr,
        WithVersion: *withVersionPtr,
        Recurse:     *recursePtr,
        Sources:     *sourcesPtr,
    }

    sudo, packageTool, autoYes, reqs := rp.Parse()

    pc := PackageConfig{
        Tool:    packageTool,
        Sudo:    sudo,
        AutoYes: autoYes,
        Reqs:    reqs,
        Force:   *forcePtr,
        Quiet:   *quietPtr,
    }
    if *updatePtr || *upgradePtr {
        pc.Update()
    }
    if *upgradePtr {
        pc.Upgrade()
    }
    pc.Install()
}
