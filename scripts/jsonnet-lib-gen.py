from __future__ import annotations

import enum
import json
import os
import subprocess
import tempfile
from pathlib import Path
from typing import Dict, Iterable, List, Mapping, Optional, Tuple

import jinja2
import prance
import typer
import yaml
from loguru import logger

JINJA2_TEMPLATES = {}
JINJA2_TEMPLATES[
    "definition.libsonnet"
] = """
{{- definition | portsImports -}}

{
  new():: {
  },
{%- if definition.has_ports() %}
  {{ definition | portsBlock | indent(2) }}
{%- endif %}
{%- if definition.enum %}
  {%- for enum in definition.enum %}
  {{ enum }}: '{{ enum }}',
  {%- endfor %}
{%- elif definition.properties %}
  {%- for prop_name, prop in definition.properties.items() %}
  {{ prop_name | withName }}({{ prop_name | escapeJsonnetKeyword }}):: {
    {%- if prop.type_.name == "ARRAY" %}
    {{ prop_name | escapeJsonnetKeyword }}:
      if std.isArray({{ prop_name }})
      then {{ prop_name }}
      else [{{ prop_name }}],
    {%- else %}
    {{ prop_name | escapeJsonnetKeyword }}: {{ prop_name | escapeJsonnetKeyword }},
    {%- endif %}
  },
  {{ prop_name | withNameMixin }}({{ prop_name | escapeJsonnetKeyword }}):: {
    {{ prop_name | escapeJsonnetKeyword }}+: {{ prop_name | escapeJsonnetKeyword }},
  },
  {%- endfor %}
{%- endif %}
}

"""
JINJA2_TEMPLATES[
    "gen.libsonnet"
] = """{
  {%- for name, import in imports %}
  {{ name }}: import '{{ import }}',
  {%- endfor %}
}

"""

JINJA2_TEMPLATES[
    "in_out_ports_block.libsonnet"
] = """
{%- if in_ports %}
in_ports: {
  {%- for port in in_ports.definition.properties.keys() %}
  {{ port }}: error 'Port {{ port }} is missing',
  {%- endfor %}
},
{%- endif %}
{%- if out_ports %}
out_ports: {
  {%- for port in out_ports.definition.properties.keys() %}
  {{ port }}: error 'Port {{ port }} is missing',
  {%- endfor %}
},
{%- endif %}
"""

JSONNET_KEYWORDS = ["error"]


class JsonnetType(enum.Enum):
    BOOLEAN = enum.auto()
    NUMBER = enum.auto()
    STRING = enum.auto()
    ENUM = enum.auto()
    ARRAY = enum.auto()
    OBJECT = enum.auto()

    @classmethod
    def from_swagger(cls, swagger_obj: Mapping) -> JsonnetType:
        swagger_type = swagger_obj["type"]
        swagger_enum = swagger_obj.get("enum")
        match swagger_type:
            case "object":
                return JsonnetType.OBJECT
            case "array":
                return JsonnetType.ARRAY
            case "string":
                if swagger_enum:
                    return JsonnetType.ENUM
                else:
                    return JsonnetType.STRING
            case "boolean":
                return JsonnetType.BOOLEAN
            case "integer" | "number":
                return JsonnetType.NUMBER
            case _:
                raise AssertionError(f"Unsupported type: {swagger_type}")


class JsonnetObjectProperty:
    name: str
    deferred: bool
    definition_ref: Optional[str]
    definition: Optional[JsonnetDefinition]
    type_: JsonnetType

    def __init__(self, name: str, prop: Mapping):
        self.name = name
        self.definition = None
        self.definition_ref = prop.get("$ref")
        if not self.definition_ref:
            self.type_ = JsonnetType.from_swagger(prop)
            self.deferred = False
        else:
            self.deferred = True

    @classmethod
    def from_swagger(cls, name: str, prop: Mapping) -> JsonnetObjectProperty:
        return cls(name, prop)


