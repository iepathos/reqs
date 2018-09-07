package main

import (
    "bufio"
    "bytes"
    "flag"
    "fmt"
    log "github.com/sirupsen/logrus"
    "io/ioutil"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "strings"
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
    log.Info("Installing homebrew")
    cmd := exec.Command("/usr/bin/ruby",
        "-e",
        "\"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)\"")
    err := cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
}

func runShell(code string) {
    cmd := exec.Command("/bin/sh", "-c", code)
    err := cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
}

func updatePackages(packageTool string) {
    log.Info("Updating " + packageTool + " packages")
    if packageTool == "apt" {
        runShell("sudo apt update -y")
    } else if packageTool == "brew" {
        runShell("brew update")
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
            fileNames = append(fileNames, dirPath+"/"+f.Name())
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

func getSysRequirementsMultipleDirs(dirPaths []string, packageTool string, recurse bool) (reqs string) {
    for _, dirPath := range dirPaths {
        reqs = newLineIfNotEmpty(reqs, getSysRequirements(dirPath, packageTool, recurse))
    }
    return reqs
}

func getInstalledAptRequirements(withVersion bool) (reqs string) {
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

func parseRequirements(dirPath, filePath, packageTool string,
    outputArg, useStdin, withVersion, recurse bool) (reqs string) {
    if dirPath != "" {
        if strings.Contains(dirPath, ",") {
            reqs = getSysRequirementsMultipleDirs(strings.Split(dirPath, ","), packageTool, recurse)
        } else {
            reqs = getSysRequirements(dirPath, packageTool, recurse)
        }
    } else if filePath != "" {
        b, err := ioutil.ReadFile(filePath)
        if err != nil {
            log.Fatal(err)
        }
        reqs = string(b)
    } else if useStdin {
        reader := bufio.NewReader(os.Stdin)
        reqs, _ = reader.ReadString('\n')
    } else if outputArg {
        if packageTool == "apt" {
            reqs = getInstalledAptRequirements(withVersion)
        } else if packageTool == "brew" {
            reqs = getInstalledBrewRequirements()
        }
        fmt.Print(reqs)
        os.Exit(0)
    } else {
        // parse the current directory
        reqs = getSysRequirements(".", packageTool, recurse)
    }
    reqs = strings.TrimSpace(strings.Replace(reqs, "\n", " ", -1))
    return reqs
}

func installRequirements(reqs, packageTool, autoYes, sudo string, quiet, force bool) {
    log.Info("Installing system requirements with " + packageTool)
    log.Info(sudo + packageTool + " install " + autoYes + reqs)
    forceArg := ""
    if force {
        if packageTool == "brew" {
            forceArg = "--force "
        } else {
            forceArg = "-f "
        }
    }
    cmd := exec.Command("/bin/sh", "-c", sudo+packageTool+" install "+forceArg+autoYes+reqs)
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    err := cmd.Run()
    if !quiet {
        fmt.Print(string(out.String()))
    }
    if err != nil {
        log.Fatal(stderr.String())
    }
}

func determinePackageTooling(useStdout bool) (sudo, autoYes, packageTool string) {
    if runtime.GOOS == "linux" {
        if !useStdout {
            log.Info("Linux system detected")
        }
        linuxTools := []string{
            "apt",
            "dnf",
            "yum",
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
        if !useStdout {
            log.Info("Darwin system detected")
        }
        if !isCommandAvailable("brew") {
            installHomebrew()
        }
        packageTool = "brew"
    } else if runtime.GOOS == "windows" {
        log.Fatal("Windows system detected, abandon all hope")
    }

    return sudo, autoYes, packageTool
}

func main() {
    // if arg -d then check the directory for <sys>-requirements.txt files and use them
    // if arg -f then use the specified file for requirements
    // if no args check the current directory

    dirPtr := flag.String("d", "", "directory or comma separated directories with requirements files")
    filePtr := flag.String("f", "", "specific requirements file to read from")
    useStdoutPtr := flag.Bool("o", false, "stdout the currently installed requirements for apt or brew")
    useStdinPtr := flag.Bool("i", false, "use stdin for requirements")
    withVersionPtr := flag.Bool("v", false, "save version with output requirements command")
    quiet := flag.Bool("q", false, "silence logging to error level")
    recurse := flag.Bool("r", false, "recurse down directories to find requirements")
    update := flag.Bool("u", false, "update package tool before install")
    force := flag.Bool("force", false, "force reinstall packages")
    flag.Parse()

    if !*quiet {
        log.SetLevel(log.DebugLevel)
    } else {
        log.SetLevel(log.ErrorLevel)
    }

    sudo, autoYes, packageTool := determinePackageTooling(*useStdoutPtr)

    reqs := parseRequirements(*dirPtr, *filePtr, packageTool,
        *useStdoutPtr, *useStdinPtr, *withVersionPtr, *recurse)

    if *update {
        updatePackages(packageTool)
    }

    installRequirements(reqs, packageTool, autoYes, sudo, *quiet, *force)
}
