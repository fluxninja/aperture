#!/usr/bin/env bash
set -euo pipefail

git_root=$(git rev-parse --show-toplevel)
blueprints_root=${git_root}/blueprints

FIND="find"
SED="sed"

if [[ "$OSTYPE" == "darwin"* ]]; then
	FIND="gfind"
	SED="gsed"
fi

# run jb install in the blueprints_root
pushd "${blueprints_root}" >/dev/null
jb install
popd >/dev/null

# remove md files in "${blueprints_root}"/../docs/content/reference/policies/bundled-blueprints' subdirectories
$FIND "${blueprints_root}"/../docs/content/reference/policies/bundled-blueprints -mindepth 2 -type f -name '*.md' -delete

# for all subdirectories within "$blueprints_root"/lib containing config.libsonnet, generate README
$FIND "$blueprints_root" -type f -name config.libsonnet | while read -r files; do
	dir=$(dirname "$files")
	echo "Generating README and Sample Values for $dir"
	python "${git_root}"/scripts/blueprint-readme-generator.py "$dir"
	npx prettier --write "$dir"/README.md
	npx prettier --write "$dir"/values.yaml
	# extract the name of the blueprint from the path
	blueprint_name=$(basename "$dir")
	# extract the relative path from the "$blueprints_root"
	# shellcheck disable=SC2001
	relative_path=$(echo "$dir" | $SED "s|$blueprints_root||")
	# remove the last dir from $relative_path
	relative_path=$(dirname "$relative_path")
	docs_dir="${blueprints_root}"/../docs/content/reference/policies/bundled-blueprints/"$relative_path"
	mkdir -p "$docs_dir"
	docs_file="$docs_dir"/"$blueprint_name".md
	# generate docs
	echo "Generating $blueprint_name.md"
	# generate docusaurus frontmatter
	echo "---" >"$docs_file"
	# title is the first line of the README.md, strip the leading '#'
	# shellcheck disable=SC2129
	echo "title: $(head -n 1 "$dir"/README.md | $SED 's/# //')" >>"$docs_file"
	# end of frontmatter
	echo "---" >>"$docs_file"
	# new line
	echo -e "\n" >>"$docs_file"
	# add mdx code block
	echo -e "\`\`\`mdx-code-block" >>"$docs_file"

	# count the number of levels in the relative_path (e.g. /policies/static-rate-limiting is 2 levels)
	num_levels=$(echo "$relative_path" | $SED 's/[^/]//g' | wc -c)

	# add the correct number of ../ to the import statement and additional 3 (reference/policies/bundled-blueprints) levels to find apertureVersion.js
	import_levels="../../../"
	for ((i = 1; i < num_levels; i++)); do
		import_levels="../$import_levels"
	done

	# shellcheck disable=SC2129
	echo -e "import {apertureVersion} from '${import_levels}apertureVersion.js';" >>"$docs_file"
	echo -e "\`\`\`" >>"$docs_file"
	echo "## Blueprint Location" >>"$docs_file"
	echo -e "\n" >>"$docs_file"
	echo "GitHub: <a href={\`https://github.com/fluxninja/aperture/tree/\${apertureVersion}/blueprints/$relative_path/$blueprint_name\`}>$blueprint_name</a>" >>"$docs_file"
	echo -e "\n" >>"$docs_file"
	echo "## Introduction" >>"$docs_file"
	echo -e "\n" >>"$docs_file"
	# copy README.md to the blueprint named md file, except the first line (title)
	tail -n +2 "$dir"/README.md >>"$docs_file"
	# run prettier on the docs file
	npx prettier --write "$docs_file"
done
