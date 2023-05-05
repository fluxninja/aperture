#!/bin/bash

git_root=$(git rev-parse --show-toplevel)

pushd "$git_root" || exit 1

./scripts/manage_asdf_tools.sh setup gcloud kubectl kind tilt helm jb tanka jsonnet kustomize
./scripts/manage_go_tools.sh
./scripts/manage_python_tools.sh

popd || exit 1
