#!/usr/bin/env bash
set -e

mkdir buildtardir
cp LICENSE ./buildtardir
cp ./dist/aperturectl ./buildtardir
cd  ./buildtardir
rm_v=${APERTURECTL_BUILD_VERSION}
rm_v=${rm_v#v}
touch aperturectl-"${rm_v}"-"${GOOS}"-"${GOARCH}".tar.gz
tar --exclude=aperturectl-"${rm_v}"-"${GOOS}"-"${GOARCH}".tar.gz -czvf aperturectl-"${rm_v}"-"${GOOS}"-"${GOARCH}".tar.gz .
cp -r -- *.tar.gz "$HOME"/project/dist/packages/
#Remove the .deb and .rpm package for darwin OS.
#As this packges don't work on macos and we have .tar.gz file which has aperturectl compiled binary for this systems.
if [ "${GOOS}" == "darwin" ]
then
    rm -rf "$HOME"/project/dist/packages/*.deb "$HOME"/project/dist/packages/*.rpm
fi
ls -al "$HOME"/project/dist/packages/
