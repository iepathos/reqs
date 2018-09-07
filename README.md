# sysreqs

Abstract cross-platform package install tool.  Takes a requirements file, like a pip requirements.txt with package names each on a new line and it tries to install the packages listed in it using either apt, brew, or dnf depending on the operating system.


Usage:

```
sysreqs < sys-requirements.txt

```