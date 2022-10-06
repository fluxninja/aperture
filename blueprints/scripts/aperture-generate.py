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


def parse_user_args(output: str, config: str) -> Tuple[Path, Path]:
    output_path = Path(output).absolute()
    config_path = Path(config).absolute()
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

    return output_path, config_path


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

    for policy_name, policy in policies["policies"].items():
        policy_path = output_path / "policies" / policy_name
        if not policy_path.parent.exists():
            policy_path.parent.mkdir()

        LOG.info(f"Creating {policy_path}")
        policy_path.write_text(yaml.dump(policy, default_flow_style=False))

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

    args = parser.parse_args()
    configure_logging(args.verbose)

    output_path, config_path = parse_user_args(
        args.output, args.config
    )

    script_path = Path(__file__).absolute()
    repo_root_path = script_path.parent.parent

    check_jsonnet()
    check_blueprint_configfile(repo_root_path, config_path)

    policies_json = generate_policies_json(repo_root_path, config_path)
    create_policies_yaml(output_path, policies_json)


if __name__ == "__main__":
    main()
