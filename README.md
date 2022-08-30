## Aperture Helm Chart
---

```
helm repo add aperture https://fluxninja.github.io/aperture/
helm install aperture --namespace aperture-system aperture/aperture --create-namespace
```
To install operator use this command
```
helm install aperture-operator --namespace aperture-system aperture/aperture-operator --create-namespace
```