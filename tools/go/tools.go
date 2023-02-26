//go:build tools
// +build tools

package tools

import (
	_ "github.com/dmarkham/enumer"
	_ "github.com/favadi/protoc-go-inject-tag"
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
	_ "github.com/golang/mock/mockgen"
	_ "gotest.tools/gotestsum"
	_ "sigs.k8s.io/controller-runtime/tools/setup-envtest"
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen"
)
