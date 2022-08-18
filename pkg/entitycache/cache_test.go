package entitycache_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	heartbeatv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/plugins/fluxninja/v1"
	"github.com/fluxninja/aperture/pkg/entitycache"
)

const (
	testPrefix = "test"
)

var _ = Describe("Cache", func() {
	var ec *entitycache.EntityCache

	BeforeEach(func() {
		ec = entitycache.NewEntityCache()
	})

	Context("by IP", func() {
		It("reads entity properly", func() {
			ip := "1.2.3.4"
			entity := testEntity("foo", "foo", ip, nil)
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
			entity := testEntity("foo", "foo", ip, nil)
			ec.Put(entity)

			removed := ec.Remove(entity)
			Expect(removed).To(BeTrue())

			found := ec.GetByIP(ip)
			Expect(found).To(BeNil())
		})

		It("returns false if trying to remove a nonexistent entity", func() {
			ip := "1.2.3.4"
			otherIP := "192.168.0.1"
			entity := testEntity("foo", "foo", ip, nil)
			ec.Put(entity)

			otherEntity := testEntity("foo2", "foo", otherIP, nil)
			removed := ec.Remove(otherEntity)
			Expect(removed).To(BeFalse())

			found := ec.GetByIP(ip)
			Expect(found).To(Equal(entity))
		})
	})

	Context("by Name", func() {
		It("reads entity properly", func() {
			uid := "foo"
			name := nameFromUid(testPrefix, uid)
			entity := testEntity(uid, "foo", "", nil)
			ec.Put(entity)
			actual := ec.GetByName(name)
			Expect(actual).To(Equal(entity))
		})

		It("returns nil when no entity found", func() {
			uid := "foo"
			name := nameFromUid(testPrefix, uid)
			actual := ec.GetByName(name)
			Expect(actual).To(BeNil())
		})

		It("removes an entity properly", func() {
			uid := "bar"
			name := nameFromUid(testPrefix, uid)
			entity := testEntity(uid, "foo", "", nil)
			ec.Put(entity)

			removed := ec.Remove(entity)
			Expect(removed).To(BeTrue())

			found := ec.GetByName(name)
			Expect(found).To(BeNil())
		})

		It("returns false if trying to remove a nonexistent entity", func() {
			uid := "bar"
			name := nameFromUid(testPrefix, uid)
			otherUid := "baz"
			entity := testEntity(uid, "foo", "1.1.1.1", nil)
			ec.Put(entity)

			otherEntity := testEntity(otherUid, "foo", "1.1.1.2", nil)
			removed := ec.Remove(otherEntity)
			Expect(removed).To(BeFalse())

			found := ec.GetByName(name)
			Expect(found).To(Equal(entity))
		})
	})

	It("clears all entities from the map", func() {
		ip := "1.2.3.4"
		entity := testEntity("foo", "foo", "", nil)
		ec.Put(entity)
		ec.Clear()
		found := ec.GetByIP(ip)
		Expect(found).To(BeNil())
	})

	Context("Services", func() {
		It("reads same service from two entities", func() {
			ec.Put(testEntity("1", "foo", "1.1.1.1", []string{"baz"}))
			ec.Put(testEntity("2", "foo", "1.1.1.2", []string{"baz"}))
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
			ec.Put(testEntity("1", "foo", ip, serviceNames))
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
			ec.Put(testEntity("1", "foo", ip, serviceNames))
			ec.Clear()
			services, _ := ec.Services()
			Expect(services).To(HaveLen(0))
		})
	})
})

func testEntity(uid, agentGroup, ipAddress string, services []string) *entitycache.Entity {
	entity := entitycache.NewEntity(entitycache.EntityID{
		Prefix: "test",
		UID:    uid,
	}, ipAddress, services)
	entity.SetAgentGroup(agentGroup)
	return entity
}

func nameFromUid(prefix, uid string) string {
	return fmt.Sprintf("%v-%v", prefix, uid)
}
