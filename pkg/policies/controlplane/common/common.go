package common

import (
	"aperture.tech/aperture/pkg/config"
)

// FxOptionsFuncTag allows sub-modules to provide their options to per policy apps independently.
var FxOptionsFuncTag = config.GroupTag("policy-fx-funcs")
