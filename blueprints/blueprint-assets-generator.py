#!/usr/bin/env python3
#
# This script is used to generate the assets for the blueprints
# - values.yaml
# - dynamic-values.yaml
# - json schema
# - documentation
from __future__ import annotations

import dataclasses
import json
import re
import subprocess
import sys
import tempfile
from pathlib import Path
from typing import Any, Dict, List, Optional

import jinja2
import typer
import yaml
from loguru import logger
from ordered_set import OrderedSet
from slugify import slugify


class ExitException(RuntimeError):
    ret_code: int

    def __init__(self, ret_code=0):
        self.ret_code = ret_code


ANNOTATION_RE = re.compile(r".*@param.*|.*@schema.*")
ANNOTATION_DETAILED_RE = re.compile(
    r".*(?P<annotation_type>@schema|@param) \((?P<param_name>[\w.\[\]]+): (?P<param_type>[\w.\[\]/:-]+) ?(?P<param_required>\w+)?\) ?(?P<param_description>.+)?"
)


@dataclasses.dataclass
class Parameter:
    annotation_type: str = "@param"
    param_name: str = ""
    param_type: str = "intermediate_node"
    is_complex_type: bool = False
    json_schema_link: str = ""
    docs_link: str = ""
    description: str = ""
    required: bool = False
    default: Optional[Any] = None


@dataclasses.dataclass
class ParameterNode:
    parameter: Parameter = dataclasses.field(default_factory=Parameter)
    children: Dict[str, ParameterNode] = dataclasses.field(default_factory=dict)
    required_children: OrderedSet[str] = dataclasses.field(default_factory=OrderedSet)


