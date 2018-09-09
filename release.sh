#!/bin/bash
go get -u github.com/golang/dep
dep ensure

distarg=""
if [[ -d "dist" ]]; then
	distarg="--rm-dist"
fi

go get -u github.com/goreleaser/goreleaser
goreleaser $distarg
