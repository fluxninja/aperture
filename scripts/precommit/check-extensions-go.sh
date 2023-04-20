#!/usr/bin/env bash

set -e

git_root=$(git rev-parse --show-toplevel)

files=(
	"$git_root"/cmd/aperture-agent/extensions.go
	"$git_root"/cmd/aperture-controller/extensions.go
)

# check whether the argument is generate
if [[ "$1" == "generate" ]]; then
	# generate md5sum for the files
	for file in "${files[@]}"; do
		if [[ -f "${file}" ]]; then
			md5sum=$(md5sum "${file}" | awk '{print $1}')
			md5sum_file="${file}.md5sum"
			echo "${md5sum}" >"${md5sum_file}"
		fi
	done
fi

for file in "${files[@]}"; do
	if [[ -f "${file}" ]]; then
		mdg5sum=$(md5sum "${file}" | awk '{print $1}')
		md5sum_file="${file}.md5sum"
		if [[ -f "${md5sum_file}" ]]; then
			md5sum=$(cat "${md5sum_file}")
			if [[ "${mdg5sum}" != "${md5sum}" ]]; then
				echo "ERROR: ${file} has been modified. Please run 'make extensions_md5sum' to update ${file}.md5sum"
				exit 1
			fi
		else
			echo "ERROR: ${file}.md5sum does not exist. Please run 'make extensions_md5sum' to generate ${file}.md5sum"
			exit 1
		fi
	fi
done
