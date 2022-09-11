<p align="center">
  <img src="docs/content/assets/img/aperture_logo.png" alt="Fluxninja Aperture" width="75%">
  <br/>

  <a href="https://dl.circleci.com/status-badge/img/gh/fluxninja/aperture/tree/main.svg?style=svg&circle-token=cf4312657fbc2f4833fee89328a3f27ab5f39c10">
    <img alt="Build Status" src="https://img.shields.io/circleci/build/github/fluxninja/aperture/main?token=cf4312657fbc2f4833fee89328a3f27ab5f39c10&style=for-the-badge">
  </a>
  <a href="https://goreportcard.com/report/github.com/fluxninja/aperture">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/fluxninja/aperture?style=for-the-badge">
  </a>
  <a href="https://codecov.io/gh/fluxninja/aperture/branch/main/">
    <img alt="Codecov Status" src="https://img.shields.io/codecov/c/github/fluxninja/aperture?style=for-the-badge">
  </a>
  <a href="https://pkg.go.dev/github.com/fluxninja/aperture">
    <img alt="Godoc Reference" src="https://img.shields.io/badge/godoc-reference-brightgreen?style=for-the-badge">
  </a>
</p>

## What is Aperture?

Aperture is the first open-source flow control and reliability management
platform for modern web applications.

Aperture enables flow control through observing, analyzing, and actuating,
facilitated by agents and a controller.

<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/fluxninja/aperture/main/docs/content/assets/img/oaadark.png" />
    <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/fluxninja/aperture/main/docs/content/assets/img/oaalight.png" />
    <img alt="Observe. Analyze. Actuate." src="https://raw.githubusercontent.com/fluxninja/aperture/main/docs/content/assets/img/oaalight.png">
  </picture>
</p>

For more detailed description, refer to [the docs](https://docs.fluxninja.com/).

## Playground

To try aperture in local K8s environment, refer to
[Playground docs](playground/README.md).

## Contributing

We would really appreciate your help!

[![Slack](https://img.shields.io/badge/Join%20Our%20Community-Slack-brightgreen)](https://join.slack.com/t/aperturetech/shared_invite/zt-1ewkfjfy9-~wF4EryoDyJ6kaPRTNZPyA)

See our [Code of Conduct](CODE_OF_CONDUCT.md).

### Reporting bugs or requesting features

Reporting bugs helps us improve Aperture to be more reliable and user friendly.
Please make sure to include all the required information to reproduce and
understand the bug you are reporting.

Follow helper questions in bug report template to make it easier.

If you see a way to improve Aperture, use the feature request template to create
an issue. Make sure to explain the problem you are trying to solve and what is
the expected behavior.

### Creating Pull Requests

When you are ready to contribute, pick an issue you'd like to solve. Try
starting with:

- [Good First Issue](https://github.com/fluxninja/aperture/issues?utf8=%E2%9C%93&q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)

For your convenience, all the needed development tools can be installed with
[`asdf`](playground/README.md#tools).

Before committing, install `pre-commit` hooks, which will automatically check if
your code meets our standards.

```
pre-commit install --hook-type={pre-commit,commit-msg,prepare-commit-msg}
pre-commit install-hooks
```

After your first PR is created you would be asked to sign our
[Contributor License Agreement](https://cla-assistant.io/fluxninja/aperture).

## Resources

For better understanding of Aperture, refer to the following resources.

- [Documentation](https://docs.fluxninja.com/)