@dataclasses.dataclass
class Blueprint:
    # nested dictionary of parameters
    nested_parameters: ParameterNode = dataclasses.field(default_factory=ParameterNode)
    # nested dictionary of required parameters
    nested_required_parameters: ParameterNode = dataclasses.field(
        default_factory=ParameterNode
    )
    # deprecated is a string
    deprecation_message: Optional[str] = None

    @classmethod
    def _resolve_param_links(
        cls, blueprints_root_relative_path: str, policies_relative_path: str, param: str
    ) -> tuple[str, str, bool]:
        # if param starts with [] then it is an array of objects, remove the []
        if param.startswith("[]"):
            param = param[2:]

        if param.startswith("aperture.spec"):
            component = param.split(".")[-1]
            component_slug = camel_to_kebab_case(component)
            docs_link = f"{policies_relative_path}/spec#{component_slug}"
            json_schema_link = f"{blueprints_root_relative_path}/gen/jsonschema/_definitions.json#/definitions/{component}"
            return (docs_link, json_schema_link, True)
        elif param.count(":") > 0:
            if param.count(":") != 2:
                logger.error(
                    f"Unable to parse @param: `{param}`. Expecting format <blueprint>:<annotation_type>:<param>"
                )
                raise typer.Exit(1)
            # parameter is of the form <blueprint>:<annotation_type>:<param>
            blueprint, annotation_type, param = param.split(":")
            docs_link = f"{policies_relative_path}/bundled-blueprints/{blueprint}#{slugify(param)}"
            parts = param.split(".")
            json_schema_link = (
                f"{blueprints_root_relative_path}/{blueprint}/gen/definitions.json#"
            )
            if annotation_type == "param":
                json_schema_link += f"/properties/{parts[0]}"
            elif annotation_type == "schema":
                json_schema_link += f"/$defs/{parts[0]}"
            else:
                logger.error(
                    f"Unable to parse @param: `{param}`. Expecting format <blueprint>:<annotation_type>:<param>. <annotation_type> must be either @param or @schema."
                )
                raise typer.Exit(1)

            for part in parts[1:]:
                json_schema_link += f"/properties/{part}"
            return (docs_link, json_schema_link, True)

        # check if it's not a base type bool, float32, float64, int32, int64, string
        elif param not in [
            "bool",
            "float32",
            "float64",
            "int32",
            "int64",
            "string",
        ]:
            # local schema
            docs_link = f"#{slugify(param)}"
            parts = param.split(".")
            json_schema_link = f"#/$defs/{parts[0]}"
            for part in parts[1:]:
                json_schema_link += f"/properties/{part}"
            return (docs_link, json_schema_link, True)
        else:
            return "", "", False

    @classmethod
    def from_comment(
        cls,
        blueprints_root_relative_path: str,
        policies_relative_path: str,
        comment: List[str],
    ) -> Blueprint:
        nested_parameters = ParameterNode()
        nested_required_parameters = ParameterNode()
        for line in comment:
            if ANNOTATION_RE.match(line.strip()):
                inner = ANNOTATION_DETAILED_RE.match(line.strip())
                if not inner:
                    inner = ""
                    logger.error(f"Unable to parse @param: `{line.strip()}`")
                    raise typer.Exit(1)

                groups = inner.groupdict()
                param_name, param_type = groups["param_name"], groups["param_type"]
                (
                    docs_link,
                    json_schema_link,
                    is_complex_type,
                ) = cls._resolve_param_links(
                    blueprints_root_relative_path, policies_relative_path, param_type
                )

                annotation_type = groups["annotation_type"]
                param_required = groups.get("param_required", "") == "required"
                param_description = groups.get("param_description", "")
                # tokenize param_name and create nested_parameters
                parts = param_name.split(".")
                parent = nested_parameters.children
                parent_required = nested_required_parameters.children
                if param_required and annotation_type == "@param":
                    nested_parameters.required_children.add(parts[0])
                    nested_required_parameters.required_children.add(parts[0])
                for idx, part in enumerate(parts):
                    if idx == len(parts) - 1:
                        node = ParameterNode(
                            Parameter(
                                annotation_type,
                                part,
                                param_type,
                                is_complex_type,
                                json_schema_link,
                                docs_link,
                                param_description,
                                param_required,
                            )
                        )
                        parent[part] = node
                        if param_required:
                            parent_required[part] = node
                    else:
                        next_part = parts[idx + 1]
                        node = ParameterNode(Parameter(annotation_type, part))
                        if part not in parent:
                            parent[part] = node
                        if param_required:
                            parent[part].required_children.add(next_part)
                        parent = parent[part].children
                        if param_required:
                            if part not in parent_required:
                                parent_required[part] = node
                            parent_required[part].required_children.add(next_part)
                            parent_required = parent_required[part].children

        if not nested_parameters or not nested_required_parameters:
            logger.error("Unable to find parameters in comments")
            raise ValueError()

        return cls(
            nested_parameters,
            nested_required_parameters,
        )


def command_with_exit_code(func):
    def wrapper():
        try:
            func()
        except ExitException as ex:
            sys.exit(ex.ret_code)

    return wrapper


def update_param_defaults(
    repository_root: Path,
    config_path: Path,
    parameters: Blueprint,
    jsonnet_path: Path = Path(),
):
    jsonnet_data = f"local config = import '{config_path}';\n"
    if jsonnet_path != Path():
        jsonnet_data += f"local fn = import '{jsonnet_path}';\n"

    if jsonnet_path != Path():
        jsonnet_data += f"fn(config)\n"
    jsonnet_data += "{_config::: config}\n"

    rendered_config = None
    with tempfile.NamedTemporaryFile(suffix=".libsonnet") as tmp:
        tmppath = Path(tmp.name)
        tmppath.write_text(jsonnet_data)

        jsonnet_jpaths = [
            "-J",
            repository_root / "blueprints",
            "-J",
            repository_root / "blueprints" / "vendor",
        ]

        try:
            result = subprocess.run(
                ["jsonnet", *jsonnet_jpaths, str(tmppath)],
                capture_output=True,
                check=True,
            )
        except subprocess.CalledProcessError as ex:
            logger.error(f"Error while rendering jsonnet: {ex.stderr}")
            # log file for debugging
            logger.error(f"Jsonnet file: {jsonnet_data}")
            raise typer.Exit(1)

        rendered_config = json.loads(result.stdout)

    def get_param_default_from_rendered_config(root: Dict, name: str) -> Any:
        parts = name.split(".")
        config = root
        for idx, part in enumerate(parts):
            if idx == len(parts) - 1:
                try:
                    return config[part]
                except KeyError:
                    # fatal exit
                    logger.error(f"Unable to find param {name} in rendered config")
                    raise typer.Exit(1)
            else:
                try:
                    config = config[part]
                except KeyError:
                    # the param is not present in the config, so we return None
                    # Also, when specific param is a map (map[string]type) and there is no default
                    # then we return None here, which will be converted into an empty map later.
                    return None

    logger.trace(rendered_config)
    # walk nested_parameters and update defaults

    def update_nested_param_defaults(node, prefix=""):
        if node.parameter.param_type != "intermediate_node":
            default = get_param_default_from_rendered_config(
                rendered_config["_config"], prefix
            )
            if default is not None:
                node.parameter.default = default
        for key, child in node.children.items():
            if prefix != "":
                keyPrefix = f"{prefix}.{key}"
            else:
                keyPrefix = key
            update_nested_param_defaults(child, keyPrefix)

    update_nested_param_defaults(parameters.nested_parameters)
    update_nested_param_defaults(parameters.nested_required_parameters)


