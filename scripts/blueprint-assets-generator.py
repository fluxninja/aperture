#!/usr/bin/env python3
#
# This script is used to generate generate blueprint config.libsonnet based on
# dashboards and policies that this policy include, as well as README with
# configuration documentation section.
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


class ExitException(RuntimeError):
    ret_code: int

    def __init__(self, ret_code=0):
        self.ret_code = ret_code


@dataclasses.dataclass
class DocBlockParam:
    param_name: str
    param_type: str
    json_schema_link: str
    spec_link: str
    description: str
    required: bool
    default: Optional[Any] = None


SECTION_RE = re.compile(r".*@section (?P<section>.+)")
SUBSECTION_RE = re.compile(r".*@subsection (?P<subsection>.+)")
PARAMETER_RE = re.compile(r".*@param.*")
PARAMETER_DETAILED_RE = re.compile(
    r".*@param \((?P<param_name>[\w.\[\]]+): (?P<param_type>[\w.\[\]]+) ?(?P<param_required>\w+)?\) (?P<param_description>.+)"
)

# DocBlockNode for nested_parameters


@dataclasses.dataclass
class DocBlockNode:
    parameter: DocBlockParam
    children: Dict[str, DocBlockNode]
    required_children: OrderedSet[str]


@dataclasses.dataclass
class DocBlock:
    section: str
    subsection: Optional[str]
    description: Optional[str]
    # flat view of parameters
    parameters: Dict[str, DocBlockParam]
    # nested dictionary of parameters
    nested_parameters: DocBlockNode
    # nested dictionary of required parameters
    nested_required_parameters: DocBlockNode

    @classmethod
    def _resolve_param_links(
        cls, aperture_json_schema_path: str, spec_path: str, param: str
    ) -> tuple[str, str]:
        if not param.startswith("aperture.spec") and not param.startswith(
            "[]aperture.spec"
        ):
            return "", ""

        component = param.split(".")[-1]

        # Transform CamelCase into camel-case
        parts = [[component[0].lower()]]
        for letter in component[1:]:
            if letter.isupper():
                parts.append(list(letter.lower()))
            else:
                parts[-1].append(letter)

        component_final = "-".join(["".join(l) for l in parts])

        return (
            f"{spec_path}#{component_final}",
            f"{aperture_json_schema_path}{component}",
        )

    @classmethod
    def from_comment(
        cls, aperture_json_schema_path: str, spec_path: str, comment: List[str]
    ) -> DocBlock:
        section = None
        subsection = None
        description = ""
        description_parsed = False
        parameters = {}
        nested_parameters = DocBlockNode(
            DocBlockParam("", "intermediate_node", "", "", "", False),
            {},
            OrderedSet([]),
        )
        nested_required_parameters = DocBlockNode(
            DocBlockParam("", "intermediate_node", "", "", "", False),
            {},
            OrderedSet([]),
        )
        for line in comment:
            if matched := SECTION_RE.match(line.strip()):
                if description:
                    description_parsed = True
                section = matched.group(1)
            elif matched := SUBSECTION_RE.match(line.strip()):
                if description:
                    description_parsed = True
                subsection = matched.group(1)
            elif matched := PARAMETER_RE.match(line.strip()):
                if description:
                    description_parsed = True

                inner = PARAMETER_DETAILED_RE.match(line.strip())
                if not inner:
                    logger.error(f"Unable to parse @param: `{line.strip()}`")
                    raise typer.Exit(1)

                groups = inner.groupdict()
                param_name, param_type = groups["param_name"], groups["param_type"]
                spec_link, json_schema_link = cls._resolve_param_links(
                    aperture_json_schema_path, spec_path, param_type
                )
                param_required = groups.get("param_required", "") == "required"
                param_description = groups["param_description"]
                parameters[param_name] = DocBlockParam(
                    param_name,
                    param_type,
                    json_schema_link,
                    spec_link,
                    param_description,
                    param_required,
                )
                # tokenize param_name and create nested_parameters
                parts = param_name.split(".")
                parent = nested_parameters.children
                parent_required = nested_required_parameters.children
                if param_required:
                    nested_parameters.required_children.add(parts[0])
                    nested_required_parameters.required_children.add(parts[0])
                for idx, part in enumerate(parts):
                    if idx == len(parts) - 1:
                        node = DocBlockNode(
                            DocBlockParam(
                                part,
                                param_type,
                                json_schema_link,
                                spec_link,
                                param_description,
                                param_required,
                            ),
                            {},
                            OrderedSet([]),
                        )
                        parent[part] = node
                        if param_required:
                            parent_required[part] = node
                    else:
                        next_part = parts[idx + 1]
                        node = DocBlockNode(
                            DocBlockParam(part, "intermediate_node", "", "", "", False),
                            {},
                            OrderedSet([]),
                        )
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

            else:
                stripped = line.lstrip(" ")
                stripped = stripped.removeprefix("* ")
                stripped = stripped.removeprefix("*")
                if description:
                    assert not description_parsed
                    description += stripped + "\n"
                elif not description and stripped:
                    description = stripped + "\n"

        if not section or not parameters:
            logger.error("Unable to find section and parameters in docblock")
            raise ValueError()

        return cls(
            section,
            subsection,
            description,
            parameters,
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


def update_docblock_param_defaults(
    repository_root: Path,
    config_path: Path,
    blocks: List[DocBlock],
    jsonnet_path: Path = Path(),
    config_key: str = "",
):
    jsonnet_data = f"local config = import '{config_path}';\n"
    if jsonnet_path != Path():
        jsonnet_data += f"local fn = import '{jsonnet_path}';\n"

    if jsonnet_path != Path():
        jsonnet_data += f"local cfg = config.common + config.{config_key};\n"
        jsonnet_data += f"fn(cfg)\n"
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
                return config[part]
            else:
                try:
                    config = config[part]
                except KeyError:
                    # When specific param is a map (map[string]type) and there is no default
                    # then we return None here, which will be converted into an empty map later.
                    return None

    logger.trace(rendered_config)
    # process defaults
    for block in blocks:
        for param in block.parameters.values():
            param.default = get_param_default_from_rendered_config(
                rendered_config["_config"], param.param_name
            )
            logger.trace(param)

    # walk nested_parameters and update defaults
    def update_nested_param_defaults(node, prefix):
        if node.parameter.param_type != "intermediate_node":
            node.parameter.default = get_param_default_from_rendered_config(
                rendered_config["_config"], prefix
            )
        for key, child in node.children.items():
            if prefix != "":
                keyPrefix = f"{prefix}.{key}"
            else:
                keyPrefix = key
            update_nested_param_defaults(child, keyPrefix)

    for block in blocks:
        update_nested_param_defaults(block.nested_parameters, "")
    for block in blocks:
        update_nested_param_defaults(block.nested_required_parameters, "")


def update_docblock_sections(blocks: List[DocBlock], section: str):
    """Updates all blocks to use blueprint-level section as section, and move block section to subsection"""
    for block in blocks:
        block.subsection = block.section
        block.section = section


SECTION_TPL = """
{% for section, blocks in sections.items() %}

<h3 class="blueprints-h3">{{ section }}</h3>
  {%- for block in blocks %}

{%- if block.subsection %}

<h4 class="blueprints-h4">{{ block.subsection }}</h4>
{%- endif %}
{%- if block.description %}

{{ block.description }}
{%- endif %}

{%- for param in block.parameters.values() %}

<ParameterDescription
    name="{{ param.param_name }}"
    type="{{ param.param_type }}"
    reference="{{ param.spec_link }}"
    value="{{ param.default | quoteValueDocs }}"
    description='{{ param.description }}' />

{%- endfor %}
{%- endfor %}

{%- endfor %}
"""

JSON_SCHEMA_TPL = """
{%- macro render_type(param_type, json_schema_link) %}
{%- if param_type.startswith('aperture.spec') %}
"type": "object",
{{ ' ' * 4 }}"$ref": "{{- json_schema_link }}"
{%- elif param_type == 'bool' %}
"type": "boolean"
{%- elif param_type == 'float32' %}
"type": "number",
"format": "float"
{%- elif param_type == 'float64' %}
"type": "number",
"format": "double"
{%- elif param_type == 'int32' %}
"type": "integer",
"format": "int32"
{%- elif param_type == 'int64' %}
"type": "integer",
"format": "int64"
{%- elif param_type.startswith('[]') %}
"type": "array",
"items": { {{- render_type(param_type[2:], json_schema_link) }} }
{%- elif param_type.startswith('map[') %}
"type": "object",
"additionalProperties": { {{- render_type(param_type[4:-1], json_schema_link) }}}
{%- else %}
"type": "{{ param_type }}"
{%- endif %}
{%- endmacro %}

{%- macro render_properties(node) %}
"{{ node.parameter.param_name }}": {
{%- if node.parameter.param_type == 'intermediate_node' %}
"type": "object",
"additionalProperties": false,
{%- if node.required_children %}
"required": [{%- for child_name in node.required_children %}"{{ child_name }}"{%- if not loop.last %},{% endif %}{%- endfor %}],
{%- endif %}
"properties": {
{%- for child_name, child_node in node.children.items() %}
{{ render_properties(child_node) }}
{%- if not loop.last %},{% endif %}
{%- endfor %}
}
{%- else %}
"description": "{{ node.parameter.description }}",
"default": {{ node.parameter.default | quoteValueJSON }},
{{- render_type(node.parameter.param_type, node.parameter.json_schema_link) }}
{%- endif %}
}
{%- endmacro %}

{
"$schema": "http://json-schema.org/draft-07/schema#",
"type": "object",
"title": "{{ blueprint_name }} blueprint",
"additionalProperties": false,
{%- if json_schema_data.required_children %}
"required": [{%- for child_name in json_schema_data.required_children %}"{{ child_name }}"{%- if not loop.last %},{% endif %}{%- endfor %}],
{%- endif %}
"properties": {
{%- for child_name, child_node in json_schema_data.children.items() %}
{{ render_properties(child_node) }}
{%- if not loop.last %},{% endif %}
{%- endfor %}
}
}
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
{%- else %}
{{- value | quoteValueYAML }}
{%- endif %}
{%- endmacro %}
{%- macro render_node(node, level) %}
{%- if node.parameter.param_type != 'intermediate_node' %}
{{ '  ' * level }}# {{ node.parameter.description }}
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
{%- endmacro %}
{%- for child_name, child_node in sample_config_data.children.items() %}
{{ render_node(child_node, 0) }}
{%- endfor %}
"""


def quoteValueYAML(value: str) -> str:
    # if value is __REQUIRED_FIELD__ return as unquoted string
    if value == "__REQUIRED_FIELD__":
        return value
    return quoteValueJSON(value)


def quoteValueJSON(value: str) -> str:
    return json.dumps(value)


def quoteValueDocs(value: str) -> str:
    # if value is __REQUIRED_FIELD__ return as unquoted string
    if value == "__REQUIRED_FIELD__":
        return value

    if isinstance(value, bool):
        return str(value).lower()

    try:
        int(value)
        return value
    except (ValueError, TypeError):
        pass

    try:
        float(value)
        return value
    except (ValueError, TypeError):
        pass

    if isinstance(value, list) or isinstance(value, dict):
        return value

    return f"'{value}'"


def get_jinja2_environment() -> jinja2.Environment:
    JINJA2_TEMPLATES = {
        "section.md.j2": SECTION_TPL,
        "values.yaml.j2": YAML_TPL,
        "definitions.json.j2": JSON_SCHEMA_TPL,
    }
    loader = jinja2.DictLoader(JINJA2_TEMPLATES)
    env = jinja2.Environment(loader=loader)
    env.filters["quoteValueYAML"] = quoteValueYAML
    env.filters["quoteValueJSON"] = quoteValueJSON
    env.filters["quoteValueDocs"] = quoteValueDocs
    return env


MDX_TEMPLATE = """
```mdx-code-block

export const ParameterHeading = ({children}) => (
  <span style={{fontWeight: "bold"}}>{children}</span>
);

export const WrappedDescription = ({children}) => (
  <span style={{wordWrap: "normal"}}>{children}</span>
);

export const RefType = ({type, reference}) => (
  <a href={reference}>{type}</a>
);

export const ParameterDescription = ({name, type, reference, value, description}) => (
  <table class="blueprints-params">
  <tr>
    <td><ParameterHeading>Parameter</ParameterHeading></td>
    <td><code>{name}</code></td>
  </tr>
  <tr>
    <td><ParameterHeading>Type</ParameterHeading></td>
    <td><em>{reference == "" ? type : <RefType type={type} reference={reference} />}</em></td>
  </tr>
  <tr>
    <td class="blueprints-default-heading"><ParameterHeading>Default Value</ParameterHeading></td>
    <td><code>{value}</code></td>
  </tr>
  <tr>
    <td class="blueprints-description"><ParameterHeading>Description</ParameterHeading></td>
    <td class="blueprints-description"><WrappedDescription>{description}</WrappedDescription></td>
  </tr>
</table>
);
```
"""


def update_readme_markdown(
    readme_path: Path,
    config_blocks: List[DocBlock],
    dynamic_config_blocks: List[DocBlock],
    blueprint_name: Path,
    aperture_version_path,
):
    """Find configuration marker in README and add generated content below it."""

    readme_data = readme_path.read_text()
    readme_copied = ""
    for line in readme_data.split("\n"):
        readme_copied += line + "\n"
        if line == "<!-- Configuration Marker -->":
            break

    readme_copied += f"\n{MDX_TEMPLATE}\n"

    readme_copied += "```mdx-code-block\n"
    readme_copied += (
        f"import {{apertureVersion as aver}} from '{aperture_version_path}'\n"
    )
    readme_copied += "```\n\n"

    readme_copied += f"Code: <a href={{`https://github.com/fluxninja/aperture/tree/${{aver}}/blueprints/{blueprint_name}`}}>{blueprint_name}</a>\n"
    sections = {}
    for block in config_blocks:
        if block.section not in sections:
            sections[block.section] = []
        sections[block.section].append(block)

    env = get_jinja2_environment()
    template = env.get_template("section.md.j2")

    rendered = template.render({"sections": sections})
    readme_copied += rendered

    sections = {}
    for block in dynamic_config_blocks:
        if block.section not in sections:
            sections[block.section] = []
        sections[block.section].append(block)
    if len(sections) > 0:
        readme_copied += "\n\n:::note\n\n"
        readme_copied += "The following configuration parameters can be [dynamically configured](/reference/aperturectl/apply/dynamic-config/dynamic-config.md) at runtime, without reloading the policy.\n\n"
        readme_copied += ":::\n\n"
        rendered = template.render({"sections": sections})
        readme_copied += rendered

    readme_path.write_text(readme_copied)


def render_sample_config_yaml(
    blueprint_name: Path,
    sample_config_path: Path,
    only_required: bool,
    blocks: List[DocBlock],
):
    """Render sample config YAML file from blocks"""
    # merge all nested parameters into one dict
    sample_config_data = DocBlockNode(
        DocBlockParam("", "intermediate_node", "", "", "", False), {}, OrderedSet([])
    )
    if only_required is False:
        for block in blocks:
            merge_docblocknodes(sample_config_data, block.nested_parameters)
    else:
        for block in blocks:
            merge_docblocknodes(sample_config_data, block.nested_required_parameters)

    env = get_jinja2_environment()
    template = env.get_template("values.yaml.j2")
    rendered = template.render(
        {"sample_config_data": sample_config_data, "blueprint_name": blueprint_name}
    )
    sample_config_path.write_text(rendered)


def render_json_schema(
    blueprint_name: Path, json_schema_path: Path, blocks: List[DocBlock]
):
    """Render JSON schema file from blocks"""
    json_schema_data = DocBlockNode(
        DocBlockParam("", "intermediate_node", "", "", "", False), {}, OrderedSet([])
    )
    for block in blocks:
        merge_docblocknodes(json_schema_data, block.nested_parameters)

    env = get_jinja2_environment()
    template = env.get_template("definitions.json.j2")
    rendered = template.render(
        {"json_schema_data": json_schema_data, "blueprint_name": blueprint_name}
    )
    json_schema_path.write_text(rendered)


def extract_docblock_comments(
    aperture_json_schema_path: str, spec_path: str, jsonnet_data: str
) -> List[DocBlock]:
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
                DocBlock.from_comment(
                    aperture_json_schema_path, spec_path, docblock_data
                )
            )
            docblock_data = []
        else:
            if inside_docblock:
                docblock_data.append(line.strip())
    return docblocks


