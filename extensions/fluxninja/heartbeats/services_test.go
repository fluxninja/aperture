package heartbeats

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	entitiesv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/discovery/entities/v1"
	heartbeatv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/fluxninja/v1"
	"github.com/fluxninja/aperture/v2/pkg/discovery/entities"
)

var _ = Describe("Services", func() {
	var ec *entities.Entities

	BeforeEach(func() {
		ec = entities.NewEntities()
	})

	Context("Services", func() {
		It("reads same service from two entities", func() {
			ec.PutForTest(testEntity("1", "1.1.1.1", "some_name", []string{"baz"}))
			ec.PutForTest(testEntity("2", "1.1.1.2", "some_name", []string{"baz"}))
			sl := populateServicesList(ec)
			Expect(sl.Services).To(HaveLen(1))
			Expect(sl.Services).To(ContainElement(&heartbeatv1.Service{
				Name:          "baz",
				EntitiesCount: 2,
			}))
		})

		It("reads two services from one entity", func() {
			ip := "1.1.1.1"
			serviceNames := []string{"baz1", "baz2"}
			name := "entity_1234"
			ec.PutForTest(testEntity("1", ip, name, serviceNames))
			sl := populateServicesList(ec)
			Expect(sl.Services).To(HaveLen(2))
			Expect(sl.Services).To(ContainElement(&heartbeatv1.Service{
				Name:          "baz1",
				EntitiesCount: 1,
			}))
			Expect(sl.Services).To(ContainElement(&heartbeatv1.Service{
				Name:          "baz2",
				EntitiesCount: 1,
			}))
		})

		It("returns no service after being cleared", func() {
			ip := "1.1.1.1"
			serviceNames := []string{"baz"}
			name := "entity_1234"
			ec.PutForTest(testEntity("1", ip, name, serviceNames))
			ec.Clear()
			sl := populateServicesList(ec)
			Expect(sl.Services).To(HaveLen(0))
		})
	})
})

func testEntity(uid, ipAddress, name string, services []string) *entitiesv1.Entity {
	return &entitiesv1.Entity{
		Uid:       uid,
		IpAddress: ipAddress,
		Name:      name,
		Services:  services,
	}
}
