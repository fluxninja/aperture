# Contributing

We would really appreciate your help!

See our [Code of Conduct](CODE_OF_CONDUCT.md).

### Creating Pull Requests

When you are ready to contribute, pick an issue you'd like to solve. Try
starting with:

- [Good First Issue](https://github.com/fluxninja/aperture/issues?utf8=%E2%9C%93&q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)

For your convenience, all the needed development tools can be installed with
`asdf` by running `make install-asdf-tools`.

Before committing, install `pre-commit` hooks, which will automatically check if
your code meets our standards.

```
pre-commit install --hook-type={pre-commit,commit-msg,prepare-commit-msg}
pre-commit install-hooks
```

After your first PR is created you would be asked to sign our
[Contributor License Agreement](https://cla-assistant.io/fluxninja/aperture).
