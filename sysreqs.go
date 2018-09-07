package main

import (
    "bufio"
    "flag"
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

func main() {
    // if arg -d then check the directory for <sys>-requirements.txt files and use them
    // if arg -f then use the specified file for requirements
    // if no args use stdin

    dirPtr := flag.String("d", "", "directory holding sys-requirements.txt files")
    filePtr := flag.String("f", "", "file to read requirements from")
    flag.Parse()

    text := ""
    if *dirPtr != "" {
        log.Info(*dirPtr)
    } else if *filePtr != "" {
        b, err := ioutil.ReadFile(*filePtr)
        if err != nil {
            log.Fatal(err)
        }
        text = string(b)
    } else {
        reader := bufio.NewReader(os.Stdin)
        text, _ = reader.ReadString('\n')
    }

    text = strings.Replace(text, "\n", " ", -1)

    var packageTool string
    sudo := ""
    autoYes := ""
    if runtime.GOOS == "linux" {
        log.Info("Linux system detected")
        if isCommandAvailable("apt") {
            packageTool = "apt"
        } else if isCommandAvailable("dnf") {
            packageTool = "dnf"
        } else if isCommandAvailable("yum") {
            packageTool = "yum"
        }
        autoYes = "-y "
        sudo = "sudo "
    } else if runtime.GOOS == "darwin" {
        log.Info("Darwin system detected")
        if !isCommandAvailable("brew") {
            installHomebrew()
        }
        packageTool = "brew"
    } else if runtime.GOOS == "windows" {
        log.Fatal("Windows system detected, abandon all hope")
    }
    log.Info("Installing system requirements with " + packageTool)
    cmd := exec.Command("/bin/sh", "-c", sudo+packageTool+" install "+autoYes+text)
    err := cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
}
