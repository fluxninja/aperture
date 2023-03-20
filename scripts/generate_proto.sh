#!/usr/bin/env bash
set -euo pipefail

script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
git_root="$(git -C "$script_dir" rev-parse --show-toplevel)"
readonly script_dir="${script_dir}"
readonly git_root="${git_root}"

readonly api_dir="$git_root/api"

proto_files=(
  "aperture/flowcontrol/check/v1/check.proto"
)


# Protoc for Python generates imports from the directory structure of the proto files,
# ignoring package structure.
#
# This little bit of sed rewrites the broken
# imports that protoc produces, of the form
#     import foo_pb2 as bar
# to the form
#     from . import foo_pb2 as bar
#
# mypy_protoc doesn't yet support aio which we use (https://github.com/nipunn1313/mypy-protobuf/issues/216) so
# we 'sed' to do some replacements to support that, adding an import for grpc.aio and typing and changing
# grpc.Server to Union[grpc.Server, grpc.aio.Server].
#
# mypy_protoc also marks all server methods as abstract, they aren't really abstract because they have default
# implementations that simply return a grpc error that it's not implemented, we make use of those default
# methods so don't want type checking telling us to fill them all in.

readonly python_gen="${git_root}/sdks/aperture-py/aperture_sdk/_gen"
echo "Installing Python gRPC tools"
pip install grpcio-tools mypy-protobuf

echo "Cleanining up old generated files"
rm -rf "$python_gen"
mkdir -p "$python_gen"

echo "Generating Python gRPC code"
python3 -m grpc_tools.protoc -I"${api_dir}" --{python_out,mypy_out,grpc_python_out,mypy_grpc_out}="${python_gen}" "${proto_files[@]}"

echo "Fixing up Python gRPC imports"
generated_py_files_str="$(find "$python_gen" -type f -name '*.py')"
readarray -t generated_py_files <<< "$generated_py_files_str"
sed -i "s/^from aperture\..* import \([^ ]*_pb2\) as \([^ ]*\)$/from . import \1 as \2/" "${generated_py_files[@]}"
generated_pyi_files_str="$(find "$python_gen" -type f -name '*.pyi')"
readarray -t generated_pyi_files <<< "$generated_pyi_files_str"
sed -i "s/^import grpc$/import grpc\nimport grpc.aio\nimport typing/" "${generated_pyi_files[@]}"
sed -i "s/: grpc\.Server/: typing.Union[grpc.Server, grpc.aio.Server]/" "${generated_pyi_files[@]}"

echo "Generating __init__.py files"
find "$python_gen" -type d -exec touch "{}/__init__.py" \;

echo "Done"
