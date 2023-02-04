#!/usr/bin/env bash
set -euo pipefail

script_root=$(dirname "$0")
blueprints_root=${script_root}/..

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
$FIND "$blueprints_root"/lib/1.0 -type f -name config.libsonnet | while read -r files; do
	dir=$(dirname "$files")
	echo "Generating README for $dir"
	python "${blueprints_root}"/scripts/blueprint-readme-generator.py "$dir"
	npx prettier --write "$dir"/README.md
	# extract the name of the blueprint from the path
	blueprint_name=$(basename "$dir")
	# extract the relative path from the "$blueprints_root"/lib/1.0
	# shellcheck disable=SC2001
	relative_path=$(echo "$dir" | $SED "s|$blueprints_root/lib/1.0/||")
	blueprint_dir=$(dirname "$relative_path")
	docs_dir="${blueprints_root}"/../docs/content/reference/policies/bundled-blueprints/"$blueprint_dir"
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

	# count the number of levels in the blueprint_dir
	num_levels=$(echo "$blueprint_dir" | $SED 's/[^/]//g' | wc -c)

	# add the correct number of ../ to the import statement and additional 3 (reference/policies/bundled-blueprints) levels to find apertureVersion.js
	import_levels="../../../"
	for ((i = 0; i < num_levels; i++)); do
		import_levels="../$import_levels"
	done

	# shellcheck disable=SC2129
	echo -e "import {apertureVersion} from '${import_levels}apertureVersion.js';" >>"$docs_file"
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
