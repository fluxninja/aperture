#!/usr/bin/env bash
# Port-forward using kubectl, with simple healthchecks to work around bug in kubectl
set -euo pipefail

check_var_set() {
	local -r var_name="${1?}"
	if [ -z "${!var_name:-}" ]; then
		printf 'Required variable %s is not set!\n' "${var_name}" >&2
		return 1
	fi
}

check_var_set RESOURCE_NAME
check_var_set REMOTE_PORT

: "${NAMESPACE:-}"
: "${RESOURCE_TYPE:=svc}"
: "${LOCAL_PORT:=${REMOTE_PORT}}"

: "${REQUEST_PATH:=/}"
: "${INITIAL_DELAY:=5s}"
: "${TIMEOUT:=1}"
: "${PERIOD:=5s}"
: "${FAILURE_THRESHOLD:=1}"

port_forward_args=()
if [ -n "${NAMESPACE:-}" ]; then
	port_forward_args+=(--namespace "${NAMESPACE}")
fi
port_forward_args+=("${RESOURCE_TYPE}/${RESOURCE_NAME}" "${LOCAL_PORT}:${REMOTE_PORT}")
curl_args=(
	--head                  # Just retrieve headers - ignore body
	--silent                # Don't show stuff on stdout
	--show-error            # But still show errors
	--max-time "${TIMEOUT}" # Make curl exit if it can't finish the operation on time
	"localhost:${LOCAL_PORT}${REQUEST_PATH}"
)

while sleep 1; do
	kubectl port-forward "${port_forward_args[@]}" &
	process_id="${!}"
	trap 'kill "${process_id}"' EXIT
	sleep "${INITIAL_DELAY}" || break
	failures=0
	while :; do
		if curl "${curl_args[@]}" | head -n1; then
			failures=0
		else
			exit_code="${?}"
			if [ "${exit_code}" -eq 124 ]; then
				printf 'Timed out waiting for response\n' >&2
			else
				printf 'Curl failed with exit code %s\n' "${exit_code}" >&2
			fi
			if ((++failures >= FAILURE_THRESHOLD)); then
				printf 'Restarting port-forwarding\n' >&2
				kill "${process_id}" || true # Ignore - might fail if the process exited by itself
				wait "${process_id}" || true # Ignore exit code
				break
			fi
		fi
		sleep "${PERIOD}" || break
	done
done