def main(
    blueprint_path: Path = typer.Argument(
        ..., help="Path to the aperture blueprint directory"
    )
):
    if not blueprint_path.exists():
        logger.error(f"No such file or directory: {blueprint_path}")
        raise typer.Exit(1)

    repository_root = Path(__file__).absolute().parent.parent

    # calculate the path of repository_root/blueprints from the blueprint_path in terms of ../
    blueprints_root = repository_root / "blueprints"
    blueprint_name = blueprint_path.relative_to(blueprints_root)

    # get parts of relative_blueprint_path
    relative_blueprint_path_parts = blueprint_name.parts

    readme_path = (
        repository_root
        / "docs/content/reference/policies/bundled-blueprints"
        / "/".join(relative_blueprint_path_parts[:-1])
        / f"{relative_blueprint_path_parts[-1]}.md"
    )

    if not readme_path.exists():
        logger.error(f"README not found: {readme_path}. Exiting.")
        raise typer.Exit(1)

    # make a prefix of ../ for each part
    spec_path = "/".join([".."] * len(relative_blueprint_path_parts)) + "/spec"
    aperture_version_path = (
        "/".join([".."] * (len(relative_blueprint_path_parts) + 2))
        + "/apertureVersion.js"
    )

    aperture_json_schema_path = (
        "/".join([".."] * (len(relative_blueprint_path_parts) + 1))
        + "/gen/jsonschema/_definitions.json#/definitions/"
    )

    blueprint_gen_path = blueprint_path / "gen"

    config_docblocks = parse_config_docblocks(
        repository_root, blueprint_path, aperture_json_schema_path, spec_path
    )
    render_json_schema(
        blueprint_name, blueprint_gen_path / "definitions.json", config_docblocks
    )
    render_sample_config_yaml(
        blueprint_name, blueprint_gen_path / "values.yaml", False, config_docblocks
    )
    render_sample_config_yaml(
        blueprint_name,
        blueprint_gen_path / "values-required.yaml",
        True,
        config_docblocks,
    )

    dynamic_config_docblocks = parse_dynamic_config_docblocks(
        repository_root, blueprint_path, aperture_json_schema_path, spec_path
    )
    render_json_schema(
        blueprint_name,
        blueprint_gen_path / "dynamic-config-definitions.json",
        dynamic_config_docblocks,
    )
    render_sample_config_yaml(
        blueprint_name,
        blueprint_gen_path / "dynamic-config-values.yaml",
        False,
        dynamic_config_docblocks,
    )
    render_sample_config_yaml(
        blueprint_name,
        blueprint_gen_path / "dynamic-config-values-required.yaml",
        True,
        dynamic_config_docblocks,
    )

    update_readme_markdown(
        readme_path,
        config_docblocks,
        dynamic_config_docblocks,
        blueprint_name,
        aperture_version_path,
    )


