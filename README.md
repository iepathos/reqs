# reqs

Reqs is a cross-platform package management tool.  It looks to make it stupid easy to manage system-level dependencies across Linux and MacOSX systems.  Initial aim is to wrap apt and homebrew so cross deployment to ubuntu and macosx systems can be easily configurable.  After that, I like Fedora so will get dnf compatible and probably yum.

Takes a requirements file, like a pip requirements.txt with package names each on a new line and it tries to install the packages listed in it using either apt, dnf, or brew.

It can gather these requirements for multiple directories and/or recursively and combine them into a single installation call.

reqs automatically determines the tool to used based on what is available.

reqs looks for files named common-requirements.txt (cross-platform same-name deps), apt-requirements.txt, dnf-requirements.txt, brew-requirements.txt and uses these files to figure out what to install.

reqs can generate the requirements for these files based on what is currently available to it on the system.


Example setup [https://github.com/iepathos/reup](https://github.com/iepathos/reup)



Usage:

Automaticaly find and detect tool-requirements.txt in the directory.  common-requirements.txt accepts for cross-platform shared system dependencies.
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

update package tool before installing requirements.
```
reqs -u
```

TODO:

+ track specific verison installed for brew output, make sure it's compatible with install
+ add dnf compatibility for fedora setups