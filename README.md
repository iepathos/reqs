# reqs

Reqs is a cross-platform Linux and MacOSX systems package management tool. Initial aim is to wrap apt and homebrew.  After that, aiming for dnf compatiblity for Fedora deployments.

Best way to use reqs is with a reqs.yml file in you repositories.

reqs.yml
```
common:
  - curl
  - git
apt:
  - golang-go
brew:
  - go
```

Then run `reqs` in your repos and it'll install your system-level dependencies for you.

Can use separate requirements files, like how pip requirements.txt work with package names each on a new line and it tries to install the packages listed in it using either apt-requirements.txt, dnf-requirements.txt, brew-requirements.txt, or common-requirements.txt.

It can gather these requirements for multiple directories and/or recursively and combine them into a single installation call.

reqs automatically determines the tool to used based on the system and what is available.


## Installation

Download the latest release for your system from [https://github.com/iepathos/reqs/releases](https://github.com/iepathos/reqs/releases)

Or install with go if you're gopher inclined
```
go get -u github.com/iepathos/reqs
```

## Usage

Automaticaly finds apt-requirements.txt, brew-requirements.txt, dnf-requirements.txt, common-requirements.txt, and reqs.yml files.  common-requirements.txt are accepted for cross-platform shared same-name system dependencies.

For an example reqs.yml see [https://github.com/iepathos/reqs/blob/master/examples/reqs.yml](https://github.com/iepathos/reqs/blob/master/examples/reqs.yml)

Example dev setup [https://github.com/iepathos/reup](https://github.com/iepathos/reup)

View reqs args and their descriptions
```
reqs -h
```

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


generate apt requirements with the versions info locked installed
```
reqs -ov > apt-requirements.txt
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

## Building

Must have Go installed.  Recent version is better.  Relies on go-dep and go-releaser.  Build script will attempt to install/update both and whatever other deps reqs has using dep.

```
./build.sh
```

## Todo

+ refactor reqs code until it's beautiful
+ test dnf compatibility for fedora setups
+ add pip, gem, npm, and bower comprehension or just stick to system packages?