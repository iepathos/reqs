# reqs

Reqs is a cross-platform package management tool.  It looks to make it stupid easy to manage system-level dependencies across Linux and MacOSX systems.  Initial aim is to wrap apt and homebrew so cross deployment to ubuntu and macosx systems can be easily configurable.  After that, I like Fedora so will get dnf compatible and probably yum.

Takes a requirements file, like a pip requirements.txt with package names each on a new line and it tries to install the packages listed in it using either apt, dnf, or brew.

It can gather these requirements for multiple directories and/or recursively and combine them into a single installation call.

reqs automatically determines the tool to used based on what is available.

reqs looks for files named common-requirements.txt (cross-platform same-name deps), apt-requirements.txt, dnf-requirements.txt, brew-requirements.txt and uses these files to figure out what to install.

reqs can generate the requirements for these files based on what is currently available to it on the system.


Example setup [https://github.com/iepathos/reup](https://github.com/iepathos/reup)

## Installation

Download the latest release for your system from [https://github.com/iepathos/reqs/releases](https://github.com/iepathos/reqs/releases)

Or install with go if you're gopher inclined
```
go get -u github.com/iepathos/reqs
```

## Usage

Automaticaly finds apt-requirements.txt, brew-requirements.txt, dnf-requirements.txt, common-requirements.txt, and reqs.yml files.  common-requirements.txt are accepted for cross-platform shared same-name system dependencies.

For an example reqs.yml see [https://github.com/iepathos/reqs/blob/master/examples/reqs.yml](https://github.com/iepathos/reqs/blob/master/examples/reqs.yml)


```
reqs
```

recurse down directories to find requirements files
```
reqs -r
```

get requirements from a specific directory, automaticaly detect appropriate <system-tool>-requirements.txt to use
```
reqs -d /some/path/
```

get requirements from a specific file
```
reqs -f tool-requirements.txt
```

get requirements from stdin
```
reqs -i < tool-requirements.txt
```


generate apt requirements from the currently installed apt packages
```
reqs -o > apt-requirements.txt
```


generate apt requirements with the specific versions installed
```
reqs -o -v > apt-requirements.txt
```

generate brew requirements from the currently install brew packages
```
reqs -o > brew-requirements.txt
```

update packages before installing requirements
```
reqs -u
```

update and upgrade packages before installing requirements
```
reqs -up
```

quiet mode squelch everything but errors
```
reqs -q
```

force reinstall of packages
```
reqs -force
```

## build

Must have Go installed.  Recent version is better.  Relies on go-dep and go-releaser.  Build script will attempt to install/update both and whatever other deps reqs has using dep.

```
./build.sh
```

TODO:
+ refactor reqs code until it's beautiful
+ add dnf compatibility for fedora setups
+ add pip, gem, npm, and bower comprehension or just stick to system packages?