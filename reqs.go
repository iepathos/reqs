package main

import (
    "bufio"
    "flag"
    "fmt"
    log "github.com/sirupsen/logrus"
    "io/ioutil"
    "os"
    "os/exec"
    "runtime"
    "strings"
    "path/filepath"
)

func newLineIfNotEmpty(text, newText string) string {
    if text == "" {
        text = newText
    } else {
        text += "\n" + newText
    }
    return text
}

func isCommandAvailable(name string) bool {
    cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
    if err := cmd.Run(); err != nil {
        return false
    }
    return true
}

func installHomebrew() {
    log.Info("Installing Homebrew")
    cmd := exec.Command("/usr/bin/ruby",
        "-e",
        "\"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)\"")
    err := cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
}

func recurseForRequirementsFiles(searchPath string) []string {
    filepathList := []string{}
    err := filepath.Walk(searchPath, func(path string, f os.FileInfo, err error) error {
        filepathList = append(filepathList, path)
        return nil
    })
    if err != nil {
        log.Fatal(err)
    }

    requirementsFilePaths := []string{}
    for _, path := range filepathList {
        if strings.Contains(path, "-requirements.txt") {
            requirementsFilePaths = append(requirementsFilePaths, path)
        }
    }
    return requirementsFilePaths
}

func getSysRequirements(dirPath, packageTool string, recurse bool) string {
    fileNames := []string{}
    if recurse {
        fileNames = recurseForRequirementsFiles(dirPath)
    } else {
        files, err := ioutil.ReadDir(dirPath)
        if err != nil {
            log.Fatal(err)
        }
        for _, f := range files {
            fileNames = append(fileNames, dirPath + "/" + f.Name())
        }
    }

    // accept packageTool-requirements.txt and common-requirements.txt
    commonRequirements := "common-requirements.txt"
    toolRequirements := packageTool + "-requirements.txt"

    text := ""
    for _, fname := range fileNames {
        if strings.Contains(fname, commonRequirements) || strings.Contains(fname, toolRequirements) {
            log.Info("Found " + fname)
            b, err := ioutil.ReadFile(fname)
            if err != nil {
                log.Fatal(err)
            }
            text += "\n" + string(b)
        }
    }
    if len(text) == 0 {
        log.Fatal("No requirements found")
    }
    return strings.TrimSpace(text)
}

func getSysRequirementsMultipleDirs(dirPaths []string, packageTool string, recurse bool) string {
    allReqs := ""
    for _, dirPath := range dirPaths {
        allReqs = newLineIfNotEmpty(allReqs, getSysRequirements(dirPath, packageTool, recurse))
    }
    return allReqs
}

func getInstalledAptRequirements(withVersion bool) string {
    reqs := ""
    out, err := exec.Command("sudo", "apt", "list", "--installed").Output()
    if err != nil {
        log.Fatal(err)
    }
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
    if err != nil {
        log.Fatal(err)
    }
    return strings.TrimSpace(string(out))
}

func parseRequirements(dirPtr, filePtr, packageTool string,
                       outputPtr, useStdinPtr,
                       withVersionPtr, quiet, recurse bool) string {
    reqs := ""
    if dirPtr != "" {
        // check if , in *dirPtr and gather from multiple directories if so
        if strings.Contains(dirPtr, ",") {
            // gather from multiple directories
            reqs = getSysRequirementsMultipleDirs(strings.Split(dirPtr, ","), packageTool, recurse)
        } else {
            reqs = getSysRequirements(dirPtr, packageTool, recurse)
        }
    } else if filePtr != "" {
        b, err := ioutil.ReadFile(filePtr)
        if err != nil {
            log.Fatal(err)
        }
        reqs = string(b)
    } else if useStdinPtr {
        reader := bufio.NewReader(os.Stdin)
        reqs, _ = reader.ReadString('\n')
    } else if outputPtr {
        if packageTool == "apt" {
            reqs = getInstalledAptRequirements(withVersionPtr)
        } else if packageTool == "brew" {
            reqs = getInstalledBrewRequirements()
        }
        fmt.Print(reqs)
        os.Exit(0)
    } else {
        reqs = getSysRequirements(".", packageTool, recurse)
    }
    reqs = strings.TrimSpace(strings.Replace(reqs, "\n", " ", -1))
    return reqs
}

func main() {
    // if arg -d then check the directory for <sys>-requirements.txt files and use them
    // if arg -f then use the specified file for requirements
    // if no args check the current directory

    dirPtr := flag.String("d", "", "directory or comma separated directories with requirements files")
    filePtr := flag.String("f", "", "requirements file to use")
    outputPtr := flag.Bool("o", false, "stdout the currently installed requirements for a specified tool apt, dnf, or brew")
    useStdinPtr := flag.Bool("i", false, "use stdin for requirements")
    withVersionPtr := flag.Bool("v", false, "save version with output requirements command")
    quiet := flag.Bool("q", false, "silence logging to error level")
    recurse := flag.Bool("r", false, "recurse down directories to find requirements")
    flag.Parse()

    if !*quiet {
        log.SetLevel(log.DebugLevel)
    } else {
        log.SetLevel(log.ErrorLevel)
    }

    var packageTool string
    sudo := ""
    autoYes := ""
    linuxTools := []string{
        "apt",
        "dnf",
    }

    // identify operating system and available package management tool
    if runtime.GOOS == "linux" {
        if !*outputPtr {
            log.Info("Linux system detected")
        }
        for _, tool := range linuxTools {
            if isCommandAvailable(tool) {
                packageTool = tool
                break
            }
        }
        sudo = "sudo "
        autoYes = "-y "
    } else if runtime.GOOS == "darwin" {
        if !*outputPtr {
            log.Info("Darwin system detected")
        }
        if !isCommandAvailable("brew") {
            installHomebrew()
        }
        packageTool = "brew"
    } else if runtime.GOOS == "windows" {
        log.Fatal("Windows system detected, abandon all hope")
    }

    reqs := parseRequirements(*dirPtr, *filePtr, packageTool, *outputPtr,
                              *useStdinPtr, *withVersionPtr,
                              *quiet, *recurse)
    log.Info(reqs)

    log.Info("Installing system requirements with " + packageTool)
    log.Info(sudo + packageTool + " install " + autoYes + reqs)
    out, err := exec.Command("/bin/sh", "-c", sudo+packageTool+" install "+autoYes+reqs).Output()
    if !*quiet {
        fmt.Print(string(out))
    }
    if err != nil {
        log.Fatal(err)
    }
}
