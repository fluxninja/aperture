#!/usr/bin/env bash
set -euo pipefail
set -x

# shellcheck disable=SC2153
readarray -t tools <<<"${TOOLS?}"
if [ "${#tools[@]}" -eq 1 ] && [ -z "${tools[0]:-}" ]; then
	# Parameter was set to empty string
	tools=()
fi

(yes || true) | ./scripts/manage_asdf_tools.sh setup "${tools[@]}"
