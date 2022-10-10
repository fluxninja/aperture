# aperture-libsonnet

This directory contains automatically generated Jsonnet Aperture Library.

## Introduction

Aperture Jsonnet Library is a set of jsonnet bindings for Aperture policies that
is generated from the OpenAPI V2/Swagger API specification that is distributed
as part of the fluxninja/aperture repository under
`docs/gen/openapiv2/aperture.swagger.yaml`.

## Requirements

Other than the [jsonnet compiler][jsonnet] the only other tool used by this
library is [jsonnet bundler][jb] that allows for freezing a specific version of
this library as a dependency in another project.

[jsonnet]: https://github.com/google/go-jsonnet
[jb]: https://github.com/jsonnet-bundler/jsonnet-bundler

## Usage

Install and freeze library with `jb`:

```sh
$ jb install github.com/fluxninja/aperture/libsonnet/1.0@main
```

Then import it in your project:

```jsonnet
local aperture = import 'github.com/fluxninja/aperture/libsonnet/1.0/main.libsonnet';
```

### On circuit components

All circuit components have a `new()` function that should be called to
initialize an empty component definition - this will allow jsonnet to verify
that for example all require in/out ports are properly set in the final object:

```jsonnet

// Don't do this:
local component = Component.withPromql(PromQL.withQueryString("vector(1)"));

// Instead, do the following:
local component = Component.withPromql(PromQL.new() + PromQL.withQueryString("vector(1)"));
```

## Development

There are two parts to this library: automatically generated bindings based on
the swagger API specification, and a set of custom bindings that build upon
those generated bindings and extend the API to provide some "syntactic sugar"
and make the library easier to use.

To update auto generated bindings, the script `jsonnet-lib-gen.py` can be used
after its dependencies (listed in `requirements.txt`) are installed. The script
requires Python 3.10+.

Example usage:

```sh
python scripts/jsonnet-lib-gen.py --output-dir 1.0/ ../docs/gen/policies/gen.yaml
```

This will update bindings from `1.0/` directory, adding and removing files in
process, to match changes in the Swagger API specification. It will also
incorporate any jsonnet custom bindings from `1.0/custom/` into the library and
make them available to the user.

Custom bindings can be used to enchant the auto generated code with additional,
more ergonomic API, that could not be generated automatically based on the
specification.
