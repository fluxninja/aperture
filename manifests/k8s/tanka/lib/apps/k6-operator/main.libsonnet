local tanka = import 'github.com/grafana/jsonnet-libs/tanka-util/main.libsonnet';
local kustomize = tanka.kustomize.new(std.thisFile);

{
  //crd: kustomize.build(path='k6/crd'),
  //manager: kustomize.build(path='k6/manager')
  default: kustomize.build(path='k6/default'),
}
