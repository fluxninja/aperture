#!/usr/bin/env bash
set -euo pipefail

script_root=$(dirname "$0")
blueprints_root=${script_root}/..

FIND="find"

if [[ "$OSTYPE" == "darwin"* ]]; then
	FIND="gfind"
fi

# accept a --force flag to force the regeneration of all graphs
# check if $1 is bounded and equal to --force
if [[ -n "${1:-}" ]] && [[ "$1" == "--force" ]]; then
	force=true
else
	force=false
fi

# run jb install in the blueprints_root
pushd "${blueprints_root}" >/dev/null
jb install
popd >/dev/null

# remove md files in "${blueprints_root}"/../docs/content/references/bundled-blueprints and subdirectories
$FIND "${blueprints_root}"/../docs/content/references/bundled-blueprints -type f -name '*.md' -delete

# for all subdirectories within "$blueprints_root"/lib containing config.libsonnet, generate README
$FIND "$blueprints_root"/lib/1.0 -type f -name config.libsonnet | while read -r files; do
	dir=$(dirname "$files")
	echo "Generating README for $dir"
	python "${blueprints_root}"/scripts/blueprint-readme-generator.py "$dir"
	npx prettier --write "$dir"/README.md
	# extract the name of the blueprint from the path
	blueprint_name=$(basename "$dir")
	# extract the relative path from the "$blueprints_root"/lib/1.0
	# shellcheck disable=SC2001
	relative_path=$(echo "$dir" | sed "s|$blueprints_root/lib/1.0/||")
	blueprint_dir=$(dirname "$relative_path")
	docs_dir="${blueprints_root}"/../docs/content/references/bundled-blueprints/"$blueprint_dir"
	mkdir -p "$docs_dir"
	docs_file="$docs_dir"/"$blueprint_name".md
	# generate docs
	echo "Generating $blueprint_name.md"
	# generate docusaurus frontmatter
	echo "---" >"$docs_file"
	# title is the first line of the README.md, strip the leading '#'
	# shellcheck disable=SC2129
	echo "title: $(head -n 1 "$dir"/README.md | sed 's/# //')" >>"$docs_file"
	# end of frontmatter
	echo "---" >>"$docs_file"
	# new line
	echo -e "\n" >>"$docs_file"
	# add mdx code block
	echo -e "\`\`\`mdx-code-block" >>"$docs_file"
	echo -e "import {apertureVersion} from '../../../apertureVersion.js';" >>"$docs_file"
	echo -e "\`\`\`" >>"$docs_file"
	echo "## Blueprint Location" >>"$docs_file"
	echo -e "\n" >>"$docs_file"
	echo "GitHub: <a href={\`https://github.com/fluxninja/aperture/tree/\${apertureVersion}/blueprints/lib/1.0/$blueprint_dir/$blueprint_name\`}>$blueprint_name</a>" >>"$docs_file"
	echo -e "\n" >>"$docs_file"
	echo "## Introduction" >>"$docs_file"
	echo -e "\n" >>"$docs_file"
	# copy README.md to the blueprint named md file, except the first line (title)
	tail -n +2 "$dir"/README.md >>"$docs_file"
	# run prettier on the docs file
	npx prettier --write "$docs_file"
done

# for all directories within "$blueprints_root"/bundles, generate examples
$FIND "$blueprints_root"/examples -mindepth 1 -maxdepth 1 -type d | while read -r dir; do
	echo "Generating examples for $dir"
	# save the contents of $dir/gen/policies/example.yaml for comparison
	old_example_yaml=$(cat "$dir"/gen/policies/example.yaml 2>/dev/null || true)

	# generate example blueprint
	python "${blueprints_root}"/scripts/generate-bundle.py --output "$dir"/gen/ \
		--config "$dir"/example.jsonnet

	npx prettier --write "$dir"/gen/.. || true

	new_example_yaml=$(cat "$dir"/gen/policies/example.yaml)

	if [[ "$old_example_yaml" != "$new_example_yaml" || "$force" == true ]]; then
		mkdir -p "$dir"/gen/graph
		# fail if commands below fails
		go run -mod=mod "${blueprints_root}"/../cmd/circuit-compiler/main.go \
			-cr "$dir"/gen/policies/example.yaml \
			-dot "$dir"/gen/graph/graph.dot
		# if exit code is not 0 then remove example.yaml
		# shellcheck disable=SC2181
		if [[ $? -ne 0 ]]; then
			rm "$dir"/gen/policies/example.yaml
			# exit script with error
			exit 1
		fi
		dot -Tsvg "$dir"/gen/graph/graph.dot >"$dir"/gen/graph/graph.svg
	fi
done