class JsonnetDefinition:
    type_: JsonnetType
    swagger_name: str
    properties: Dict[str, JsonnetObjectProperty]
    enum: Optional[List[str]]

    def __init__(self, swagger_name: str, type_: JsonnetType):
        self.swagger_name = swagger_name
        self.type_ = type_
        self.properties = {}
        self.enum = []

    def has_ports(self):
        return "in_ports" in self.properties or "out_ports" in self.properties

    @property
    def jsonnet_name(self):
        return self.swagger_name

    def _parse_properties(self, properties: Mapping):
        for prop_name, prop in properties.items():
            logger.trace(f"{prop_name}")
            self.properties[prop_name] = JsonnetObjectProperty.from_swagger(
                prop_name, prop
            )

    @classmethod
    def from_swagger(cls, name: str, definition: Mapping):
        type_ = JsonnetType.from_swagger(definition)
        obj = cls(name, type_)
        match type_:
            case JsonnetType.OBJECT:
                properties = definition.get("properties")
                if properties:
                    obj._parse_properties(properties)
            case JsonnetType.ENUM:
                obj.enum = definition.get("enum")
        return obj

    @property
    def jsonnet_fname(self) -> str:
        return f"{self.jsonnet_name.lower()}.libsonnet"

    @property
    def jsonnet_path(self) -> Path:
        return Path("v1") / self.jsonnet_fname


def escapeJsonnetKeyword(word: str) -> str:
    if word in JSONNET_KEYWORDS:
        return word + "_"
    return word


def withNameFilter(name: str) -> str:
    parts = name.split("_")
    titleCase = "".join([part.capitalize() for part in parts])
    return f"with{titleCase}"


def withNameMixinFilter(name: str) -> str:
    with_name = withNameFilter(name)
    return f"{with_name}Mixin"


def defaultPorts(definition: JsonnetDefinition) -> str:
    if definition.type_ != JsonnetType.OBJECT:
        return ""
    in_ports = definition.properties.get("in_ports")
    out_ports = definition.properties.get("out_ports")
    if not in_ports and not out_ports:
        return ""
    env = get_jinja2_environment()
    template = env.get_template("in_out_ports_block.libsonnet")
    return template.render({"in_ports": in_ports, "out_ports": out_ports})


def portsImports(definition: JsonnetDefinition) -> str:
    imports = ""
    if definition.type_ != JsonnetType.OBJECT:
        return imports

    in_ports = definition.properties.get("in_ports")
    if in_ports:
        ports_name = f"{in_ports.definition.jsonnet_name.lower()}"
        ports_fname = f"{ports_name}.libsonnet"
        imports += f"local {ports_name} = import './{ports_fname}';\n"

    out_ports = definition.properties.get("out_ports")
    if out_ports:
        ports_name = f"{out_ports.definition.jsonnet_name.lower()}"
        ports_fname = f"{ports_name}.libsonnet"
        imports += f"local {ports_name} = import './{ports_fname}';\n"

    return imports


def portsBlock(definition: JsonnetDefinition) -> str:
    block = ""
    if definition.type_ != JsonnetType.OBJECT:
        return block

    in_ports = definition.properties.get("in_ports")
    if in_ports:
        ports_name = f"{in_ports.definition.jsonnet_name.lower()}"
        block += f"inPorts:: {ports_name},\n"

    out_ports = definition.properties.get("out_ports")
    if out_ports:
        ports_name = f"{out_ports.definition.jsonnet_name.lower()}"
        block += f"outPorts:: {ports_name},\n"

    return block.removesuffix("\n")


def get_jinja2_environment() -> jinja2.Environment:
    loader = jinja2.DictLoader(JINJA2_TEMPLATES)
    env = jinja2.Environment(loader=loader)
    env.filters["withName"] = withNameFilter
    env.filters["withNameMixin"] = withNameMixinFilter
    env.filters["escapeJsonnetKeyword"] = escapeJsonnetKeyword
    env.filters["defaultPorts"] = defaultPorts
    env.filters["portsImports"] = portsImports
    env.filters["portsBlock"] = portsBlock

    return env


PROTOBUF_IGNORED_DEFS = [
    "googlerpcStatus",
    "googleprotobufAny",
]


