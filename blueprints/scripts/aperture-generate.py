#!/usr/bin/env python3
from __future__ import annotations

import argparse
import logging
import subprocess
import sys
import tempfile
import json
import os
import yaml

from pathlib import Path
import textwrap
from typing import Dict, List, Tuple


LOG = logging.getLogger(__file__)


CLI_DESCRIPTION = """
Aperture policies & dashboards generator utility.

This utility can be used to generate policies and dashboards "in-place", that
is generate json files with grafana dashboards and aperture policies, as opposed
to using blueprints library as jsonnet library. Check main README.md for more
details.
"""


class ExitException(RuntimeError):
    ret_code: int

    def __init__(self, ret_code=0):
        self.ret_code = ret_code


def command_with_exit_code(func):
    def wrapper():
        try:
            func()
        except ExitException as ex:
            sys.exit(ex.ret_code)

    return wrapper


def parse_user_args(blueprint: str, output: str, config: str) -> Tuple[Path, Path, Path]:
    blueprint_path = Path(blueprint).absolute()
    output_path = Path(output).absolute()
    config_path = Path(config).absolute()
    if not blueprint_path.exists():
        LOG.error(f"Policy directory does not exist: {blueprint_path}")
        raise ExitException(1)

    blueprint_main_path = blueprint_path / "main.libsonnet"
    if not blueprint_main_path.exists():
        LOG.error(f"Policy main.libsonnet is missing: {blueprint_main_path}")
        raise ExitException(1)

    if not config_path.exists():
        LOG.error(f"Policy jsonnet configuration does not exist: {config_path}")
        raise ExitException(1)

    if output_path.exists():
        if not output_path.is_dir():
            LOG.error("Path out the output directory exists, but is not a directory")
            raise ExitException(1)
        if not os.access(output_path, os.W_OK):
            LOG.error("Output directory not writeable")
            raise ExitException(1)
    else:
        output_parent = output_path.parent
        if not output_parent.exists():
            LOG.error("Output directory does not exist and its parent is also missing")
            raise ExitException(1)
        if not os.access(output_parent, os.W_OK):
            LOG.error("Output directory does not exist and its parent is not writeable")
            raise ExitException(1)
        output_path.mkdir()

    return blueprint_path, output_path, config_path


def configure_logging(verbose: bool):
    console_handler = logging.StreamHandler()
    LOG.addHandler(console_handler)
    if verbose:
        LOG.setLevel(logging.DEBUG)
        log_format = "%(asctime)s | %(levelname)s: %(message)s"
    else:
        LOG.setLevel(logging.INFO)
        log_format = "%(message)s"
    console_handler.setFormatter(logging.Formatter(log_format))


def get_jsonnet_library_search_paths(repository_root: Path) -> List[str]:
    search_paths = [
        "--jpath",
        repository_root / "blueprints",
        "--jpath",
        repository_root / "vendor",
    ]
    return search_paths


def check_jsonnet():
    try:
        subprocess.run(["jsonnet", "--version"], check=True, capture_output=True)
    except FileNotFoundError as ex:
        LOG.error("Couldn't find jsonnet binary in $PATH")
        raise ExitException(1)
    except subprocess.CalledProcessError as ex:
        LOG.error(f"Couldn't run jsonnet binary: {ex.stderr.decode('utf-8').strip()}")
        raise ExitException(1)


def check_blueprint_configfile(repository_root: Path, config_path: Path):
    """Parse config with jsonnet and do some very basic validation of the resulting JSON object"""

    jsonnet_search_paths = get_jsonnet_library_search_paths(repository_root)
    try:
        subprocess_cmd = ["jsonnet", *jsonnet_search_paths, str(config_path)]
        LOG.debug("Running: `" + " ".join([str(i) for i in subprocess_cmd]) + "`")
        subprocess.run(
            subprocess_cmd,
            check=True,
            capture_output=True,
        )
    except subprocess.CalledProcessError as ex:
        jsonnet_stderr = textwrap.indent(ex.stderr.decode("utf-8").strip(), prefix="  ")
        LOG.error(f"Couldn't parse blueprint configuration as jsonnet:\n{jsonnet_stderr}")
        raise ExitException(1)


