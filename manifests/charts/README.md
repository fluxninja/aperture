# Helm Charts

## Common Errors

### Error: found in Chart.yaml, but missing in charts/ directory

You didn't install chart dependencies defined in Chart.yaml.

```sh
helm dependency build ./agent
```

It might fail with `no repository definition for <URL>`,
in which case you should run:

```sh
help repo add <URL>
```

and then try again.
