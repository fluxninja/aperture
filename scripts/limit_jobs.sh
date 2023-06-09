#!/bin/bash

# Define a function to limit the number of concurrent jobs
function limit_jobs() {
  local max_jobs="$1"
  local cmd="${*:2}"
  local job_count
  job_count="$(jobs -p | wc -l)"
  while [ "$job_count" -ge "$max_jobs" ]; do
    sleep 1
    job_count="$(jobs -p | wc -l)"
  done
  $cmd &
}
