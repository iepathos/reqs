#!/bin/bash
# downloads reqs from github release binary
version="v0.3.5"
arch="$(uname)"

if [[ "$arch" == "Darwin" ]]; then
	url="https://github.com/iepathos/reqs/releases/download/${version}/reqs_${version//v}_Darwin_x86_64.tar.gz"
else
	url="https://github.com/iepathos/reqs/releases/download/${version}/reqs_${version//v}_Linux_x86_64.tar.gz"
fi
echo $url
curl -sL $url > reqs.tar.gz
tar -xzf reqs.tar.gz
rm reqs.tar.gz