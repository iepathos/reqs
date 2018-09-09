#!/bin/bash
# downloads reqs from github release binary
version="v0.2.9"
arch="$(uname)"

if [[ "$arch" == "Darwin" ]]; then
	url="https://github.com/iepathos/reqs/releases/download/v${version//v}/reqs}${version//v}_Darwin_x86_64.tar.gz"
else
	url="https://github.com/iepathos/reqs/releases/download/v${version//v}/reqs_${version//v}_Linux_x86_64.tar.gz"
fi
echo $url
curl -sL $url > reqs.tar.gz
tar -xzf reqs.tar.gz
rm reqs.tar.gz