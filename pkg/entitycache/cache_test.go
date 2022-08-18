package entitycache_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	heartbeatv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/plugins/fluxninja/v1"
	"github.com/fluxninja/aperture/pkg/entitycache"
)

var _ = Describe("Cache", func() {
	var ec *entitycache.EntityCache

	BeforeEach(func() {
		ec = entitycache.NewEntityCache()
	})

	Context("by IP", func() {
		It("reads entity properly", func() {
			ip := "1.2.3.4"
			name := "entity_1234"
			entity := testEntity("foo", "foo", ip, name, nil)
			ec.Put(entity)
			actual := ec.GetByIP(ip)
			Expect(actual).To(Equal(entity))
		})

		It("returns nil when no entity found", func() {
			ip := "1.2.3.4"
			actual := ec.GetByIP(ip)
			Expect(actual).To(BeNil())
		})

		It("removes an entity properly", func() {
			ip := "1.2.3.4"
			name := "entity_1234"
			entity := testEntity("foo", "foo", ip, name, nil)
			ec.Put(entity)

			removed := ec.Remove(entity)
			Expect(removed).To(BeTrue())

			found := ec.GetByIP(ip)
			Expect(found).To(BeNil())
		})

		It("returns false if trying to remove a nonexistent entity", func() {
			ip := "1.2.3.4"
			otherIP := "192.168.0.1"
			name := "entity_1234"
			otherName := "other_entity_4321"
			entity := testEntity("foo", "foo", ip, name, nil)
			ec.Put(entity)

			otherEntity := testEntity("foo2", "foo", otherIP, otherName, nil)
			removed := ec.Remove(otherEntity)
			Expect(removed).To(BeFalse())

			found := ec.GetByIP(ip)
			Expect(found).To(Equal(entity))
		})
	})

	Context("by Name", func() {
		It("reads entity properly", func() {
			uid := "foo"
			name := "some_name"
			entity := testEntity(uid, "foo", "", name, nil)
			ec.Put(entity)
			actual := ec.GetByName(name)
			Expect(actual).To(Equal(entity))
		})

		It("returns nil when no entity found", func() {
			name := "some_name"
			actual := ec.GetByName(name)
			Expect(actual).To(BeNil())
		})

		It("removes an entity properly", func() {
			uid := "bar"
			name := "some_name"
			entity := testEntity(uid, "foo", "", name, nil)
			ec.Put(entity)

			removed := ec.Remove(entity)
			Expect(removed).To(BeTrue())

			found := ec.GetByName(name)
			Expect(found).To(BeNil())
		})

		It("returns false if trying to remove a nonexistent entity", func() {
			uid := "bar"
			name := "some_name"
			otherUid := "baz"
			otherName := "another_name"
			entity := testEntity(uid, "foo", "1.1.1.1", name, nil)
			ec.Put(entity)

			otherEntity := testEntity(otherUid, "foo", "1.1.1.2", otherName, nil)
			removed := ec.Remove(otherEntity)
			Expect(removed).To(BeFalse())

			found := ec.GetByName(name)
			Expect(found).To(Equal(entity))
		})
	})

	It("clears all entities from the map", func() {
		ip := "1.2.3.4"
		entity := testEntity("foo", "foo", "", "some_name", nil)
		ec.Put(entity)
		ec.Clear()
		found := ec.GetByIP(ip)
		Expect(found).To(BeNil())
	})

	Context("Services", func() {
		It("reads same service from two entities", func() {
			ec.Put(testEntity("1", "foo", "1.1.1.1", "some_name", []string{"baz"}))
			ec.Put(testEntity("2", "foo", "1.1.1.2", "some_name", []string{"baz"}))
			services, _ := ec.Services()
			Expect(services).To(HaveLen(1))
			Expect(services).To(ContainElement(&heartbeatv1.Service{
				Name:          "baz",
				EntitiesCount: 2,
			}))
		})

		It("reads two services from one entity", func() {
			ip := "1.1.1.1"
			serviceNames := []string{"baz1", "baz2"}
			name := "entity_1234"
			ec.Put(testEntity("1", "foo", ip, name, serviceNames))
			services, _ := ec.Services()
			Expect(services).To(HaveLen(2))
			Expect(services).To(ContainElement(&heartbeatv1.Service{
				Name:          "baz1",
				EntitiesCount: 1,
			}))
			Expect(services).To(ContainElement(&heartbeatv1.Service{
				Name:          "baz2",
				EntitiesCount: 1,
			}))
		})

		It("returns no service after being cleared", func() {
			ip := "1.1.1.1"
			serviceNames := []string{"baz"}
			name := "entity_1234"
			ec.Put(testEntity("1", "foo", ip, name, serviceNames))
			ec.Clear()
			services, _ := ec.Services()
			Expect(services).To(HaveLen(0))
		})
	})
})

func testEntity(uid, agentGroup, ipAddress, name string, services []string) *entitycache.Entity {
	entity := entitycache.NewEntity(entitycache.EntityID{
		Prefix: "test",
		UID:    uid,
	}, ipAddress, name, services)
	entity.SetAgentGroup(agentGroup)
	return entity
}
