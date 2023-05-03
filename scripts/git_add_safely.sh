#!/usr/bin/env bash

# This script takes files to git add from the command line and does so safely with retries and sleeps if .git/index.lock is there

if [ -z "$1" ]; then
	echo "Usage: $0 <file1> [<file2> ...]"
	exit 1
fi

function git_add_safely {
	local retries=5
	local sleep_time=2

	while ((retries > 0)); do
		if git add "$@" 2>/dev/null; then
			echo "Files added successfully:"
			echo "$@"
			break
		else
			echo "Failed to add files. Retrying in $sleep_time seconds..."
			sleep "$sleep_time"
		fi
		((retries--))
	done

	if ((retries == 0)); then
		echo "Failed to add files after multiple attempts. Please try again later."
		exit 2
	fi
}

git_add_safely "$@"
