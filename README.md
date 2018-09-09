# reqs

Reqs is a cross-platform Linux and MacOSX systems package management tool. It wraps apt, homebrew, dnf and is able to automatically which tool to use based on the system it is on and installs the appropriate dependenices from any number of requirements files or reqs.yml files.  Can automatically update and/or upgrade before installing the reqs with an arg.  Has many useful args.  The main focus of reqs is system package management abstraction with pip and possibly gem support added as an after thought to ease some project deployments.  Because pip and ruby reqs generally don't differ from system-to-system abstracting the tools is not so important to reqs.

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

install pip dependencies from any found 'requirements.txt' files after system dependencies. Must specify path to the pip executable to use.  Can also specify pip requirements inside of reqs.yml under a pip section.
```
reqs -pip $(which pip)
``

## Building

Must have Go installed.  Recent version is better.  Relies on go-dep and go-releaser.  Build script will attempt to install/update both and whatever other deps reqs has using dep.

```
./build.sh
```

## Note
This is a little passion side project for me.  My goal is to have beautiful code here even if this project isn't used widely.  If you have an idea you want implemented with this tool let me know about it. I will review and accept pull requests.  I open sourced this project so feel free to hack it up however you see fit.

## Todo

+ test dnf compatibility for fedora setups
+ add gem, npm, and bower comprehension or just stick to system packages?