MARKDOWN_DOC_TPL = """
{%- macro render_type(param_type, is_complex_type) %}
{%- if param_type.startswith('[]') %}
{{- 'Array of ' + render_type(param_type[2:], is_complex_type) }}
{%- elif is_complex_type %}
{{- 'Object (' + param_type + ')' }}
{%- elif param_type == 'bool' %}
{{- 'Boolean' }}
{%- elif param_type == 'float32' %}
{{- 'Number (float)' }}
{%- elif param_type == 'float64' %}
{{- 'Number (double)' }}
{%- elif param_type == 'int32' %}
{{- 'Integer (int32)' }}
{%- elif param_type == 'int64' %}
{{- 'Integer (int64)' }}
{%- else %}
{{- param_type }}
{%- endif %}
{%- endmacro %}

{%- macro render_node(node, level, annotation_type, parent_prefix='') %}
{%- set indent = '    ' * level %}
{%- set anchor = (parent_prefix + node.parameter.param_name) | slugify %}
{%- set heading_level = '#' * (level + 1) %}
{%- if node.parameter.param_type == 'intermediate_node' %}
<!-- vale off -->

{{ heading_level }} {{ parent_prefix if annotation_type == '@param' }}{{ node.parameter.param_name }} {#{{ anchor }}}

<!-- vale on -->

{%- for child_name, child_node in node.children.items() if child_node.parameter.param_type != 'intermediate_node' %}
{{ render_node(child_node, level + 1, annotation_type,
               parent_prefix + node.parameter.param_name + '.') }}
{%- endfor %}
{%- for child_name, child_node in node.children.items() if child_node.parameter.param_type == 'intermediate_node' %}
{{ render_node(child_node, level + 1, annotation_type,
               parent_prefix + node.parameter.param_name + '.') }}
{%- endfor %}
{%- else %}
<!-- vale off -->

<a id="{{ anchor }}"></a>

<ParameterDescription
    name='{{ parent_prefix if annotation_type == '@param' }}{{ node.parameter.param_name }}'
    description='{{ node.parameter.description }}'
    type='{{- render_type(node.parameter.param_type, node.parameter.is_complex_type) }}'
    reference='{{ node.parameter.docs_link }}'
    value='{{ node.parameter.default | to_json }}'
/>

<!-- vale on -->
{%- endif %}
{%- endmacro %}

### Parameters

{%- for child_name, child_node in nested_parameters.children.items() %}
{%- if child_node.parameter.annotation_type == '@param' %}
{{ render_node(child_node, 3, '@param') }}

---

{%- endif %}
{%- endfor %}
{%- set ns = namespace(has_schema=False) %}
{%- for child_name, child_node in nested_parameters.children.items() %}
{%- if child_node.parameter.annotation_type == '@schema' %}
{%- set ns.has_schema = true %}
{%- endif %}
{%- endfor %}
{%- if ns.has_schema %}

### Schemas

{%- for child_name, child_node in nested_parameters.children.items() %}
{%- if child_node.parameter.annotation_type == '@schema' %}
{{ render_node(child_node, 3, '@schema') }}

---

{%- endif %}
{%- endfor %}
{%- endif %}
"""

