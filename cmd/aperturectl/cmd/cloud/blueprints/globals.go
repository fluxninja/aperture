package blueprints

import (
	cloudutils "github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/cloud/utils"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
)

var (
	controller cloudutils.ControllerConn
	client     utils.CloudBlueprintsClient

	valuesFile string
)
