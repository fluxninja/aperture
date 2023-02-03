#!/usr/bin/env python3
#
# This script is used to generate generate blueprint config.libsonnet based on
# dashboards and policies that this policy include, as well as README.md with
# configuration documentation section.
from __future__ import annotations

import dataclasses
import sys
import subprocess
import tempfile
import typer
import yaml
import json
import re


from pathlib import Path
from typing import Dict, List, Optional, Any

from loguru import logger
import jinja2


class ExitException(RuntimeError):
    ret_code: int

    def __init__(self, ret_code=0):
        self.ret_code = ret_code


@dataclasses.dataclass
class DocBlockParam:
    param_name: str
    param_type: str
    param_link: str
    description: str
    required: bool
    default: Optional[Any] = None


SECTION_RE = re.compile(r".*@section (?P<section>.+)")
SUBSECTION_RE = re.compile(r".*@subsection (?P<subsection>.+)")
PARAMETER_RE = re.compile(r".*@param.*")
PARAMETER_DETAILED_RE = re.compile(r".*@param \((?P<param_name>[\w.\[\]]+): (?P<param_type>[\w.\[\]]+) ?(?P<param_required>\w+)?\) (?P<param_description>.+)")


@dataclasses.dataclass
class DocBlock:
    section: str
    subsection: Optional[str]
    description: Optional[str]
    parameters: Dict[str, DocBlockParam]

    @classmethod
    def _resolve_param_to_policy_ref(cls, param: str) -> str:
        if not param.startswith("aperture.spec") and not param.startswith("[]aperture.spec"):
            return ""

        component = param.split(".")[-1]

        # Transform CamelCase into camel-case
        parts = [[component[0].lower()]]
        for letter in component[1:]:
            if letter.isupper():
                parts.append(list(letter.lower()))
            else:
                parts[-1].append(letter)

        component_final = "v1-" + "-".join(["".join(l) for l in parts])

        return f"../../spec#{component_final}"

    @classmethod
    def from_comment(cls, comment: List[str]) -> DocBlock:
        section = None
        subsection = None
        description = ""
        description_parsed = False
        parameters = {}
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
                param_link = cls._resolve_param_to_policy_ref(param_type)
                param_required = groups.get("param_required", "") == "required"
                param_description = groups["param_description"]
                parameters[param_name] = (DocBlockParam(param_name, param_type, param_link, param_description, param_required))
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

        return cls(section, subsection, description, parameters)


def command_with_exit_code(func):
    def wrapper():
        try:
            func()
        except ExitException as ex:
            sys.exit(ex.ret_code)

    return wrapper


def _get_params_for_blocks(blocks: List[DocBlock], required: bool) -> List[DocBlockParam]:
    return [param for block in blocks for param in block.parameters.values() if param.required == required]


def _generate_required_configuration(blocks: List[DocBlock]) -> str:
    """Creates a temporary config objects and renders given policy/dashboard to extract
       default values from resulting _config object."""

    required_params = _get_params_for_blocks(blocks, required=True)

    # Create an object with all required parameters for the given policy/dashboard
    root = {}
    def create_nested_object_with_value(param_name):
        parent = root
        parts = param_name.split(".")[1:]
        for idx, part in enumerate(parts):
            if idx == len(parts) - 1:
                parent[part] = "FAKE-VALUE"
            else:
                if part not in parent:
                    parent[part] = {}
                parent = parent[part]

    for param in required_params:
        create_nested_object_with_value(param.param_name)

    jsonnet_config = "{"
    def append_required_params(obj):
        result = "{"
        for key, value in obj.items():
            if isinstance(value, str):
                result += f"'{key}': '{value}',"
            else:
                inner = append_required_params(value)
                result += f"{key}+: {inner},"
        result += "}"

        return result

    jsonnet_config = append_required_params(root)

    return jsonnet_config