MARKDOWN_README_TPL = """
{%- macro render_type(param_type, is_complex_type) %}
{%- if param_type.startswith('[]') %}
{{- 'Array of ' + render_type(param_type[2:], is_complex_type) }}
{%- elif is_complex_type %}
{{- 'Object (' + param_type + ')' }}
{%- elif param_type == 'bool' %}
{{- 'Boolean' }}
{%- elif param_type == 'float32' %}
{{- 'Number (float)' }}
{%- elif param_type == 'float64' %}
{{- 'Number (double)' }}
{%- elif param_type == 'int32' %}
{{- 'Integer (int32)' }}
{%- elif param_type == 'int64' %}
{{- 'Integer (int64)' }}
{%- else %}
{{- param_type }}
{%- endif %}
{%- endmacro %}

{%- macro render_properties(node, level, annotation_type, parent_prefix='') %}
{%- set indent = '    ' * level %}
{%- set anchor = (parent_prefix + node.parameter.param_name) | slugify %}
{%- set heading_level = '#' * (level + 1) %}
{%- if node.parameter.param_type == 'intermediate_node' %}
{{ heading_level }} {{ parent_prefix if annotation_type == '@param' }}{{ node.parameter.param_name }} {#{{ anchor }}}

{%- if node.parameter.description %}
**Description**: {{ node.parameter.description }}
{%- endif %}
{%- for child_name, child_node in node.children.items() %}
{{ render_properties(child_node, level + 1, annotation_type,
                     parent_prefix + node.parameter.param_name + '.') }}
{%- endfor %}
{%- else %}
{{ heading_level }} {{ parent_prefix if annotation_type == '@param' }}{{ node.parameter.param_name }} {#{{ anchor }}}
**Description**: {{ node.parameter.description }}
**Type**: {{- render_type(node.parameter.param_type, node.parameter.is_complex_type) }}
**Default Value**:
<details>
<summary>Click to expand</summary>
```yaml
{{ node.parameter.default | to_yaml }}
```
</details>
{%- endif %}
{%- endmacro %}

### Parameters

{%- for child_name, child_node in nested_parameters.children.items() %}
{%- if child_node.parameter.annotation_type == '@param' %}
{{ render_properties(child_node, 3, '@param') }}
{%- endif %}
{%- endfor %}
{%- set ns = namespace(has_schema=False) %}
{%- for child_name, child_node in nested_parameters.children.items() %}
{%- if child_node.parameter.annotation_type == '@schema' %}
{%- set ns.has_schema = true %}
{%- endif %}
{%- endfor %}
{%- if ns.has_schema %}

### Schemas

{%- for child_name, child_node in nested_parameters.children.items() %}
{%- if child_node.parameter.annotation_type == '@schema' %}
{{ render_properties(child_node, 3, '@schema') }}
{%- endif %}
{%- endfor %}
{%- endif %}
"""


