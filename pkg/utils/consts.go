package utils

const (
	// ApertureAgent is the service name for Agent.
	ApertureAgent = "aperture-agent"
	// ApertureController is the service name for Controller.
	ApertureController = "aperture-controller"

	// ApertureCloudAgentGroup agent group name for Cloud Agents.
	ApertureCloudAgentGroup = "aperture-cloud"

	// InstallationModeKubernetesSidecar for sidecar installation mode.
	InstallationModeKubernetesSidecar = "KUBERNETES_SIDECAR"
	// InstallationModeKubernetesDaemonSet for Kubernetes DaemonSet installation mode.
	InstallationModeKubernetesDaemonSet = "KUBERNETES_DAEMONSET"
	// InstallationModeLinuxBareMetal for bare metal installation mode.
	InstallationModeLinuxBareMetal = "LINUX_BARE_METAL"
	// InstallationModeCloudAgent for Cloud Agents.
	InstallationModeCloudAgent = "CLOUD_AGENT"
)
