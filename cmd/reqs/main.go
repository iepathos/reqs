package main

import (
    "flag"
    "github.com/iepathos/reqs"
    log "github.com/sirupsen/logrus"
)

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

    if *withVersionPtr {
        *useStdoutPtr = true
    }
    if *sourcesPtr || *useStdoutPtr {
        log.SetLevel(log.ErrorLevel)
    } else if !*quietPtr {
        log.SetLevel(log.DebugLevel)
    } else {
        log.SetLevel(log.ErrorLevel)
    }

    rp := reqs.RequirementsParser{
        Dir:         *dirPtr,
        File:        *filePtr,
        UseStdout:   *useStdoutPtr,
        UseStdin:    *useStdinPtr,
        WithVersion: *withVersionPtr,
        Recurse:     *recursePtr,
        Sources:     *sourcesPtr,
    }

    sudo, packageTool, autoYes, requirements := rp.Parse()

    pc := reqs.PackageConfig{
        Tool:    packageTool,
        Sudo:    sudo,
        AutoYes: autoYes,
        Reqs:    requirements,
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