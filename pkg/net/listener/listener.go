// +kubebuilder:validation:Optional
package listener

import (
	"context"
	"net"

	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// ListenerConfig holds configuration for socket listeners.
// swagger:model
// +kubebuilder:object:generate=true
type ListenerConfig struct {
	// Keep-alive period - 0 = enabled if supported by protocol or operating system. If negative, then keep-alive is disabled.
	KeepAlive config.Duration `json:"keep_alive" validate:"gte=0s" default:"180s"`

	// Address to bind to in the form of `[host%zone]:port`
	Addr string `json:"addr" validate:"hostname_port" default:":8080"`

	// TCP networks - `tcp`, `tcp4` (IPv4-only), `tcp6` (IPv6-only)
	Network string `json:"network,omitempty" validate:"omitempty,oneof=tcp tcp4 tcp6,omitempty" default:"tcp"`
}

func newListener(config ListenerConfig) (net.Listener, error) {
	listenConfig := net.ListenConfig{KeepAlive: config.KeepAlive.AsDuration()}

	listener, err := listenConfig.Listen(context.Background(), config.Network, config.Addr)
	if err != nil {
		log.Error().Err(err).Str("addr", config.Addr).Msg("Unable to announce on local network address")
		return nil, err
	}

	return listener, nil
}