class ApertureJsonnetGenerator:
    swagger_path: Path
    definitions: Dict[str, JsonnetDefinition]

    def __init__(self, swagger_path: Path):
        self.swagger_path = swagger_path
        self.definitions = {}

    def _first_pass(self):
        RESOLVE_NONE = 0
        parser = prance.ResolvingParser(
            str(self.swagger_path), resolve_types=RESOLVE_NONE
        )
        assert parser.specification

        for swagger_name, swagger_def in parser.specification["definitions"].items():
            if swagger_name in PROTOBUF_IGNORED_DEFS:
                continue
            self.definitions[swagger_name] = JsonnetDefinition.from_swagger(
                swagger_name, swagger_def
            )

    def _second_pass(self):
        for definition in self.definitions.values():
            if definition.type_ != JsonnetType.OBJECT:
                continue

            # make a list of props to delete
            to_delete = []
            for prop in definition.properties.values():
                if not prop.deferred:
                    continue
                assert prop.definition_ref
                definition_ref_name = prop.definition_ref.split("/")[2]
                if definition_ref_name in PROTOBUF_IGNORED_DEFS:
                    # delete this prop
                    to_delete.append(prop.name)
                    continue
                definition_resolved_ref = self.definitions.get(definition_ref_name)
                if not definition_resolved_ref:
                    raise ValueError(f"Unknown Definition: {definition_ref_name}")
                prop.definition = definition_resolved_ref
                prop.type_ = JsonnetType.OBJECT
                prop.deferred = False
            for prop_name in to_delete:
                del definition.properties[prop_name]

    def parse(self):
        self._first_pass()
        self._second_pass()


def render_jsonnet_definition(definition: JsonnetDefinition) -> str:
    env = get_jinja2_environment()
    template = env.get_template("definition.libsonnet")
    return template.render(definition=definition)


def render_spec_libsonnet(custom_patches: Iterable[Path]) -> str:
    import_str = "(import 'gen.libsonnet')"
    # sort custom_patches by name
    custom_patches = sorted(custom_patches, key=lambda p: p.name)
    for patch_path in custom_patches:
        import_str += f" + (import '{patch_path}')"
    import_str += "\n"
    return import_str


def render_gen_libsonnet(path: Path, imports: List[Tuple[str, Path]]):
    env = get_jinja2_environment()
    template = env.get_template("gen.libsonnet")
    path.write_text(template.render(imports=imports))


