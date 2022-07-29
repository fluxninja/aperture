package common

import (
	"github.com/FluxNinja/aperture/pkg/config"
)

// FxOptionsFuncTag allows sub-modules to provide their options to per policy apps independently.
var FxOptionsFuncTag = config.GroupTag("policy-fx-funcs")
