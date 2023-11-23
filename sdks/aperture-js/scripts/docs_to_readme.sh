#!/usr/bin/env bash

# This script is used to generate the README.md file from the docs folder
# 1. Copy contents from docs/README.md
# 2. Update all links to point to the docs folder
# 3. Remove the first line of the README.md file
# 4. Replace the #API Reference section in README.md

# Copy contents from docs/README.md
cp docs/README.md README.md.tmp

# Add docs/ in front of any link that starts with classes, enums, and interfaces
sed -i -e 's/(classes/(docs\/classes/g' README.md.tmp
sed -i -e 's/(enums/(docs\/enums/g' README.md.tmp
sed -i -e 's/(interfaces/(docs\/interfaces/g' README.md.tmp

# Remove the first line of the README.md file
sed -i -e '1d' README.md.tmp

# Remove the #API Reference section and everything under it from README.md
sed -i -e '/# API Reference/,$d' README.md

# Append the new #API Reference section to README.md
echo -e "\n# API Reference\n$(cat README.md.tmp)" >> README.md

# Remove temporary file
rm README.md.tmp README.md-e README.md.tmp-e

# Format README.md
prettier --write README.md
