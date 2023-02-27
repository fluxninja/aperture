#!/usr/bin/env bash

# take go version as argument
goversion=$1
gitroot=$(git rev-parse --show-toplevel)

FIND="find"
SED="sed"
if [[ "$OSTYPE" == "darwin"* ]]; then
	FIND="gfind"
	SED="gsed"
fi

# make sure go version is of the form major.minor.patch
if [[ ! "$goversion" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
	echo "Invalid go version: $goversion"
	exit 1
fi

# parse major version - 1.15.2 ->  major=1 minor=15 patch=2
major=$(echo "$goversion" | cut -d. -f1)
minor=$(echo "$goversion" | cut -d. -f2)
#patch=$(echo "$goversion" | cut -d. -f3)

major_minor="$major.$minor"

# find all go.mod files and update them with the new major version, exclude all vendor directories
# go file contain a line like: go 1.15
"$FIND" . -name go.mod -not -path "*/vendor/*" -exec "$SED" -i "s/go [0-9]\+\.[0-9]\+/go $major_minor/" {} \;

# change version in go_mod_tidy.sh by replacing -compat=1.15 with -compat=1.16
"$SED" -i "s/-compat=[0-9]\+\.[0-9]\+/-compat=$major_minor/" "$gitroot"/scripts/go_mod_tidy.sh

# run go mod tidy
"$gitroot"/scripts/go_mod_tidy.sh

# find Dockerfiles that use go, look for FROM golang:<version>-type and update the version goversion, i.e. major.minor.patch
"$FIND" . -name Dockerfile -not -path "*/vendor/*" -exec "$SED" -i "s/FROM golang:[0-9]\+\.[0-9]\+\.[0-9]\+/FROM golang:$goversion/" {} \;
# find Dockerfiles and replace cache busting version, e.g. id=<name>-<version> with id=<name>-<goversion>
"$FIND" . -name Dockerfile -not -path "*/vendor/*" -exec "$SED" -i -re "s/id=([a-zA-Z0-9]+)-([0-9]+\.[0-9]+\.[0-9]+)/id=\1-$goversion/" {} \;

# update circleci cimg/go version in .circleci/continue-workflow.yml
"$SED" -i "s/cimg\/go:[0-9]\+\.[0-9]\+\.[0-9]\+/cimg\/go:$goversion/" "$gitroot"/.circleci/continue-workflows.yml

# update golang in .tool-versions
"$SED" -i "s/golang [0-9]\+\.[0-9]\+\.[0-9]\+/golang $goversion/" "$gitroot"/.tool-versions

# update .golangci.yml, replace go: "1.15" with go: "1.16"
"$SED" -i "s/go: \"[0-9]\+\.[0-9]\+\"/go: \"$major_minor\"/" "$gitroot"/.golangci.yaml

# in .pre-commit-config.yaml, replace -compat=1.15 with -compat=1.16
"$SED" -i "s/-compat=[0-9]\+\.[0-9]\+/-compat=$major_minor/" "$gitroot"/.pre-commit-config.yaml

# invalidate asdf cache
"$gitroot"/scripts/invalidate_asdf_ci_cache.sh
