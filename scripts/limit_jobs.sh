#!/bin/bash

# Define a function to limit the number of concurrent jobs
function limit_jobs() {
	local max_jobs="$1"
	local cmd="${*:2}"
	while [ "$(jobs -p | wc -l)" -ge "$max_jobs" ]; do
		sleep 1
	done
	$cmd &
}