JSON_SCHEMA_TPL = """
{% macro render_type(param_type, ref_id, is_complex_type) %}
{% if param_type.startswith('[]') %}
type: array
items:
  {{ render_type(param_type[2:], ref_id, is_complex_type) | indent(2) }}
{% elif param_type.startswith('map[') %}
type: object
additionalProperties: true
{% elif is_complex_type %}
type: object
$ref: "{{- ref_id }}"
{% elif param_type == 'bool' %}
type: boolean
{% elif param_type == 'float32' %}
type: number
format: float
{% elif param_type == 'float64' %}
type: number
format: double
{% elif param_type == 'int32' %}
type: integer
format: int32
{% elif param_type == 'int64' %}
type: integer
format: int64
{% else %}
type: "{{ param_type }}"
{% endif %}
{% endmacro %}
{% macro render_properties(node, annotation_type, prefix='') %}
{% if node.parameter.annotation_type == annotation_type %}
{% if node.parameter.param_type == 'intermediate_node' %}
{{ node.parameter.param_name }}:
  type: object
  additionalProperties: false
  {% if node.required_children %}
  required:
  {% for child_name in node.required_children %}  - {{ child_name }}
  {% endfor %}
  {% endif %}
  properties:
  {% for child_name, child_node in node.children.items() %}
    {{ render_properties(child_node, annotation_type, prefix ~ node.parameter.param_name ~ '_') | indent(4) }}
  {% endfor %}
{% else %}
{{ node.parameter.param_name }}:
  description: "{{ node.parameter.description }}"
  default: {{ node.parameter.default | quote_value }}
  {{ render_type(node.parameter.param_type, node.parameter.json_schema_link,
                 node.parameter.is_complex_type) | indent(2, true) }}
{% endif %}
{% endif %}
{% endmacro %}
$schema: "http://json-schema.org/draft-07/schema#"
type: object
title: "{{ blueprint_name }} blueprint"
additionalProperties: false
{% if nested_parameters.required_children %}
required:
{% for child_name in nested_parameters.required_children %}- {{ child_name }}
{% endfor %}
{% endif %}
properties:
{% for child_name, child_node in nested_parameters.children.items() %}
  {{ render_properties(child_node, '@param') | indent(2) }}
{% endfor %}
$defs:
{% for child_name, child_node in nested_parameters.children.items() %}
    {{ render_properties(child_node, '@schema', '') | indent(2) }}
{% endfor %}
"""

JSON_SCHEMA_DEFINITIONS_TPL = """

"""

YAML_TPL = """
# Generated values file for {{ blueprint_name }} blueprint
# Documentation/Reference for objects and parameters can be found at:
# https://docs.fluxninja.com/reference/policies/bundled-blueprints/{{ blueprint_name }}
{%- macro render_value(value, level) %}
{%- if value is mapping %}
{%- for key, val in value.items() %}
{{ '  ' * (level) }}{{ key }}: {{ render_value(val, level+1) }}
{%- endfor %}
{%- elif value is iterable and value is not string %}
{%- if value | length == 0 %}
{{- '[]' }}
{%- else %}
{%- for item in value %}
{{ '  ' * level }}- {{ render_value(item, level+1) }}
{%- endfor %}
{%- endif %}
{%- else %}
{{- value | quote_value }}
{%- endif %}
{%- endmacro %}
{%- macro render_node(node, level) %}
{%- if node.parameter.annotation_type == '@param' %}
{%- if node.parameter.param_type != 'intermediate_node' %}
{%- if node.parameter.description %}
{{ '  ' * level }}# {{ node.parameter.description }}
{%- endif %}
{{ '  ' * level }}# Type: {{ node.parameter.param_type }}
{%- if node.parameter.required != False %}
{{ '  ' * level }}# Required: {{ node.parameter.required }}
{%- endif %}
{%- endif %}
{%- if node.param_name != '' %}
{%- if node.children | length == 0 %}
{{ '  ' * level }}{{ node.parameter.param_name }}: {{ render_value(node.parameter.default, level+1) }}
{%- else %}
{{ '  ' * level }}{{ node.parameter.param_name }}:
{%- endif %}
{%- endif %}
{%- for child_name, child_node in node.children.items() %}
{{- render_node(child_node, level + 1) }}
{%- endfor %}
{%- endif %}
{%- endmacro %}
{%- for child_name, child_node in sample_config_data.children.items() %}
{{ render_node(child_node, 0) }}
{%- endfor %}
"""


def quote_value(value: str) -> str:
    # if value is __REQUIRED_FIELD__ return as unquoted string
    if value == "__REQUIRED_FIELD__":
        return value
    return json.dumps(value)


def to_yaml(value: Any) -> str:
    return yaml.dump(value, default_flow_style=False)


def to_json(value: Any) -> str:
    return json.dumps(value)


