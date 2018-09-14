package main

import (
    "flag"
    "github.com/iepathos/reqs"
    log "github.com/sirupsen/logrus"
    "os"
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
    pipPtr := flag.String("pip", "", "install pip dependencies from any 'requirements.txt' found, this arg must be given the path to the pip executable to use")
    pip3Ptr := flag.String("pip3", "", "install pip3 dependencies from any 'requirements.txt' found and any pip3 entries in reqs.yml")
    sudoPipPtr := flag.Bool("spip", false, "install pip dependencies with sudo")
    sudoPip3Ptr := flag.Bool("spip3", false, "install pip3 dependencies with sudo")
    ymlPtr := flag.Bool("yml", false, "stdout the currently installed system requirements in yml format")
    npmPtr := flag.Bool("npm", false, "install global npm dependencies reqs.yml, installs package.json files in the appropriate directories")
    sudoNpmPtr := flag.Bool("snpm", false, "install npm dependencies with sudo")
    flag.Parse()

    if *withVersionPtr {
        *useStdoutPtr = true
    }
    if *sourcesPtr || *useStdoutPtr || *ymlPtr {
        log.SetLevel(log.ErrorLevel)
    } else if !*quietPtr {
        log.SetLevel(log.DebugLevel)
    } else {
        log.SetLevel(log.ErrorLevel)
    }

    if *sudoPipPtr && *pipPtr == "" {
        *pipPtr = "pip"
    }
    if *sudoPip3Ptr && *pip3Ptr == "" {
        *pip3Ptr = "pip3"
    }
    if *sudoNpmPtr {
        *npmPtr = true
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
    if *ymlPtr {
        ymlMap := rp.GenerateReqsYml()
        reqs.StdoutReqsYml(ymlMap)
        os.Exit(0)
    }

    sudo, packageTool, autoYes, requirements := rp.Parse()

    pipRequirements := ""
    if *pipPtr != "" {
        pipRequirements = rp.ParsePip()
    }
    pip3Requirements := ""
    if *pip3Ptr != "" {
        pip3Requirements = rp.ParsePip3()
    }
    npmRequirements := ""
    if *npmPtr {
        npmRequirements = rp.ParseNpm()
    }

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
    pc.Install(*updatePtr, *upgradePtr)

    if *pipPtr != "" {
        reqs.PipInstall(pipRequirements, *pipPtr, *sudoPipPtr, *upgradePtr, *quietPtr)
    }
    if *pip3Ptr != "" {
        reqs.PipInstall(pip3Requirements, *pip3Ptr, *sudoPip3Ptr, *upgradePtr, *quietPtr)
    }
    if *npmPtr {
        globalArg := true
        fromDirectory := ""
        // install global npm requirements
        reqs.NpmInstall(npmRequirements, fromDirectory, *sudoNpmPtr, globalArg, *quietPtr)
        // any directories with package.json in them but where
        // node_modules is not part of the path run just `npm install` inside
        packageDirs := rp.FindNpmPackageDirs()

        for _, pkgDir := range packageDirs {
            reqs.NpmInstall("", pkgDir, false, false, *quietPtr)
        }
    }
}
