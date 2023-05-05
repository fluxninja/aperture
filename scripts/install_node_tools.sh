#!/usr/bin/env bash

set -euo pipefail

gitroot=$(git rev-parse --show-toplevel)

pushd "${gitroot}" >/dev/null

if ! command -v node &>/dev/null; then
	printf 'Installing NodeJS\n'
	./scripts/install_asdf_tools.sh setup nodejs
fi

printf 'Installing NodeJS tools\n'

pushd "${gitroot}" >/dev/null

# install the following tools globally if their commands are not found
# prettier@latest ; look for prettier
# @bitnami/readme-generator-for-helm ; look for readme-generator
# @mermaid-js/mermaid-cli ; look for mmdc

# make an map of tools and commands
tools=(
	"prettier@latest:prettier"
	"@bitnami/readme-generator-for-helm:readme-generator"
	"@mermaid-js/mermaid-cli:mmdc"
)

for tool in "${tools[@]}"; do
	# split the tool and command
	IFS=":" read -ra tool_and_command <<<"${tool}"
	tool_name=${tool_and_command[0]}
	command_name=${tool_and_command[1]}

	# check if the command is found
	if ! command -v "${command_name}" &>/dev/null; then
		# install the tool
		npm install -g "${tool_name}" || true
	fi
done

popd >/dev/null

if asdf where nodejs &>/dev/null; then
	asdf reshim nodejs
fi

popd >/dev/null
