#!/bin/bash

dep ensure

distarg=""
if [[ -d "dist" ]]; then
	distarg="--rm-dist"
fi

goreleaser $distarg
