package autoscale

import (
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
)

// ParsePodAutoScaler parses a PodAutoScaler and returns its nested circuit representation.
func ParsePodAutoScaler(
	podautoscaler *policylangv1.PodAutoScaler,
) (*policylangv1.NestedCircuit, error) {
	autoscaler := &policylangv1.AutoScaler{
		Scaler: &policylangv1.AutoScaler_Scaler{
			Scaler: &policylangv1.AutoScaler_Scaler_KubernetesReplicas{
				KubernetesReplicas: podautoscaler.PodScaler,
			},
		},
		MinScale:                   podautoscaler.MinReplicas,
		MaxScale:                   podautoscaler.MaxReplicas,
		ScaleOutControllers:        podautoscaler.ScaleOutControllers,
		ScaleInControllers:         podautoscaler.ScaleInControllers,
		MaxScaleOutPercentage:      podautoscaler.MaxScaleOutPercentage,
		MaxScaleInPercentage:       podautoscaler.MaxScaleInPercentage,
		ScaleOutCooldown:           podautoscaler.ScaleOutCooldown,
		ScaleInCooldown:            podautoscaler.ScaleInCooldown,
		CooldownOverridePercentage: podautoscaler.CooldownOverridePercentage,
		ScaleOutAlerterParameters:  podautoscaler.ScaleOutAlerterParameters,
		ScaleInAlerterParameters:   podautoscaler.ScaleInAlerterParameters,
	}

	if podautoscaler.OutPorts != nil {
		autoscaler.OutPorts = &policylangv1.AutoScaler_Outs{
			ActualScale:     podautoscaler.OutPorts.ActualReplicas,
			ConfiguredScale: podautoscaler.OutPorts.ConfiguredReplicas,
			DesiredScale:    podautoscaler.OutPorts.DesiredReplicas,
		}
	}

	autoscalerCircuit, err := ParseAutoScaler(autoscaler)
	autoscalerCircuit.Name = "PodAutoScaler"

	return autoscalerCircuit, err
}
