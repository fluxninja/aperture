package plugins

import (
	// This is needed to pin kube-openapi version in go mod as direct dependency.
	// Without it, loading plugin fails with `plugin was built with a different
	// version of package k8s.io/kube-openapi`.
	_ "k8s.io/kube-openapi/pkg/util"
)
