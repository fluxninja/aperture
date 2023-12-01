package servicegetter_test

import (
	"context"
	"net"
	"testing"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/peer"

	entitiesv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/discovery/entities/v1"
	discoveryentities "github.com/fluxninja/aperture/v2/pkg/discovery/entities"
	servicegetter "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service-getter"
)

func TestServiceGetter(t *testing.T) {
	entities := discoveryentities.NewEntities()
	sg := servicegetter.FromEntities(entities)

	t.Run("ServicesFromContext with no peer information", func(t *testing.T) {
		ctx := context.Background()
		services := sg.ServicesFromContext(ctx)
		assert.Empty(t, services)
	})

	t.Run("ServicesFromContext with invalid IP address", func(t *testing.T) {
		ctx := peerContext("invalid")
		services := sg.ServicesFromContext(ctx)
		assert.Empty(t, services)
	})

	t.Run("ServicesFromContext with valid IP address but no entity", func(t *testing.T) {
		ip := "192.168.1.1"
		ctx := peerContext(ip)
		services := sg.ServicesFromContext(ctx)
		assert.Empty(t, services)
	})

	t.Run("ServicesFromContext with valid IP address and entity", func(t *testing.T) {
		ip := "192.168.1.2"
		entity := discoveryentities.NewEntity(&entitiesv1.Entity{
			IpAddress: ip,
			Services:  []string{"svc1", "svc2"},
		})
		entities.Put(entity)
		defer entities.Remove(entity)

		ctx := peerContext(ip)
		services := sg.ServicesFromContext(ctx)
		assert.Equal(t, entity.Services(), services)
	})

	t.Run("ServicesFromSocketAddress with invalid IP address", func(t *testing.T) {
		addr := &corev3.SocketAddress{Address: "invalid"}
		services := sg.ServicesFromSocketAddress(addr)
		assert.Equal(t, []string(nil), services)
	})

	t.Run("ServicesFromSocketAddress with valid IP address but no entity", func(t *testing.T) {
		ip := "192.168.1.3"
		addr := &corev3.SocketAddress{Address: ip}
		services := sg.ServicesFromSocketAddress(addr)
		assert.Equal(t, []string(nil), services)
	})

	t.Run("ServicesFromSocketAddress with valid IP address and entity", func(t *testing.T) {
		ip := "192.168.1.4"
		entity := discoveryentities.NewEntity(&entitiesv1.Entity{
			IpAddress: ip,
			Services:  []string{"svc3", "svc4"},
		})
		entities.Put(entity)
		defer entities.Remove(entity)

		addr := &corev3.SocketAddress{Address: ip}
		services := sg.ServicesFromSocketAddress(addr)
		assert.Equal(t, entity.Services(), services)
	})

	t.Run("ServicesFromAddress with invalid IP address", func(t *testing.T) {
		services := sg.ServicesFromAddress("invalid")
		assert.Nil(t, services)
	})

	t.Run("ServicesFromAddress with valid IP address but no entity", func(t *testing.T) {
		ip := "192.168.1.5"
		services := sg.ServicesFromAddress(ip)
		assert.Nil(t, services)
	})

	t.Run("ServicesFromAddress with valid IP address and entity", func(t *testing.T) {
		ip := "192.168.1.6"
		entity := discoveryentities.NewEntity(&entitiesv1.Entity{
			IpAddress: ip,
			Services:  []string{"svc5", "svc6"},
		})
		entities.Put(entity)
		defer entities.Remove(entity)

		services := sg.ServicesFromAddress(ip)
		assert.Equal(t, entity.Services(), services)
	})
}

func peerContext(ip string) context.Context {
	return peer.NewContext(context.Background(), &peer.Peer{
		Addr: &net.TCPAddr{
			IP: net.ParseIP(ip),
		},
	})
}
