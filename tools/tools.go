//go:build tools
// +build tools

package tools

import (
	_ "github.com/dmarkham/enumer"
	_ "github.com/favadi/protoc-go-inject-tag"
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/google/pprof"
	_ "golang.org/x/tools/cmd/godoc"
	_ "golang.org/x/tools/cmd/goimports"
	_ "gotest.tools/gotestsum"
)
