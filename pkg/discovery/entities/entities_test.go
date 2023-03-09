package entities_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	entitiesv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/discovery/entities/v1"
	"github.com/fluxninja/aperture/pkg/discovery/entities"
)

var _ = Describe("Cache", func() {
	var ec *entities.Entities

	BeforeEach(func() {
		ec = entities.NewEntities()
	})

	Context("by IP", func() {
		It("reads entity properly", func() {
			ip := "1.2.3.4"
			name := "entity_1234"
			entity := testEntity("foo", ip, name, nil)
			ec.Put(entity)
			actual, err := ec.GetByIP(ip)
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(Equal(entity))
		})

		It("returns nil when no entity found", func() {
			ip := "1.2.3.4"
			actual, err := ec.GetByIP(ip)
			Expect(err).To(Not(BeNil()))
			Expect(actual).To(BeNil())
		})

		It("removes an entity properly", func() {
			ip := "1.2.3.4"
			name := "entity_1234"
			entity := testEntity("foo", ip, name, nil)
			ec.Put(entity)

			removed := ec.Remove(entity)
			Expect(removed).To(BeTrue())

			found, err := ec.GetByIP(ip)
			Expect(err).To(Not(BeNil()))
			Expect(found).To(BeNil())
		})

		It("returns false if trying to remove a nonexistent entity", func() {
			ip := "1.2.3.4"
			otherIP := "192.168.0.1"
			name := "entity_1234"
			otherName := "other_entity_4321"
			entity := testEntity("foo", ip, name, nil)
			ec.Put(entity)

			otherEntity := testEntity("foo2", otherIP, otherName, nil)
			removed := ec.Remove(otherEntity)
			Expect(removed).To(BeFalse())

			found, err := ec.GetByIP(ip)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(Equal(entity))
		})
	})

	Context("by Name", func() {
		It("reads entity properly", func() {
			uid := "foo"
			name := "some_name"
			entity := testEntity(uid, "", name, nil)
			ec.Put(entity)
			actual, err := ec.GetByName(name)
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(Equal(entity))
		})

		It("returns nil when no entity found", func() {
			name := "some_name"
			actual, err := ec.GetByName(name)
			Expect(err).To(Not(BeNil()))
			Expect(actual).To(BeNil())
		})

		It("removes an entity properly", func() {
			uid := "bar"
			name := "some_name"
			entity := testEntity(uid, "", name, nil)
			ec.Put(entity)

			removed := ec.Remove(entity)
			Expect(removed).To(BeTrue())

			found, err := ec.GetByName(name)
			Expect(err).To(Not(BeNil()))
			Expect(found).To(BeNil())
		})

		It("returns false if trying to remove a nonexistent entity", func() {
			uid := "bar"
			name := "some_name"
			otherUid := "baz"
			otherName := "another_name"
			entity := testEntity(uid, "1.1.1.1", name, nil)
			ec.Put(entity)

			otherEntity := testEntity(otherUid, "1.1.1.2", otherName, nil)
			removed := ec.Remove(otherEntity)
			Expect(removed).To(BeFalse())

			found, err := ec.GetByName(name)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(Equal(entity))
		})
	})

	It("clears all entities from the map", func() {
		ip := "1.2.3.4"
		entity := testEntity("foo", "", "some_name", nil)
		ec.Put(entity)
		ec.Clear()
		found, err := ec.GetByIP(ip)
		Expect(err).To(Not(BeNil()))
		Expect(found).To(BeNil())
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