def get_jinja2_environment() -> jinja2.Environment:
    JINJA2_TEMPLATES = {
        "markdown.doc.md.j2": MARKDOWN_DOC_TPL,
        "markdown.readme.md.j2": MARKDOWN_README_TPL,
        "values.yaml.j2": YAML_TPL,
        "definitions.json.j2": JSON_SCHEMA_TPL,
    }
    loader = jinja2.DictLoader(JINJA2_TEMPLATES)
    env = jinja2.Environment(
        loader=loader, comment_start_string="<%--", comment_end_string="--%>"
    )
    env.filters["slugify"] = slugify
    env.filters["quote_value"] = quote_value
    env.filters["to_yaml"] = to_yaml
    env.filters["to_json"] = to_json
    return env


def update_readme_markdown(
    readme_path: Path,
    config_parameters: Blueprint,
    dynamic_config_parameters: Blueprint,
):
    """Find configuration marker in and add generated content below it."""
    config_marker = "<!-- Configuration Marker -->"

    if not readme_path.exists():
        # create a new file with config marker
        readme_path.write_text(config_marker)

    readme_data = readme_path.read_text()

    # add a config marker to the end of the file
    if config_marker not in readme_data:
        readme_data += f"\n{config_marker}"

    readme_copied = ""
    for line in readme_data.split("\n"):
        readme_copied += line + "\n"
        if line == config_marker:
            break

    env = get_jinja2_environment()
    template = env.get_template("markdown.readme.md.j2")
    readme_copied += template.render(
        {"nested_parameters": config_parameters.nested_parameters}
    )
    if len(dynamic_config_parameters.nested_parameters.children) > 0:
        readme_copied += "\n\n## Dynamic Configuration\n\n"
        readme_copied += template.render(
            {"nested_parameters": dynamic_config_parameters.nested_parameters}
        )
    readme_path.write_text(readme_copied)


def update_docs_markdown(
    readme_path: Path,
    config_parameters: Blueprint,
    dynamic_config_parameters: Blueprint,
    blueprint_name: Path,
    docs_root_relative_path,
):
    """Find configuration marker in and add generated content below it."""
    config_marker = "<!-- Configuration Marker -->"

    if not readme_path.exists():
        # create a new file with config marker
        readme_path.write_text(config_marker)

    readme_data = readme_path.read_text()

    # add a config marker to the end of the file
    if config_marker not in readme_data:
        readme_data += f"\n{config_marker}"

    readme_copied = ""
    for line in readme_data.split("\n"):
        readme_copied += line + "\n"
        if line == config_marker:
            break

    aperture_version_path = docs_root_relative_path + "/apertureVersion.js"
    parameter_components_path = docs_root_relative_path + "/parameterComponents.js"

    readme_copied += "```mdx-code-block\n"
    readme_copied += (
        f"import {{apertureVersion as aver}} from '{aperture_version_path}'\n"
    )
    readme_copied += (
        f"import {{ParameterDescription}} from '{parameter_components_path}'\n"
    )
    readme_copied += "```\n\n"

    # if the blueprint is deprecated, show the warning
    if config_parameters.deprecation_message:
        readme_copied += f":::danger\n"
        readme_copied += (
            f"This blueprint is deprecated and will be removed in a future release.\n"
        )
        readme_copied += f"{config_parameters.deprecation_message}\n"
        readme_copied += f"\n:::\n\n"

    readme_copied += f"## Configuration\n"
    readme_copied += f"<!-- vale off -->\n"
    readme_copied += f"\nBlueprint name: <a href={{`https://github.com/fluxninja/aperture/tree/${{aver}}/blueprints/{blueprint_name}`}}>{blueprint_name}</a>\n\n"
    readme_copied += f"<!-- vale on -->\n"

    env = get_jinja2_environment()
    template = env.get_template("markdown.doc.md.j2")
    rendered = template.render(
        {"nested_parameters": config_parameters.nested_parameters}
    )

    readme_copied += rendered
    if len(dynamic_config_parameters.nested_parameters.children) > 0:
        readme_copied += "\n\n## Dynamic Configuration\n\n"
        readme_copied += "\n\n:::note\n\n"
        readme_copied += "The following configuration parameters can be [dynamically configured](/reference/aperturectl/apply/dynamic-config/dynamic-config.md) at runtime, without reloading the policy.\n\n"
        readme_copied += ":::\n\n"
        rendered = template.render(
            {"nested_parameters": dynamic_config_parameters.nested_parameters}
        )
        readme_copied += rendered

    readme_path.write_text(readme_copied)


