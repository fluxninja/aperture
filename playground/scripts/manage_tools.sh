#!/bin/bash

git_root=$(git rev-parse --show-toplevel)

pushd "$git_root" || exit 1

./scripts/manage_tools.sh setup gcloud kubectl kind tilt helm jb tanka jsonnet kustomize

popd || exit 1