def update_docblock_param_defaults(repository_root: Path, jsonnet_path: Path, blocks: List[DocBlock]):
    jsonnet_data = f"local fn = import '{jsonnet_path}';\n"
    required_config = _generate_required_configuration(blocks)
    jsonnet_data += f"local cfg = {required_config};\n"
    jsonnet_data += f"fn(cfg)\n"
    jsonnet_data += "+ { _config::: super._config, policy:: super.policy, circuit:: super.circuit }\n"

    rendered_config = None
    with tempfile.NamedTemporaryFile(suffix=".libsonnet") as tmp:
        tmppath = Path(tmp.name)
        tmppath.write_text(jsonnet_data)

        jsonnet_jpaths = [
            "-J", repository_root,
            "-J", repository_root / "vendor",
        ]

        try:
            result = subprocess.run(["jsonnet", *jsonnet_jpaths, str(tmppath)], capture_output=True, check=True)
        except subprocess.CalledProcessError as ex:
            logger.error(f"Error while rendering jsonnet: {ex.stderr}")
            # log file for debugging
            logger.error(f"Jsonnet file: {jsonnet_data}")
            raise typer.Exit(1)

        rendered_config = json.loads(result.stdout)

    def get_param_default_from_rendered_config(root: Dict, name: str) -> Any:
        parts = name.split(".")[1:]
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
    params_with_defaults = _get_params_for_blocks(blocks, required=False)
    for param in params_with_defaults:
        param.default = get_param_default_from_rendered_config(rendered_config["_config"], param.param_name)
        logger.trace(param)


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
    reference="{{ param.param_link }}"
    value='{% if param.value %}{{ param.default | quoteValue }}{% endif %}'
    description='{{ param.description }}' />

{%- endfor %}
{%- endfor %}

{%- endfor %}
"""

def quoteValue(value: str) -> str:
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

    return f"\"{value}\""


def get_jinja2_environment() -> jinja2.Environment:
    JINJA2_TEMPLATES = {"section.md.j2": SECTION_TPL}
    loader = jinja2.DictLoader(JINJA2_TEMPLATES)
    env = jinja2.Environment(loader=loader)
    env.filters["quoteValue"] = quoteValue


    return env


def update_docblock_param_names(blocks: List[DocBlock], prefix: str):
    """When rendering README.md, parameter names should be prefixed with blueprint prefix"""
    for block in blocks:
        for param in block.parameters.values():
            param.param_name = f"{prefix}.{param.param_name}"


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
    <td><code>{value != '' ? value : "REQUIRED VALUE"}</code></td>
  </tr>
  <tr>
    <td class="blueprints-description"><ParameterHeading>Description</ParameterHeading></td>
    <td class="blueprints-description"><WrappedDescription>{description}</WrappedDescription></td>
  </tr>
</table>
);
```
"""


def update_readme_markdown(readme_path, blocks: List[DocBlock]):
    """Find configuration marker in README.md and append all blocks after it"""

    readme_data = readme_path.read_text()
    readme_copied = ""
    for line in readme_data.split("\n"):
        if line == "<!-- Configuration Marker -->":
            readme_copied += line + "\n"
            break
        readme_copied += line + "\n"

    readme_copied += f"\n{MDX_TEMPLATE}\n"

    sections = {}
    for block in blocks:
        if block.section not in sections:
            sections[block.section] = []
        sections[block.section].append(block)

    env = get_jinja2_environment()
    template = env.get_template("section.md.j2")
    rendered = template.render({"sections": sections})
    readme_copied += rendered
    readme_path.write_text(readme_copied)


def extract_docblock_comments(jsonnet_data: str) -> List[DocBlock]:
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
            docblocks.append(DocBlock.from_comment(docblock_data))
            docblock_data = []
        else:
           if inside_docblock:
               docblock_data.append(line.strip())
    return docblocks


def main(blueprint_path: Path = typer.Argument(..., help="Path to the aperture blueprint directory")):
    repository_root = Path(__file__).absolute().parent.parent

    if not blueprint_path.exists():
        logger.error(f"No such file or directory: {blueprint_path}")
        raise typer.Exit(1)

    readme_path = blueprint_path / "README.md"
    if not readme_path.exists():
        logger.error(f"README.md not found: {readme_path}. Exiting.")
        raise typer.Exit(1)

    config_path = blueprint_path / "config.libsonnet"

    metadata_path = blueprint_path / "metadata.yaml"

    metadata = yaml.safe_load(metadata_path.read_text())

    docblocks = extract_docblock_comments(config_path.read_text())

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
        path = metadata["sources"][section]["path"]
        # append common section to blocks
        blocks.extend(sections["Common"])
        update_docblock_param_defaults(repository_root, path, blocks)

    update_readme_markdown(readme_path, docblocks)

if __name__ == "__main__":
    typer.run(main)
