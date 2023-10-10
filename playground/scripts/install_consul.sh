#!/bin/bash

set -euo pipefail

git_root=$(git rev-parse --show-toplevel)

pushd "$git_root" || exit 1

if ! command -v consul-k8s &>/dev/null; then
	printf 'Installing Consul CLI\n'
	brew install hashicorp/tap/consul-k8s && consul-k8s version
fi

consul-k8s install --auto-approve --namespace default --config-file "${git_root}/playground/resources/consul/values.yaml" && kubectl apply -f "${git_root}/playground/resources/consul/otel.yaml" || exit 1

popd || exit 1
