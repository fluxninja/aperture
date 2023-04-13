package static

import (
	"encoding/json"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	entitiesv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/discovery/entities/v1"
	"github.com/fluxninja/aperture/pkg/config"
	sdconfig "github.com/fluxninja/aperture/pkg/discovery/static/config"
	"github.com/fluxninja/aperture/pkg/mocks"
	"github.com/fluxninja/aperture/pkg/notifiers"
)

var (
	ctrl            *gomock.Controller
	mockEventWriter *mocks.MockEventWriter
)

var _ = BeforeEach(func() {
	ctrl = gomock.NewController(GinkgoT())
	mockEventWriter = mocks.NewMockEventWriter(ctrl)
})

var _ = Describe("Static service discovery", func() {
	Context("Discovery from config", func() {
		It("Writes no entities with nil service list", func() {
			cfg := &sdconfig.StaticDiscoveryConfig{
				Entities: nil,
			}

			sd, trackers := CreateStaticDiscoveryWithFakeTracker(cfg)

			trackers.EXPECT().WriteEvent(gomock.Any(), gomock.Any()).MaxTimes(0)

			err := sd.start()
			Expect(err).NotTo(HaveOccurred())
		})

		It("Correctly reads a config entity", func() {
			someIPAddress := "1.2.3.4"
			someUID := "foo"
			someName := "some_entity"
			someService := "svc1"

			cfg := &sdconfig.StaticDiscoveryConfig{
				Entities: []entitiesv1.Entity{
					{
						IpAddress: someIPAddress,
						Uid:       someUID,
						Services:  []string{someService},
						Name:      someName,
					},
				},
			}

			expectedEntity := &entitiesv1.Entity{
				IpAddress: someIPAddress,
				Uid:       someUID,
				Services:  []string{someService},
				Name:      someName,
			}

			expectedEntityKey := notifiers.Key(someUID)
			serializedExpectedEntity, err := json.Marshal(expectedEntity)
			Expect(err).NotTo(HaveOccurred())

			sd, trackers := CreateStaticDiscoveryWithFakeTracker(cfg)

			trackers.EXPECT().WriteEvent(expectedEntityKey, serializedExpectedEntity).Times(1)

			err = sd.start()
			Expect(err).NotTo(HaveOccurred())
		})

		It("Writes all entities defined in config", func() {
			bytes := []byte(`
      {
        "entities": [
          {
            "ip_address": "1.2.3.4",
            "uid": "foo",
            "name": "someName",
            "services": ["svc1"]
          },
          {
            "ip_address": "1.2.3.5",
            "uid": "bar",
            "name": "someName",
            "services": ["svc1"]
          }
        ]
      }
      `)

			// use unmarshaller
			unmarshaller, err := config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller(bytes)
			Expect(err).NotTo(HaveOccurred())
			var cfg sdconfig.StaticDiscoveryConfig
			err = unmarshaller.Unmarshal(&cfg)
			Expect(err).NotTo(HaveOccurred())

			sd, trackers := CreateStaticDiscoveryWithFakeTracker(&cfg)

			trackers.EXPECT().WriteEvent(gomock.Any(), gomock.Any()).Times(2)

			err = sd.start()
			Expect(err).NotTo(HaveOccurred())
		})

		It("Writes one entity if it's defined for multiple services", func() {
			someIPAddress := "1.2.3.4"
			someUID := "foo"
			someName := "some_entity"
			someService := "svc1"
			someOtherService := "svc2"

			cfg := &sdconfig.StaticDiscoveryConfig{
				Entities: []entitiesv1.Entity{
					{
						IpAddress: someIPAddress,
						Uid:       someUID,
						Services:  []string{someService, someOtherService},
						Name:      someName,
					},
				},
			}

			expectedEntity := &entitiesv1.Entity{
				IpAddress: someIPAddress,
				Uid:       someUID,
				Services:  []string{someService, someOtherService},
				Name:      someName,
			}

			expectedEntityKey := notifiers.Key(someUID)
			serializedExpectedEntity, err := json.Marshal(expectedEntity)
			Expect(err).NotTo(HaveOccurred())

			sd, trackers := CreateStaticDiscoveryWithFakeTracker(cfg)

			trackers.EXPECT().WriteEvent(expectedEntityKey, serializedExpectedEntity).Times(1)

			err = sd.start()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

func CreateStaticDiscoveryWithFakeTracker(
	config *sdconfig.StaticDiscoveryConfig,
) (*StaticDiscovery, *mocks.MockEventWriter) {
	return newStaticServiceDiscovery(mockEventWriter, config), mockEventWriter
}
