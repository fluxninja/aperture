#!/usr/bin/env bash
set -euo pipefail

# find git root
pushd "$(git rev-parse --show-toplevel)" >/dev/null

CR_PREAMBLE_YAML=$(
	cat <<-END
		apiVersion: fluxninja.com/v1alpha1
		kind: Policy
		metadata:
		  name: service1
	END
)

# append contents of blueprints/blueprints/latency-gradient/example/gen/policies/example.yaml to CR_PREAMBLE_YAML in the next line under spec key
yq eval-all 'select(fileIndex == 0) * {"spec": select(fileIndex == 1)}' <(echo "$CR_PREAMBLE_YAML") blueprints/blueprints/latency-gradient/example/gen/policies/example.yaml >operator/config/samples/fluxninja.com_v1alpha1_policy.yaml

popd >/dev/null
