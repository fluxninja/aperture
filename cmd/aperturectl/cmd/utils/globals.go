package utils

import (
	"path/filepath"
)

var (
	// AllowDeprecated is a flag to allow deprecated blueprints to be listed.
	AllowDeprecated = false

	// AperturectlRootDir is the root directory for aperturectl.
	AperturectlRootDir = ".aperturectl"
	// BlueprintsCacheRoot is the root directory for blueprints cache.
	BlueprintsCacheRoot = filepath.Join(AperturectlRootDir, "blueprints")
	// BuilderCacheRoot is the root directory for builder cache.
	BuilderCacheRoot = filepath.Join(AperturectlRootDir, "build")
)
