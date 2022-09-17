package static

import (
	"encoding/json"
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	entitycachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/entitycache/v1"
	"github.com/fluxninja/aperture/pkg/mocks"
	"github.com/fluxninja/aperture/pkg/notifiers"
)

var (
	ctrl         *gomock.Controller
	mockTrackers *mocks.MockTrackers
)

var _ = BeforeEach(func() {
	ctrl = gomock.NewController(GinkgoT())
	mockTrackers = mocks.NewMockTrackers(ctrl)
})

var _ = Describe("Static service discovery", func() {
	Context("Discovery from config", func() {
		It("Writes no entities with nil service list", func() {
			config := StaticDiscoveryConfig{
				Services: nil,
			}

			sd, trackers := CreateStaticDiscoveryWithFakeTracker(config)

			trackers.EXPECT().WriteEvent(gomock.Any(), gomock.Any()).MaxTimes(0)

			err := sd.start()
			Expect(err).NotTo(HaveOccurred())
		})

		It("Correctly reads a config entity", func() {
			someIPAddress := "1.2.3.4"
			someUID := "foo"
			someName := "some_entity"
			someService := "svc1"

			config := StaticDiscoveryConfig{
				Services: []*ServiceConfig{
					{
						Name: someService,
						Entities: []*EntityConfig{
							{
								IPAddress: someIPAddress,
								UID:       someUID,
								Name:      someName,
							},
						},
					},
				},
			}

			expectedEntity := &entitycachev1.Entity{
				IpAddress: someIPAddress,
				Prefix:    staticEntityTrackerPrefix,
				Uid:       someUID,
				Services:  []string{someService},
				Name:      someName,
			}

			expectedEntityKey := notifiers.Key(fmt.Sprintf("%v.%v", staticEntityTrackerPrefix, someUID))
			serializedExpectedEntity, err := json.Marshal(expectedEntity)
			Expect(err).NotTo(HaveOccurred())

			sd, trackers := CreateStaticDiscoveryWithFakeTracker(config)

			trackers.EXPECT().WriteEvent(expectedEntityKey, serializedExpectedEntity).Times(1)

			err = sd.start()
			Expect(err).NotTo(HaveOccurred())
		})

		It("Writes all entities defined in config", func() {
			config := StaticDiscoveryConfig{
				Services: []*ServiceConfig{
					{
						Name: "svc1",
						Entities: []*EntityConfig{
							{
								IPAddress: "1.2.3.4",
								UID:       "foo",
								Name:      "someName",
							},
							{
								IPAddress: "1.2.3.5",
								UID:       "bar",
								Name:      "someName",
							},
						},
					},
				},
			}

			sd, trackers := CreateStaticDiscoveryWithFakeTracker(config)

			trackers.EXPECT().WriteEvent(gomock.Any(), gomock.Any()).Times(2)

			err := sd.start()
			Expect(err).NotTo(HaveOccurred())
		})

		It("Writes one entity if it's defined for multiple services", func() {
			someIPAddress := "1.2.3.4"
			someUID := "foo"
			someName := "some_entity"
			someService := "svc1"
			someOtherService := "svc2"

			config := StaticDiscoveryConfig{
				Services: []*ServiceConfig{
					{
						Name: someService,
						Entities: []*EntityConfig{
							{
								IPAddress: someIPAddress,
								UID:       someUID,
								Name:      someName,
							},
						},
					},
					{
						Name: someOtherService,
						Entities: []*EntityConfig{
							{
								IPAddress: someIPAddress,
								UID:       someUID,
								Name:      someName,
							},
						},
					},
				},
			}

			expectedEntity := &entitycachev1.Entity{
				IpAddress: someIPAddress,
				Prefix:    staticEntityTrackerPrefix,
				Uid:       someUID,
				Services:  []string{someService, someOtherService},
				Name:      someName,
			}

			expectedEntityKey := notifiers.Key(fmt.Sprintf("%v.%v", staticEntityTrackerPrefix, someUID))
			serializedExpectedEntity, err := json.Marshal(expectedEntity)
			Expect(err).NotTo(HaveOccurred())

			sd, trackers := CreateStaticDiscoveryWithFakeTracker(config)

			trackers.EXPECT().WriteEvent(expectedEntityKey, serializedExpectedEntity).Times(1)

			err = sd.start()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

func CreateStaticDiscoveryWithFakeTracker(config StaticDiscoveryConfig) (*StaticDiscovery, *mocks.MockTrackers) {
	sd, err := newStaticServiceDiscovery(mockTrackers, config)
	Expect(err).NotTo(HaveOccurred())
	return sd, mockTrackers
}