def render_sample_config_yaml(
    blueprint_name: Path,
    sample_config_path: Path,
    only_required: bool,
    parameters: Blueprint,
):
    """Render sample config YAML file from blocks"""
    sample_config_data = ParameterNode()
    if only_required is False:
        sample_config_data = parameters.nested_parameters
    else:
        sample_config_data = parameters.nested_required_parameters

    env = get_jinja2_environment()
    template = env.get_template("values.yaml.j2")
    rendered = template.render(
        {"sample_config_data": sample_config_data, "blueprint_name": blueprint_name}
    )
    sample_config_path.write_text(rendered)


def render_json_schema(
    blueprint_name: Path, json_schema_path: Path, parameters: Blueprint
):
    """Render JSON schema file from blocks"""
    nested_parameters = parameters.nested_parameters

    env = get_jinja2_environment()
    template = env.get_template("definitions.json.j2")
    rendered = template.render(
        {"nested_parameters": nested_parameters, "blueprint_name": blueprint_name}
    )
    # convert yaml to json
    rendered = yaml.safe_load(rendered)
    rendered = json.dumps(rendered, indent=2)
    json_schema_path.write_text(rendered)


def parse_annotations(
    blueprints_root_relative_path: str, policies_relative_path: str, jsonnet_data: str
) -> Blueprint:
    docblock_start_re = r".*\/\*\*$"
    docblock_end_re = r".*\*\/$"

    docblocks = []
    inside_docblock = False
    docblock_data = []
    for line in jsonnet_data.split("\n"):
        if re.match(docblock_start_re, line):
            assert not inside_docblock
            inside_docblock = True
        elif re.match(docblock_end_re, line):
            assert inside_docblock
            inside_docblock = False
            docblocks.append(
                Blueprint.from_comment(
                    blueprints_root_relative_path, policies_relative_path, docblock_data
                )
            )
            docblock_data = []
        else:
            if inside_docblock:
                docblock_data.append(line.strip())

    # merge docblocks
    merged_parameters = Blueprint()

    for block in docblocks:
        merge_parameternodes(
            merged_parameters.nested_parameters, block.nested_parameters
        )
        merge_parameternodes(
            merged_parameters.nested_required_parameters,
            block.nested_required_parameters,
        )

    return merged_parameters