def parse_config_docblocks(
    repository_root: Path,
    blueprint_path: Path,
    aperture_json_schema_path: str,
    spec_path: str,
) -> List[DocBlock]:
    config_path = blueprint_path / "config.libsonnet"

    if not config_path.exists():
        logger.error(f"config.libsonnet not found: {config_path}. Exiting.")
        raise typer.Exit(1)

    metadata_path = blueprint_path / "metadata.yaml"

    metadata = yaml.safe_load(metadata_path.read_text())

    docblocks = extract_docblock_comments(
        aperture_json_schema_path, spec_path, config_path.read_text()
    )

    sections = {section: [] for section in metadata["sources"].keys()}
    # append Common to sections
    sections["Common"] = []

    for block in docblocks:
        if block.section not in sections and block.section != "Common":
            logger.error(f"Unknown docblocks @section: {block.section}")
            raise typer.Exit(1)
        sections[block.section].append(block)

    # Make sure that all non-required parameters have their default values updated based on library defaults
    for section, blocks in sections.items():
        # Skip common section
        if section == "Common":
            continue
        jsonnet_path = metadata["sources"][section]["path"]
        config_key = metadata["sources"][section]["config_key"]
        # append common section to blocks
        blocks.extend(sections["Common"])
        update_docblock_param_defaults(
            repository_root, config_path, blocks, jsonnet_path, config_key
        )
    return docblocks


def parse_dynamic_config_docblocks(
    repository_root: Path,
    blueprint_path: Path,
    aperture_json_schema_path: str,
    spec_path: str,
) -> List[DocBlock]:
    config_path = blueprint_path / "dynamic-config.libsonnet"
    if not config_path.exists():
        return []

    docblocks = extract_docblock_comments(
        aperture_json_schema_path, spec_path, config_path.read_text()
    )

    sections = {}
    for block in docblocks:
        sections.setdefault(block.section, []).append(block)

    for _, blocks in sections.items():
        update_docblock_param_defaults(repository_root, config_path, blocks)
    return docblocks


# merge 2nd docblock into 1st


def merge_docblocknodes(docblock1: DocBlockNode, docblock2: DocBlockNode):
    docblock1.required_children.update(docblock2.required_children)
    # recursive merge docblock1.children and docblock2.children
    for key, value in docblock2.children.items():
        if key in docblock1.children:
            merge_docblocknodes(docblock1.children[key], value)
        else:
            docblock1.children[key] = value


if __name__ == "__main__":
    typer.run(main)