def generate_blueprint_entrypoint(blueprint_relative_path: Path, config_path: Path) -> str:
    config_text = config_path.read_text()
    blueprint_main_path = str(blueprint_relative_path / "main.libsonnet")
    data = f"local blueprint = import '{blueprint_main_path}';\n"
    data += f"local config = ({config_text});\n"
    data += "blueprint + { _config:: config }\n"
    return data


def generate_policies_json(repository_root: Path, config_path: Path) -> Dict:
    jsonnet_search_paths = get_jsonnet_library_search_paths(repository_root)
    try:
        subprocess_cmd = ["jsonnet", *jsonnet_search_paths, str(config_path)]
        LOG.debug("Running: `" + " ".join([str(i) for i in subprocess_cmd]) + "`")
        result = subprocess.run(
            subprocess_cmd,
            check=True,
            capture_output=True,
        )
    except subprocess.CalledProcessError as ex:
        jsonnet_stderr = textwrap.indent(ex.stderr.decode("utf-8").strip(), prefix="  ")
        LOG.error(f"Couldn't parse blueprint configuration as jsonnet:\n{jsonnet_stderr}")
        raise ExitException(1)

    return json.loads(result.stdout)


def create_policies_yaml(output_path: Path, policies: Dict):
    for dashboard_name, dashboard in policies["dashboards"].items():
        dashboard_path = output_path / "dashboards" / dashboard_name
        if not dashboard_path.parent.exists():
            dashboard_path.parent.mkdir()

        LOG.info(f"Creating {dashboard_path}")
        dashboard_path.write_text(json.dumps(dashboard["dashboard"], indent=4))

    for blueprint_name, blueprint in policies["policies"].items():
        blueprint_path = output_path / "policies" / blueprint_name
        if not blueprint_path.parent.exists():
            blueprint_path.parent.mkdir()

        LOG.info(f"Creating {blueprint_path}")
        blueprint_path.write_text(yaml.dump(blueprint["policy"], default_flow_style=False))
        #blueprint_path.write_text(json.dumps(blueprint, indent=4))


@command_with_exit_code
def main():
    parser = argparse.ArgumentParser(
        description=CLI_DESCRIPTION,
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument(
        "--verbose",
        default=False,
        action="store_true",
        help="Whether to log verbose messages to stderr",
    )
    parser.add_argument(
        "--output", default="_gen", help="Output directory for json files"
    )
    parser.add_argument(
        "--config", help="jsonnet file with blueprint configuration"
    )
    parser.add_argument(
        "blueprint", metavar="BLUEPRINT", type=str, help="Aperture blueprint path"
    )

    args = parser.parse_args()
    configure_logging(args.verbose)

    blueprint_path, output_path, config_path = parse_user_args(
        args.blueprint, args.output, args.config
    )

    script_path = Path(__file__).absolute()
    repo_root_path = script_path.parent.parent
    policies_root_path = repo_root_path / "blueprints"
    blueprint_relative_path = blueprint_path.resolve().relative_to(policies_root_path.resolve())

    check_jsonnet()
    check_blueprint_configfile(repo_root_path, config_path)

    with tempfile.TemporaryDirectory() as tmp:
        tempdir = Path(tmp)
        entrypoint_path = tempdir / "main.jsonnet"
        entrypoint_text = generate_blueprint_entrypoint(blueprint_relative_path, config_path)
        LOG.debug(entrypoint_text)
        entrypoint_path.write_text(entrypoint_text)

        policies_json = generate_policies_json(repo_root_path, entrypoint_path)

        create_policies_yaml(output_path, policies_json)


if __name__ == "__main__":
    main()
