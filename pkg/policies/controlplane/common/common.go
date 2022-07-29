package common

import (
	"github.com/fluxninja/aperture/pkg/config"
)

// FxOptionsFuncTag allows sub-modules to provide their options to per policy apps independently.
var FxOptionsFuncTag = config.GroupTag("policy-fx-funcs")
