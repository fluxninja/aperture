package heartbeats

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	entitycachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/entitycache/v1"
	heartbeatv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/plugins/fluxninja/v1"
	"github.com/fluxninja/aperture/pkg/entitycache"
)

var _ = Describe("Services", func() {
	var ec *entitycache.EntityCache

	BeforeEach(func() {
		ec = entitycache.NewEntityCache()
	})

	Context("Services", func() {
		It("reads same service from two entities", func() {
			ec.Put(testEntity("1", "1.1.1.1", "some_name", []string{"baz"}))
			ec.Put(testEntity("2", "1.1.1.2", "some_name", []string{"baz"}))
			entityCache := populateServicesList(ec)
			Expect(entityCache.Services).To(HaveLen(1))
			Expect(entityCache.Services).To(ContainElement(&heartbeatv1.Service{
				Name:          "baz",
				EntitiesCount: 2,
			}))
		})

		It("reads two services from one entity", func() {
			ip := "1.1.1.1"
			serviceNames := []string{"baz1", "baz2"}
			name := "entity_1234"
			ec.Put(testEntity("1", ip, name, serviceNames))
			entityCache := populateServicesList(ec)
			Expect(entityCache.Services).To(HaveLen(2))
			Expect(entityCache.Services).To(ContainElement(&heartbeatv1.Service{
				Name:          "baz1",
				EntitiesCount: 1,
			}))
			Expect(entityCache.Services).To(ContainElement(&heartbeatv1.Service{
				Name:          "baz2",
				EntitiesCount: 1,
			}))
		})

		It("returns no service after being cleared", func() {
			ip := "1.1.1.1"
			serviceNames := []string{"baz"}
			name := "entity_1234"
			ec.Put(testEntity("1", ip, name, serviceNames))
			ec.Clear()
			entityCache := populateServicesList(ec)
			Expect(entityCache.Services).To(HaveLen(0))
		})
	})
})

func testEntity(uid, ipAddress, name string, services []string) *entitycachev1.Entity {
	return &entitycachev1.Entity{
		Prefix:    "test",
		Uid:       uid,
		IpAddress: ipAddress,
		Name:      name,
		Services:  services,
	}
}
