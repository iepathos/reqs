# reqs examples


We have three projects here with different system requirements and pip dependencies.  If we want to deploy one to multiple systems, some Ubuntu, some Fedora, some MacOSX then we can use reqs to manage all their dependency installation.


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