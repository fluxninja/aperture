#!/usr/bin/env bash

# merge plugin swaggers
set -e

git_root=$(git rev-parse --show-toplevel)
docs_root="$git_root"/docs

FIND="find"
if [[ "$OSTYPE" == "darwin"* ]]; then
	FIND="gfind"
fi

cp "$docs_root"/gen/config/extensions/fluxninja/extension-swagger.yaml merged-extension-swagger.yaml
dirs=$(find . -name 'extension-swagger.yaml' -exec dirname {} \;)
for dir in $dirs; do
	echo "Merging $dir/extension-swagger.yaml"
	yq eval-all --inplace "select(fileIndex==0).definitions *= select(fileIndex==1).definitions | select(fileIndex==0)" merged-extension-swagger.yaml "$dir"/extension-swagger.yaml
	yq eval-all --inplace "select(fileIndex==0).paths *= select(fileIndex==1).paths | select(fileIndex==0)" merged-extension-swagger.yaml "$dir"/extension-swagger.yaml
done
dirs=$($FIND "$docs_root"/gen/config -name 'config-swagger.yaml' -exec dirname {} \;)
for dir in $dirs; do
	echo generating markdown for "$dir"/config-swagger.yaml
	basename=$(basename "$dir")
	cp "$dir"/config-swagger.yaml "$dir"/gen.yaml
	yq eval-all --inplace "select(fileIndex==0).definitions *= select(fileIndex==1).definitions | select(fileIndex==0)" "$dir"/gen.yaml merged-extension-swagger.yaml
	yq eval-all --inplace "select(fileIndex==0).paths *= select(fileIndex==1).paths | select(fileIndex==0)" "$dir"/gen.yaml merged-extension-swagger.yaml
	swagger flatten \
		--with-flatten=remove-unused "$dir"/gen.yaml \
		--format=yaml --output "$dir"/gen.yaml
	swagger generate markdown \
		--spec "$dir"/gen.yaml \
		--target "$dir" \
		--skip-validation \
		--quiet \
		--with-flatten=remove-unused \
		--tags=common-configuration \
		--tags=extension-configuration \
		--tags=agent-configuration \
		--tags=controller-configuration \
		--allow-template-override \
		--template-dir "$docs_root"/tools/swagger/swagger-templates \
		--config-file "$docs_root"/tools/swagger/markdown-config.yaml \
		--output "$basename".md
	rm "$dir"/gen.yaml
	cat "$dir"/metadata "$dir"/"$basename".md >"$dir"/"$basename".md.tmp
	mv "$dir"/"$basename".md.tmp "$dir"/"$basename".md
	npx prettier --prose-wrap="preserve" --write "$dir"/"$basename".md
	mv "$dir"/"$basename".md "$docs_root"/content/reference/configuration
done
rm merged-extension-swagger.yaml

# policy markdown
echo generating policy markdown
policy_dir="$docs_root"/gen/policy
cp "$docs_root"/content/assets/openapiv2/aperture.swagger.yaml "$policy_dir"/
yq -i eval 'del(.paths)' "$policy_dir"/aperture.swagger.yaml
yq -i eval 'del(.tags)' "$policy_dir"/aperture.swagger.yaml
# 'mixin' is mostly used for --keep-spec-order
swagger mixin "$policy_dir"/config-swagger.yaml "$policy_dir"/aperture.swagger.yaml --keep-spec-order --format=yaml -o "$policy_dir"/policy.yaml
# Fixup .info, which is altered by 'mixin'
yq -i eval-all 'select(fileIndex == 0).info = select(fileIndex == 1).info' \
	"$policy_dir"/policy.yaml "$policy_dir"/config-swagger.yaml
swagger flatten --with-flatten=remove-unused "$policy_dir"/policy.yaml --format=yaml --output "$policy_dir"/policy.yaml
swagger generate markdown --spec "$policy_dir"/policy.yaml --target "$policy_dir" \
	--skip-validation \
	--quiet \
	--with-flatten=remove-unused \
	--tags=policy-configuration \
	--allow-template-override --template-dir "$docs_root"/tools/swagger/swagger-templates \
	--config-file "$docs_root"/tools/swagger/markdown-config.yaml
rm "$policy_dir"/aperture.swagger.yaml
# append gen/policy/metadata on top of gen/policy/policy.md
cat "$policy_dir"/metadata "$policy_dir"/policy.md >"$policy_dir"/policy.md.tmp
mv "$policy_dir"/policy.md.tmp "$policy_dir"/policy.md
npx prettier --prose-wrap="preserve" --write "$policy_dir"/policy.md
mv "$policy_dir"/policy.md "$docs_root"/content/reference/policies/spec.md
