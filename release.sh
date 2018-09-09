#!/bin/bash
go get -u github.com/golang/dep
dep ensure

distarg=""
if [[ -d "dist" ]]; then
	distarg="--rm-dist"
fi

go get -u github.com/goreleaser/goreleaser
goreleaser $distarg

# update download version to latest git tag
previoustag="$(git tag | tail -n 2 | head -n 1)"
latesttag="$(git tag | tail -n 1)"
sed "s/$previoustag/$latesttag/" download.sh
git add download.sh
git commit -m "updated download.sh from $previoustag to $latesttag"
git push origin master