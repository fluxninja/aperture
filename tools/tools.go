//go:build tools
// +build tools

package tools

import (
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/bufbuild/buf/cmd/protoc-gen-buf-lint"
	_ "github.com/dmarkham/enumer"
	_ "github.com/favadi/protoc-go-inject-tag"
	_ "github.com/fullstorydev/grpcurl/cmd/grpcurl"
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/google/pprof"
	_ "github.com/mikefarah/yq/v4"
	_ "github.com/vektra/mockery/v2"
	_ "golang.org/x/tools/cmd/godoc"
	_ "golang.org/x/tools/cmd/stringer"
)
