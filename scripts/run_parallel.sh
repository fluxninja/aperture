#!/usr/bin/env bash

declare -a pids=()
declare -a cmds=()
any_failed=0

run_command() {
	local cmd="$1"
	echo "Running command: $cmd"
	eval "$cmd" &
	local pid=$!
	pids+=("$pid")
	cmds+=("$cmd")
}

for cmd in "$@"; do
	run_command "$cmd"
done

index=0
for pid in "${pids[@]}"; do
	wait "$pid"
	exit_code=$?
	if [ "$exit_code" -ne 0 ]; then
		echo "Command failed: ${cmds[$index]}"
		any_failed=1
	fi
	((index++))
done

if [ "$any_failed" -eq 1 ]; then
	exit 1
fi
exit 0
