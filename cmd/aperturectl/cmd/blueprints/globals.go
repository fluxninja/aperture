package blueprints

import policyv1alpha1 "github.com/fluxninja/aperture/operator/api/policy/v1alpha1"

var (
	blueprintName string
	outputDir     string
	valuesFile    string
	graphDir      string
	applyPolicy   bool
	kubeConfig    string
	validPolicies []*policyv1alpha1.Policy
	all           bool
	onlyRequired  bool
	skipPull      bool
)
