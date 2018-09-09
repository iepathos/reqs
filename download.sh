#!/bin/bash
# downloads reqs from github release binary
arch="$(uname)"
if [[ "$arch" == "Darwin" ]]; then
	curl -sL https://github.com/iepathos/reqs/releases/download/v0.2.4/reqs_0.2.4_Darwin_x86_64.tar.gz > reqs.tar.gz
else
	curl -sL https://github.com/iepathos/reqs/releases/download/v0.2.4/reqs_0.2.4_Linux_x86_64.tar.gz > reqs.tar.gz
fi
tar -xzf reqs.tar.gz
rm reqs.tar.gz