# Classification extractors

(formerly known as "YAML classification policies")

Companion package for aperture.tech/aperture/api/gen/proto/go/aperture/classification/v1

Label extractors (shortly – extractors) are a high-level way to specify how to
extract a flow label given http request metadata, without a need to write rego
code.

Eg.:

```yaml
extractor:
  json:
    from: request.http.body
    pointer: /query
```

Internally, multiple extractors will be compiled to a single rego module.
