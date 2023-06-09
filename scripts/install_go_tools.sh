#!/usr/bin/env bash

set -euo pipefail

gitroot=$(git rev-parse --show-toplevel)

pushd "${gitroot}" >/dev/null

aperturectl="$("$gitroot"/scripts/build_aperturectl.sh)"
# download jsonnet dependencies as well
"$aperturectl" blueprints list --uri="$gitroot"/blueprints >/dev/null

if ! command -v go &>/dev/null; then
	printf 'Installing Go\n'
	./scripts/install_asdf_tools.sh setup golang
fi

printf 'Installing Go tools\n'

pushd "$gitroot"/tools/go >/dev/null

# first run go mod download
go mod download

tools=$(grep _ tools.go | awk -F'"' '{print $2}')

echo "$tools" | while IFS= read -r tool; do
    go install "$tool" &
done

wait  # Wait for all background jobs to complete


popd >/dev/null

if asdf where golang &>/dev/null; then
	asdf reshim golang
fi

popd >/dev/null
