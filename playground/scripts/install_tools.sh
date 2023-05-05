#!/bin/bash

git_root=$(git rev-parse --show-toplevel)

pushd "$git_root" || exit 1

./scripts/install_asdf_tools.sh setup gcloud kubectl kind tilt helm jb tanka jsonnet kustomize
./scripts/install_go_tools.sh
./scripts/install_python_tools.sh

popd || exit 1
