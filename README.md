# reqs

Abstract cross-platform package install tool.  Takes a requirements file, like a pip requirements.txt with package names each on a new line and it tries to install the packages listed in it using either apt, brew, or dnf depending on the operating system.

Usage:

Automaticaly find and detect tool-requirements.txt in the directory.  common-requirements.txt accepts for cross-platform shared system dependencies.
```
reqs
```

get requiremetns from a specific directory, automaticaly detect appropriate <system-tool>-requirements.txt to use
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

