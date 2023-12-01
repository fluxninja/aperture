package status

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/wrapperspb"

	statusv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/status/v1"
	"github.com/fluxninja/aperture/v2/pkg/alerts"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

var rootRegistry Registry

var _ = Describe("Status Registry", func() {
	BeforeEach(func() {
		alerter := alerts.NewSimpleAlerter(100)
		rootRegistry = NewRegistry(log.GetGlobalLogger(), alerter)
	})

	Context("there is a single level root registry", func() {
		It("returns the root registry itself", func() {
			Expect(rootRegistry.Parent()).To(BeNil())
			Expect(rootRegistry.Root()).To(Equal(rootRegistry))
			Expect(rootRegistry.ChildIfExists("test", rootRegistry.Value())).To(BeNil())
		})
		It("returns empty status information", func() {
			Expect(rootRegistry.GetStatus()).To(Equal(&statusv1.Status{}))
			Expect(rootRegistry.GetGroupStatus()).To(Equal(&statusv1.GroupStatus{
				Status: &statusv1.Status{},
				Groups: make(map[string]*statusv1.GroupStatus),
			}))
		})
		It("returns updated status information", func() {
			test_status := NewStatus(nil, errors.New("test status"))
			rootRegistry.SetStatus(test_status)
			Expect(rootRegistry.GetStatus()).To(Equal(test_status))
			Expect(rootRegistry.HasError()).To(BeTrue())
			rootRegistry.SetStatus(nil)
			Expect(rootRegistry.GetStatus()).To(Equal(&statusv1.Status{}))
			Expect(rootRegistry.HasError()).To(BeFalse())

			test_groupstatus := &statusv1.GroupStatus{
				Status: rootRegistry.GetStatus(),
				Groups: make(map[string]*statusv1.GroupStatus),
			}
			rootRegistry.SetGroupStatus(test_groupstatus)
			Expect(rootRegistry.GetGroupStatus()).To(Equal(test_groupstatus))
		})
		It("creates a new child registry and then detaches it", func() {
			rootRegistry.Detach()
			child_registry := rootRegistry.Child("test", "child")
			Expect(rootRegistry.ChildIfExists("test", "child")).To(Equal(child_registry))
			Expect(child_registry.URI()).To(Equal("/test/child"))
			rootRegistry.Child("test", "child").Detach()
			Expect(rootRegistry.ChildIfExists("test", "child")).To(BeNil())
		})
		It("returns the key of the root registry", func() {
			Expect(rootRegistry.Value()).To(Equal("root"))
		})
	})

	var (
		child1      Registry
		child2      Registry
		grandChild1 Registry
		grandChild2 Registry
	)
	Context("there are multiple registries in hierarchy", func() {
		BeforeEach(func() {
			child1 = rootRegistry.Child("test", "child1")
			child2 = rootRegistry.Child("test", "child2")
			grandChild1 = child1.Child("test", "grandChild1")
			grandChild2 = child2.Child("test", "grandChild2")
		})

		It("returns root, parent, child registry", func() {
			Expect(child1.Parent()).To(Equal(rootRegistry))
			Expect(child2.Parent()).To(Equal(rootRegistry))
			Expect(grandChild1.Parent()).To(Equal(child1))
			Expect(grandChild2.Parent()).To(Equal(child2))
			Expect(child1.Root()).To(Equal(rootRegistry))
			Expect(child2.Root()).To(Equal(rootRegistry))
			Expect(grandChild1.Root()).To(Equal(rootRegistry))
			Expect(grandChild2.Root()).To(Equal(rootRegistry))
			Expect(child1.Child("test", "grandChild1")).To(Equal(grandChild1))
			Expect(grandChild1.URI()).To(Equal("/test/child1/test/grandChild1"))
			Expect(child2.Child("test", "grandChild2")).To(Equal(grandChild2))
			Expect(grandChild2.URI()).To(Equal("/test/child2/test/grandChild2"))
			Expect(child1.ChildIfExists("test", "grandChild1")).To(Equal(grandChild1))
			Expect(child2.ChildIfExists("test", "grandChild2")).To(Equal(grandChild2))
			Expect(rootRegistry.Child("test", "child1")).To(Equal(child1))
			Expect(rootRegistry.Child("test", "child2")).To(Equal(child2))
			Expect(rootRegistry.ChildIfExists("test", "child1")).To(Equal(child1))
			Expect(rootRegistry.ChildIfExists("test", "child2")).To(Equal(child2))
		})
		It("returns updated status information", func() {
			test_status1 := NewStatus(nil, errors.New(""))
			rootRegistry.SetStatus(test_status1)
			Expect(rootRegistry.GetStatus()).To(Equal(test_status1))
			Expect(rootRegistry.HasError()).To(BeTrue())

			test_status2 := NewStatus(wrapperspb.String("test status2"), nil)
			test_groupstatus1 := &statusv1.GroupStatus{
				Status: test_status2,
				Groups: make(map[string]*statusv1.GroupStatus),
			}
			child1.SetStatus(test_status1)
			child1.SetGroupStatus(test_groupstatus1)
			Expect(child1.HasError()).To(BeFalse())
			Expect(child1.GetStatus()).To(Equal(test_status2))
			Expect(grandChild1.GetStatus()).To(Equal(&statusv1.Status{}))

			test_status3 := NewStatus(nil, errors.New("test status3"))
			test_groupstatus2 := &statusv1.GroupStatus{
				Status: test_status3,
				Groups: make(map[string]*statusv1.GroupStatus),
			}

			child2.SetGroupStatus(test_groupstatus2)
			Expect(child2.HasError()).To(BeTrue())
			Expect(child2.GetStatus()).To(Equal(test_status3))
			Expect(grandChild2.GetStatus()).To(Equal(&statusv1.Status{}))
			Expect(rootRegistry.HasError()).To(BeTrue())

			rootRegistry.SetGroupStatus(test_groupstatus2)
			rootGroupStatus := rootRegistry.GetGroupStatus().Status
			Expect(rootGroupStatus).To(Equal(test_status3))
			rootRegistry.SetStatus(nil)
			Expect(rootRegistry.HasError()).To(BeTrue())
		})
		It("creates multiple child registries then detaches them", func() {
			grandChild3 := child1.Child("test", "grandChild3")
			grandChild4 := child1.Child("test", "grandChild4")
			ggrandChild1 := grandChild3.Child("test", "ggrandChild1")
			ggrandChild2 := grandChild4.Child("test", "ggrandChild2")
			Expect(child1.ChildIfExists("test", "grandChild3")).To(Equal(grandChild3))
			Expect(child1.ChildIfExists("test", "grandChild4")).To(Equal(grandChild4))
			Expect(grandChild3.ChildIfExists("test", "ggrandChild1")).To(Equal(ggrandChild1))
			Expect(grandChild4.ChildIfExists("test", "ggrandChild2")).To(Equal(ggrandChild2))

			child1.Child("test", "grandChild3").Detach()
			Expect(child1.ChildIfExists("test", "grandChild3")).To(BeNil())
			grandChild4.Child("test", "ggrandChild2").Detach()
			Expect(grandChild4.ChildIfExists("test", "ggrandChild2")).To(BeNil())
			child1.Child("test", "grandChild4").Detach()
			Expect(child1.ChildIfExists("test", "grandChild4")).To(BeNil())
		})
		It("returns the key of the registry that is registered with the parent", func() {
			Expect(child1.Value()).To(Equal("child1"))
			Expect(child2.Value()).To(Equal("child2"))
			Expect(grandChild1.Value()).To(Equal("grandChild1"))
			Expect(grandChild2.Value()).To(Equal("grandChild2"))
		})
	})
})