def main(
    output_dir: Path = typer.Option(..., help="Output path for the generated library"),
    aperture_swagger_path: Path = typer.Argument(
        ..., help="Location of the aperture.swagger.yaml"
    ),
):
    if not aperture_swagger_path.exists():
        logger.error(f"No such file or directory: {aperture_swagger_path}")
        raise typer.Exit(1)

    generator = ApertureJsonnetGenerator(aperture_swagger_path)
    generator.parse()

    output_gen_dir = output_dir
    if output_gen_dir.exists():
        libsonnet_files = output_gen_dir.rglob("*.libsonnet")
        for path in libsonnet_files:
            if path.is_file():
                os.remove(path)

    definitions = generator.definitions.values()
    for definition in definitions:
        rendered_jsonnet = render_jsonnet_definition(definition)
        libsonnet_path = output_dir / definition.jsonnet_path
        if not libsonnet_path.parent.exists():
            libsonnet_path.parent.mkdir(parents=True)
        libsonnet_path.write_text(rendered_jsonnet)

    custom_patches_basedir = output_dir / "custom"
    custom_patches = [
        path.relative_to(output_dir)
        for path in custom_patches_basedir.glob("*.libsonnet")
    ]

    spec_libsonnet_path = output_dir / "spec.libsonnet"
    spec_libsonnet_data = render_spec_libsonnet(custom_patches)
    spec_libsonnet_path.write_text(spec_libsonnet_data)

    gen_libsonnet_path = output_dir / "gen.libsonnet"
    imports = [("v1", Path("v1/gen.libsonnet"))]
    render_gen_libsonnet(gen_libsonnet_path, imports)

    gen_libsonnet_path = output_dir / "v1" / "gen.libsonnet"
    imports = [
        (definition.jsonnet_name, Path(definition.jsonnet_fname))
        for definition in definitions
    ]
    render_gen_libsonnet(gen_libsonnet_path, imports)

    # rewrite exclusiveMinimum, exclusiveMaximum to use numbers based on minimum and maximum
    # instead of booleans in aperture_swagger_path, save to temp file
    with open(aperture_swagger_path, "r") as f:
        aperture_swagger_data = f.read()
        # read aperture_swagger_data as yaml
        aperture_swagger_yaml = yaml.safe_load(aperture_swagger_data)
        # iterate over all definitions
        for _, definition in aperture_swagger_yaml["definitions"].items():
            # iterate over all properties in definition
            for _, property in definition["properties"].items():
                # if the property has exclusiveMinimum or exclusiveMaximum, set it to the value of minimum or maximum
                if (
                    "exclusiveMinimum" in property
                    and property["exclusiveMinimum"] == True
                ):
                    property["exclusiveMinimum"] = property["minimum"]
                    # remove minimum
                    del property["minimum"]
                if (
                    "exclusiveMaximum" in property
                    and property["exclusiveMaximum"] == True
                ):
                    property["exclusiveMaximum"] = property["maximum"]
                    # remove maximum
                    del property["maximum"]
        # write aperture_swagger_yaml as yaml to temp file
        with tempfile.NamedTemporaryFile(mode="w", delete=False) as f:
            yaml.dump(aperture_swagger_yaml, f)

            # next, generate json schema using openapi2jsonschema
            # execute openapi2jsonschema with arguments -
            # 1. path in aperture_swagger_path variable
            # 2. --strict flag
            # 3. --output flag with value as output_dir/jsonschema
            jsonschema_dir = output_dir / "jsonschema"
            exit_code = subprocess.call(
                [
                    "openapi2jsonschema",
                    f.name,
                    "--strict",
                    "--output",
                    str(jsonschema_dir),
                ]
            )
            # remove temp file
            os.remove(f.name)
            if exit_code != 0:
                logger.error(
                    f"openapi2jsonschema exited with non-zero exit code: {exit_code}"
                )
                raise typer.Exit(1)
            # remove all files in output_dir/jsonschema except for the _definitions.json file
            for path in jsonschema_dir.rglob("*"):
                if path.is_file() and path.name != "_definitions.json":
                    os.remove(path)

            # inject k8s custom resource definition into _definitions.json
            jsonschema_definitions_path = jsonschema_dir / "_definitions.json"
            with open(jsonschema_definitions_path, "r") as f:
                jsonschema_definitions = json.load(f)
                jsonschema_definitions["definitions"][
                    "PolicyCustomResource"
                ] = json.loads(CUSTOM_RESOURCE_DEFINITION)
            with open(jsonschema_definitions_path, "w") as f:
                json.dump(jsonschema_definitions, f, indent=2)


CUSTOM_RESOURCE_DEFINITION = """
{
  "description": "CustomResourceDefinition represents a resource that should be exposed on the API server.  Its name MUST be in the format <.spec.name>.<.spec.group>.",
  "type": "object",
  "title": "Policy CustomResourceDefinition",
  "additionalProperties": false,
  "properties": {
    "apiVersion": {
      "description": "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
      "x-order": 0,
      "type": ["string", "null"],
      "enum": ["fluxninja.com/v1alpha1"]
    },
    "kind": {
      "description": "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
      "x-order": 1,
      "type": ["string", "null"],
      "enum": ["Policy"]
    },
    "metadata": {
      "x-order": 2,
      "$ref": "https://kubernetesjsonschema.dev/v1.18.1/_definitions.json#/definitions/io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta"
    },
    "spec": {
      "description": "Aperture Policy Object",
      "x-order": 3,
      "$ref": "#/definitions/Policy"
    },
    "dynamicConfig": {
        "description": "DynamicConfig provides dynamic configuration for the policy.",
      "x-order": 4,
        "type": ["object", "null"]
    },
    "status": {
      "description": "Status indicates the actual state of the CustomResourceDefinition",
      "x-order": 5,
      "$ref": "https://kubernetesjsonschema.dev/v1.18.1/_definitions.json#/definitions/io.k8s.apiextensions-apiserver.pkg.apis.apiextensions.v1beta1.CustomResourceDefinitionStatus"
    }
  }
}
"""


if __name__ == "__main__":
    typer.run(main)
