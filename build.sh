#!/bin/bash

dep ensure

if [[ ! -d "bin" ]]; then
	mkdir bin
fi

# build for apple macbook
GOOS=darwin GOARCH=386 go build reqs.go
mv reqs bin/reqs-darwin-386
GOOS=darwin GOARCH=amd64 go build reqs.go
mv reqs bin/reqs-darwin-amd64

# build for linux 
GOOS=linux GOARCH=386 go build reqs.go
mv reqs bin/reqs-linux-386
GOOS=linux GOARCH=amd64 go build reqs.go
mv reqs bin/reqs-linux-amd64