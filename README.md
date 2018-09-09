# reqs

Reqs is a cross-platform Linux and MacOSX systems package management tool. It wraps apt, homebrew, dnf and is able to automatically determine the right tool to use based on the system.  It checks requirements files and/or reqs.yml files.  

The main focus of reqs is system package management abstraction with pip and possibly gem support added as an after thought to ease some project deployments.  Because pip and ruby reqs generally don't differ from system-to-system abstracting those tools is not so important to reqs.

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
dnf:
  - golang
```

Then run `reqs` in your repos and it'll install your system-level dependencies for you.

Can use separate requirements files, like how pip requirements.txt work with package names each on a new line and it tries to install the packages listed in it using either apt-requirements.txt, dnf-requirements.txt, brew-requirements.txt, or common-requirements.txt.

It can gather these requirements for multiple directories and/or recursively and combine them into a single installation call.

reqs automatically determines the tool to used based on the system and what is available.


## Installation

Download the latest release for your system from [https://github.com/iepathos/reqs/releases](https://github.com/iepathos/reqs/releases)

Or install with go if you're gopher inclined
```
go get -u github.com/iepathos/reqs/cmd/reqs
```

## Usage

Automaticaly finds apt-requirements.txt, brew-requirements.txt, dnf-requirements.txt, common-requirements.txt, and reqs.yml files.  common-requirements.txt are accepted for cross-platform shared same-name system dependencies.

For an example reqs.yml see [https://github.com/iepathos/reqs/blob/master/examples/reqs.yml](https://github.com/iepathos/reqs/blob/master/examples/reqs.yml)

Example dev setup [https://github.com/iepathos/reup](https://github.com/iepathos/reup)

view reqs args and their descriptions
```
reqs -h
```

install all of the example projects with system pip and system npm
```
reqs -r -d examples -spip -snpm
```

install requirements in the current directory
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

## Releasing

Must have Go installed.  Recent version is better.  Relies on go-dep and go-releaser.  `release.sh` will attempt to install/update both  go packages and whatever other deps reqs has using dep.  git tag the current commit you wish to release with the next appropriate version tag and run
```
./release.sh
```

Must export GITHUB_TOKEN with permission to push to origin master for the git repo.  If you just fork off github.com/iepathos/reqs and then use a personal access github token with repo permission you should be groovy.

## Todo

+ refactor reqs code until it's beautiful
+ add gem, npm, and bower comprehension or just stick to system packages?