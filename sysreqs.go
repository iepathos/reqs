package main

import (
    "bufio"
    log "github.com/sirupsen/logrus"
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
    reader := bufio.NewReader(os.Stdin)
    text, _ := reader.ReadString('\n')
    text = strings.Replace(text, "\n", " ", -1)

    if runtime.GOOS == "linux" {
        log.Info("Linux system detected, installing system requirements with apt")
        cmd := exec.Command("/bin/sh", "-c", "sudo apt install -y"+text)
        err := cmd.Run()
        if err != nil {
            log.Fatal(err)
        }
    } else if runtime.GOOS == "darwin" {
        log.Info("Darwin system detected")
        if !isCommandAvailable("brew") {
            installHomebrew()
        }
        log.Info("Installing system requirements with brew")
        cmd := exec.Command("/bin/sh", "-c", "brew install "+text)
        err := cmd.Run()
        if err != nil {
            log.Fatal(err)
        }
    }
}
