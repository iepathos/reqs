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
)

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

func getSysRequirements(dirPath, packageTool string) string {
    files, err := ioutil.ReadDir(dirPath)
    if err != nil {
        log.Fatal(err)
    }

    // accept packageTool-requirements.txt and common-requirements.txt
    commonRequirements := "common-requirements.txt"
    toolRequirements := packageTool + "-requirements.txt"

    text := ""
    for _, f := range files {
        if f.Name() == commonRequirements || f.Name() == toolRequirements {
            fpath := dirPath + "/" + f.Name()
            log.Info("Found " + fpath)
            b, err := ioutil.ReadFile(fpath)
            if err != nil {
                log.Fatal(err)
            }
            text += "\n" + string(b)
        }
    }
    if len(text) == 0 {
        log.Fatal("No requirements found")
    }
    return text
}

func getAptRequirements() string {
    reqs := ""
    out, err := exec.Command("sudo", "apt", "list", "--installed").Output()
    if err != nil {
        log.Fatal(err)
    }
    // log.Info(string(out))
    for _, line := range strings.Split(string(out), "\n") {
        if strings.Contains(line, "/") {
            lSplit := strings.Split(string(line), "/")
            if reqs == "" {
                reqs += lSplit[0]
            } else {
                reqs += "\n" + lSplit[0]
            }
        }
    }
    return reqs
}

func getBrewRequirements() string {
    out, err := exec.Command("brew", "list").Output()
    if err != nil {
        log.Fatal(err)
    }
    return strings.TrimSpace(string(out))
}

func main() {
    // if arg -d then check the directory for <sys>-requirements.txt files and use them
    // if arg -f then use the specified file for requirements
    // if no args check the current directory

    dirPtr := flag.String("d", "", "directory holding sys-requirements.txt files")
    filePtr := flag.String("f", "", "file to read requirements from")
    generatePtr := flag.String("g", "", "stdout the currently installed requirements for a specified tool apt, dnf, or brew")
    useStdin := flag.Bool("i", false, "use stdin for requirements")
    flag.Parse()

    var packageTool string
    sudo := ""
    autoYes := ""
    if runtime.GOOS == "linux" {
        if *generatePtr == "" {
            log.Info("Linux system detected")
        }
        if isCommandAvailable("apt") {
            packageTool = "apt"
        } else if isCommandAvailable("dnf") {
            packageTool = "dnf"
        }
        sudo = "sudo "
        autoYes = "-y "
    } else if runtime.GOOS == "darwin" {
        if *generatePtr == "" {
            log.Info("Darwin system detected")
        }
        if !isCommandAvailable("brew") {
            installHomebrew()
        }
        packageTool = "brew"
    } else if runtime.GOOS == "windows" {
        log.Fatal("Windows system detected, abandon all hope")
    }

    reqs := ""
    if *dirPtr != "" {
        reqs = getSysRequirements(*dirPtr, packageTool)
    } else if *filePtr != "" {
        b, err := ioutil.ReadFile(*filePtr)
        if err != nil {
            log.Fatal(err)
        }
        reqs = string(b)
    } else if *useStdin {
        reader := bufio.NewReader(os.Stdin)
        reqs, _ = reader.ReadString('\n')
    } else if *generatePtr != "" {
        // stdout requirements
        if packageTool == "apt" {
            reqs = getAptRequirements()
        } else if packageTool == "brew" {
            reqs = getBrewRequirements()
        }
        fmt.Print(reqs)
        os.Exit(0)
    } else {
        reqs = getSysRequirements(".", packageTool)
    }
    reqs = strings.Replace(reqs, "\n", " ", -1)
    log.Info(reqs)

    log.Info("Installing system requirements with " + packageTool)
    log.Info(sudo + packageTool + " install " + autoYes + reqs)
    cmd := exec.Command("/bin/sh", "-c", sudo+packageTool+" install "+autoYes+reqs)
    err := cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
}
