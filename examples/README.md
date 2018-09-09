# reqs examples


We have three projects here with different system requirements and pip dependencies.  If we want to deploy to multiple systems, some Ubuntu, some Fedora, some MacOSX then reqs is a good candidate for managing the installation of their system dependencies.

reqs can recurse and accept multiple directories to check for various requirements files.  This makes it especially good at dealing with combinations of services with varied dependencies.  It'll read all the requirements files, use the appropriate tool depending on the system type it's deployed on, and run one install step removing the duplicate dependencies from the call.


All 3 projects on a system can be installed like
```
reqs -r -d /path/to/one,/path/to/two,/path/to/three
```

Or if all the projects are inside a parent directory like

- /parent/dir
	+ service-project1
	+ service-project2
	+ service-project3

Then reqs can handle installing their dependencies like
```
reqs -r -d /parent/dir
```

If we also want to install the pip dependencies after it installs the system dependencies.  We specify the pip executable to use for the pip install step, currently does not allow for specify multiple pip environments for the pip step.
```
reqs -r -d /parent/dir -pip pip
```

Let's also run update and upgrade for the system packages.
```
reqs -r -d /parent/dir -pip pip -up
```

And let's make it all quiet so it doesn't spam everything.  Errors will still get logged.
```
reqs -r -d /parent/dir -pip pip -up -q
```