#!/bin/bash
# downloads reqs from github release binary
# can pass version as argument to download script
# to get a particular, should get the latest by default
if [[ -z "$1" ]]; then
	version="v0.4.0"
else
	version="$1"
fi
arch="$(uname)"

if [[ "$arch" == "Darwin" ]]; then
	url="https://github.com/iepathos/reqs/releases/download/${version}/reqs_${version//v}_Darwin_x86_64.tar.gz"
else
	url="https://github.com/iepathos/reqs/releases/download/${version}/reqs_${version//v}_Linux_x86_64.tar.gz"
fi
echo $url
curl -sL $url > reqs.tar.gz
tar -xzf reqs.tar.gz reqs
rm reqs.tar.gz