def main(
    # blueprint_path is a path and is a required argument
    blueprint_path: Path = typer.Argument(
        ...,
        help="Path to the aperture blueprint directory",
        exists=True,
        dir_okay=True,
        file_okay=False,
        resolve_path=True,
    )
):
    repository_root = Path(__file__).absolute().parent.parent

    # calculate the path of repository_root/blueprints from the blueprint_path in terms of ../
    blueprints_root = repository_root / "blueprints"
    blueprint_name = blueprint_path.relative_to(blueprints_root)

    # get parts of relative_blueprint_path
    relative_blueprint_path_parts = blueprint_name.parts

    # make a prefix of ../ for each part
    policies_relative_path = "/".join([".."] * len(relative_blueprint_path_parts))
    docs_root_relative_path = "/".join(
        [".."] * (len(relative_blueprint_path_parts) + 2)
    )

    blueprints_root_relative_path = "/".join(
        [".."] * (len(relative_blueprint_path_parts) + 1)
    )

    blueprint_gen_path = blueprint_path / "gen"

    # create the gen directory if it doesn't exist
    blueprint_gen_path.mkdir(parents=True, exist_ok=True)

    config_parameters = parse_config_parameters(
        repository_root,
        blueprint_path,
        blueprints_root_relative_path,
        policies_relative_path,
    )
    render_json_schema(
        blueprint_name, blueprint_gen_path / "definitions.json", config_parameters
    )
    render_sample_config_yaml(
        blueprint_name, blueprint_gen_path / "values.yaml", False, config_parameters
    )
    render_sample_config_yaml(
        blueprint_name,
        blueprint_gen_path / "values-required.yaml",
        True,
        config_parameters,
    )

    dynamic_config_parameters = parse_dynamic_config_docblocks(
        repository_root,
        blueprint_path,
        blueprints_root_relative_path,
        policies_relative_path,
    )
    render_json_schema(
        blueprint_name,
        blueprint_gen_path / "dynamic-config-definitions.json",
        dynamic_config_parameters,
    )
    render_sample_config_yaml(
        blueprint_name,
        blueprint_gen_path / "dynamic-config-values.yaml",
        False,
        dynamic_config_parameters,
    )
    render_sample_config_yaml(
        blueprint_name,
        blueprint_gen_path / "dynamic-config-values-required.yaml",
        True,
        dynamic_config_parameters,
    )

    blueprints_docs_root_path = (
        repository_root / "docs/content/reference/policies/bundled-blueprints"
    )
    # check whether the blueprint_docs_root_path exists
    if blueprints_docs_root_path.exists():
        readme_path = (
            blueprints_docs_root_path
            / "/".join(relative_blueprint_path_parts[:-1])
            / f"{relative_blueprint_path_parts[-1]}.md"
        )

        update_docs_markdown(
            readme_path,
            config_parameters,
            dynamic_config_parameters,
            blueprint_name,
            docs_root_relative_path,
        )
    else:
        readme_path = blueprint_path / "README.md"
        update_readme_markdown(
            readme_path, config_parameters, dynamic_config_parameters
        )


def parse_config_parameters(
    repository_root: Path,
    blueprint_path: Path,
    blueprints_root_relative_path: str,
    policies_relative_path: str,
) -> Blueprint:
    config_path = blueprint_path / "config.libsonnet"

    if not config_path.exists():
        logger.error(f"config.libsonnet not found: {config_path}. Exiting.")
        raise typer.Exit(1)

    metadata_path = blueprint_path / "metadata.yaml"

    metadata = yaml.safe_load(metadata_path.read_text())

    parameters = parse_annotations(
        blueprints_root_relative_path, policies_relative_path, config_path.read_text()
    )

    # read deprecated property (string) from metadata
    parameters.deprecation_message = metadata.get("deprecated", None)

    # set defaults for nested parameters
    for source in metadata["sources"].keys():
        jsonnet_path = metadata["sources"][source]
        update_param_defaults(repository_root, config_path, parameters, jsonnet_path)

    return parameters


def parse_dynamic_config_docblocks(
    repository_root: Path,
    blueprint_path: Path,
    blueprints_root_relative_path: str,
    policies_relative_path: str,
) -> Blueprint:
    config_path = blueprint_path / "dynamic-config.libsonnet"
    if not config_path.exists():
        return Blueprint()

    dynamic_config_parameters = parse_annotations(
        blueprints_root_relative_path, policies_relative_path, config_path.read_text()
    )

    update_param_defaults(repository_root, config_path, dynamic_config_parameters)
    return dynamic_config_parameters


# merge 2nd ParameterNode into 1st


def merge_parameternodes(params1: ParameterNode, params2: ParameterNode):
    params1.required_children.update(params2.required_children)
    # recursive merge params1.children and params2.children
    for key, value in params2.children.items():
        if key in params1.children:
            merge_parameternodes(params1.children[key], value)
        else:
            params1.children[key] = value


def camel_to_kebab_case(s):
    # split on uppercase letters
    parts = re.split(r"(?=[A-Z])", s)
    # check if first part is empty
    if parts[0] == "":
        # remove first part
        parts = parts[1:]
    # lowercase all parts
    parts = [part.lower() for part in parts]
    # join with - all parts
    return "-".join(parts)


if __name__ == "__main__":
    typer.run(main)
