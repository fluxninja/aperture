#!/bin/bash

# Example
# FORMAT=png ./docs/tools/mermaid/generate_mermaid.sh ./docs/content/applying-policies/load-scheduling/assets/workload-prioritization/graph.mmd

gitroot=$(git rev-parse --show-toplevel)
docsdir=$gitroot/docs
# default to svg FORMAT if not provided
FORMAT=${FORMAT:-svg}

set -euo pipefail
mmd_file="$1"
echo "generating $FORMAT file for $mmd_file"
mmdc --quiet --input "$mmd_file" --configFile "$docsdir"/tools/mermaid/mermaid-theme.json --cssFile "$docsdir"/tools/mermaid/mermaid.css --scale 2 --output "$mmd_file"."$FORMAT" --backgroundColor transparent
