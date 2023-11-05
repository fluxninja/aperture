#!/usr/bin/env bash

# merge plugin swaggers
set -e

SED="sed"
# check whether we are using macOS
if [ "$(uname)" == "Darwin" ]; then
	SED="gsed"
fi
export SED

git_root=$(git rev-parse --show-toplevel)
metrics_root="$git_root"/docs/content/self-hosting/integrations/metrics
SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

export metrics_root
export SCRIPT_DIR

function generate_otel_docs() {
	key="$1"
	# Remove '- ' from key
	key="${key:2}"
	echo "Generating $key.md"
	rm -rf "$metrics_root"/"$key".md
	camel_case="$(yq e ".$key.camel_case" "$SCRIPT_DIR"/metadata.yaml)"
	metric_name="$(yq e ".$key.metric_name" "$SCRIPT_DIR"/metadata.yaml)"
	receiver_name="$(yq e ".$key.receiver_name" "$SCRIPT_DIR"/metadata.yaml)"
	content=$(cat "$SCRIPT_DIR"/template.md)
	content=$(echo "$content" | $SED "s/CAMEL_CASE/$camel_case/g")
	content=$(echo "$content" | $SED "s/METRIC_NAME/$metric_name/g")
	content=$(echo "$content" | $SED "s/RECEIVER_NAME/$receiver_name/g")
	echo "$content" >"$metrics_root"/"$key".md
	prettier "$metrics_root"/"$key".md --write
}

export -f generate_otel_docs


# Read output of yq command line by line
yq eval 'keys' "$SCRIPT_DIR"/metadata.yaml | while IFS= read -r key; do
    generate_otel_docs "$key"
done
