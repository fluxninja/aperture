package config

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
)

type simpleTestConfig struct {
	Name string
	Type string
}

type nestedTestConfig struct {
	Name       string
	TestConfig struct {
		testName string
		testBool bool
	}
	Type int
}

type defaultTestConfig struct {
	Type   *string
	Name   string
	Values []int `json:"values" default:"[1,2,3]"`
	Val    int
}

type testStruct struct {
	Configs []defaultTestConfig `json:"configs"`
}

var _ = Describe("Koanf-unmarshaller", func() {
	Context("simple config struct", func() {
		var koanf *KoanfUnmarshaller

		BeforeEach(func() {
			bytes := []byte("koanf-unmarshaller: test")
			unmarshaller, err := KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller(bytes)
			Expect(err).NotTo(HaveOccurred())

			koanf = unmarshaller.(*KoanfUnmarshaller)
		})
		It("config key is set and updated accordingly", func() {
			Expect(koanf.IsSet("koanf-unmarshaller")).To(BeTrue())
			Expect(koanf.Get("koanf-unmarshaller")).To(Equal("test"))

			err := koanf.Reload([]byte("koanf-test: unmarshaller"))
			Expect(err).NotTo(HaveOccurred())
			Expect(koanf.IsSet("koanf-test")).To(BeTrue())
			Expect(koanf.Get("koanf-test")).To(Equal("unmarshaller"))
		})

		It("unmarshal using the underlying koanf", func() {
			var testConfig simpleTestConfig
			err := koanf.Unmarshal(&testConfig)
			Expect(err).NotTo(HaveOccurred())

			cm := koanf.Get("")
			koanfMap := map[string]interface{}{
				"koanf-unmarshaller": "test",
			}
			Expect(cm).To(Equal(koanfMap))
		})

		It("unmarshalKey binds the interface to key path in config map", func() {
			var testConfig simpleTestConfig
			err := koanf.UnmarshalKey("koanf2", &testConfig)
			Expect(err).NotTo(HaveOccurred())

			cm := koanf.Get("")
			koanfMap := map[string]interface{}{
				"koanf-unmarshaller": "test",
			}
			Expect(cm).To(Equal(koanfMap))
		})
	})

	Context("nested config struct", func() {
		var koanf *KoanfUnmarshaller

		BeforeEach(func() {
			bytes := []byte("koanf1: test1\nkoanf2: test2\nkoanf3: test3")
			unmarshaller, err := KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller(bytes)
			Expect(err).NotTo(HaveOccurred())

			koanf = unmarshaller.(*KoanfUnmarshaller)
		})

		It("config key is set and updated accordingly", func() {
			Expect(koanf.IsSet("koanf1")).To(BeTrue())
			Expect(koanf.Get("koanf2")).To(Equal("test2"))
			Expect(koanf.Get("koanf3")).To(Equal("test3"))

			err := koanf.Reload([]byte("koanf4: test4\nkoanf5: test5\nkoanf6: test6"))
			Expect(err).NotTo(HaveOccurred())
			Expect(koanf.IsSet("koanf4")).To(BeTrue())
			Expect(koanf.Get("koanf5")).To(Equal("test5"))
			Expect(koanf.Get("koanf6")).To(Equal("test6"))
		})

		It("unmarshal using the underlying koanf", func() {
			var testConfig nestedTestConfig
			err := koanf.Unmarshal(&testConfig)
			Expect(err).NotTo(HaveOccurred())

			cm := koanf.Get("")
			koanfMap := map[string]interface{}{
				"koanf1": "test1",
				"koanf2": "test2",
				"koanf3": "test3",
			}
			Expect(cm).To(Equal(koanfMap))
		})

		It("unmarshalKey binds the interface to key path in config map", func() {
			var testConfig nestedTestConfig
			err := koanf.UnmarshalKey("koanf", &testConfig)
			Expect(err).NotTo(HaveOccurred())

			cm := koanf.Get("")
			koanfMap := map[string]interface{}{
				"koanf1": "test1",
				"koanf2": "test2",
				"koanf3": "test3",
			}
			Expect(cm).To(Equal(koanfMap))
		})
	})

	Context("BindEnv is enabled", func() {
		var koanf *KoanfUnmarshaller
		const nameKey string = "CONFIG_TEST_NAME"
		const typeKey string = "CONFIG_TEST_TYPE"
		const nameValue string = "name"
		const typeValue string = "type"

		BeforeEach(func() {
			unmarshaller, err := KoanfUnmarshallerConstructor{
				EnableEnv: true,
			}.NewKoanfUnmarshaller(nil)
			Expect(err).NotTo(HaveOccurred())

			koanf = unmarshaller.(*KoanfUnmarshaller)

			err = os.Setenv(nameKey, nameValue)
			Expect(err).NotTo(HaveOccurred())
			err = os.Setenv(typeKey, typeValue)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			err := os.Unsetenv(nameKey)
			Expect(err).NotTo(HaveOccurred())
			err = os.Unsetenv(typeKey)
			Expect(err).NotTo(HaveOccurred())
		})

		It("binds environment variables correctly", func() {
			var testConfig simpleTestConfig
			err := koanf.Unmarshal(&testConfig)
			Expect(err).NotTo(HaveOccurred())

			Expect(testConfig.Name).To(Equal(nameValue))
			Expect(testConfig.Type).To(Equal(typeValue))
		})
	})

	Context("BindEnv is enabled and non supported parateres are passed", func() {
		var koanf *KoanfUnmarshaller
		var nestedConfig nestedTestConfig
		s := "test"
		bytes := []byte("")
		unmarshaller, err := KoanfUnmarshallerConstructor{
			EnableEnv: true,
		}.NewKoanfUnmarshaller(bytes)
		Expect(err).NotTo(HaveOccurred())
		koanf = unmarshaller.(*KoanfUnmarshaller)

		It("covers the case in bindEnvsKey when a pointer is nil and points it to a new Zero Value", func() {
			nilConfig := &defaultTestConfig{
				Name: s,
				Type: nil,
			}
			err := koanf.Unmarshal(nilConfig)
			Expect(err).ToNot(HaveOccurred())
		})

		It("covers the case when an interface has a non nil pointer which gets converted into a value that can be used, and other type of datatypes are also evaluated", func() {
			nilConfig := &defaultTestConfig{
				Name: s,
				Type: &s,
			}
			err := koanf.Unmarshal(nilConfig)
			Expect(err).ToNot(HaveOccurred())
		})
		It("covers the case in bindEnvsKey when an embedded struct is passed, and check if any of them is duration or timestamp", func() {
			err := koanf.Unmarshal(&nestedConfig)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when there is incoming map on top of existing map", func() {
		var koanf *KoanfUnmarshaller

		var a, b map[string]interface{}
		var ok bool

		BeforeEach(func() {
			bytes := []byte("koanf-unmarshaller: test")
			unmarshaller, err := KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller(bytes)
			Expect(err).NotTo(HaveOccurred())
			koanf = unmarshaller.(*KoanfUnmarshaller)

			cm := koanf.Get("")
			if a, ok = cm.(map[string]interface{}); ok {
				Expect(a["koanf-unmarshaller"]).To(Equal("test"))
			}
		})

		It("merges when both source and target keys are both maps", func() {
			bytes := []byte("koanf-merge: maps")
			unmarshaller, err := KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller(bytes)
			Expect(err).NotTo(HaveOccurred())
			koanfB := unmarshaller.(*KoanfUnmarshaller)

			cmB := koanfB.Get("")
			if b, ok = cmB.(map[string]interface{}); ok {
				Expect(b["koanf-merge"]).To(Equal("maps"))
			}

			merge(a, b, true)
			cm := map[string]interface{}{
				"koanf-merge":        "maps",
				"koanf-unmarshaller": "test",
			}
			Expect(b).To(Equal(cm))
		})

		It("overrides when source key and target key are equal", func() {
			bytes := []byte("koanf-unmarshaller: maps")
			unmarshaller, err := KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller(bytes)
			Expect(err).NotTo(HaveOccurred())
			koanfB := unmarshaller.(*KoanfUnmarshaller)

			cmB := koanfB.Get("")
			if b, ok = cmB.(map[string]interface{}); ok {
				Expect(b["koanf-unmarshaller"]).To(Equal("maps"))
			}

			merge(a, b, true)
			cm := map[string]interface{}{
				"koanf-unmarshaller": "test",
			}
			Expect(a).To(Equal(cm))
			Expect(a).To(Equal(cm))
		})
	})

	Context("When Testing unmarshaller generic Get functions", func() {
		var testConfig nestedTestConfig
		bytes := []byte("koanf-unmarshaller: test")
		unmarshaller, err := KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller(bytes)
		Expect(err).NotTo(HaveOccurred())
		It("returns the correct key values parsed from the unmarshaller", func() {
			val := GetValue(unmarshaller, "koanf-unmarshaller", &testConfig)
			Expect(val).To(Equal("test"))
			str := GetStringValue(unmarshaller, "koanf-unmarshaller", "test")
			Expect(str).To(Equal("test"))
			integer := GetIntValue(unmarshaller, "koanf-unmarshaller", 0)
			Expect(integer).To(Equal(0))
			bool := GetBoolValue(unmarshaller, "koanf-unmarshaller", true)
			Expect(bool).To(BeFalse())
		})
	})

	Context("using a basic unmarshaller", func() {
		bytes := []byte("koanf-unmarshaller: test")
		mp := map[string]interface{}{
			"koanf-unmarshaller": "test",
		}
		It("returns an error when unmarshalling a map", func() {
			err := UnmarshalYAML(bytes, &mp)
			Expect(err).To(HaveOccurred())
		})
	})
	Context("when testing for default overrides", func() {
		s := "test"
		defaultConfig := defaultTestConfig{
			Name: s,
			Type: &s,
			Val:  5,
		}
		bytes := []byte("Val: 0")
		unmarshaller, err := KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller(bytes)
		Expect(err).NotTo(HaveOccurred())
		It("returns the overridden value", func() {
			err = unmarshaller.Unmarshal(&defaultConfig)
			Expect(err).ToNot(HaveOccurred())
			Expect(defaultConfig.Val).To(Equal(int(0)))
		})
		It("should reload Val and change its value", func() {
			bytes2 := []byte("Val: 97")
			unmarshaller.Reload(bytes2)
			err = unmarshaller.Unmarshal(&defaultConfig)
			Expect(err).ToNot(HaveOccurred())
			Expect(defaultConfig.Val).To(Equal(int(97)))
		})
	})
	Context("when testing for a slices of ints within a struct", func() {
		s := "test"
		defaultConfig := &defaultTestConfig{
			Name: s,
			Type: &s,
			Val:  5,
		}

		bytes := []byte("values: [4,5,6]")
		unmarshaller, err := KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller(bytes)
		Expect(err).NotTo(HaveOccurred())
		It("overrides values within a struct", func() {
			err = unmarshaller.Unmarshal(defaultConfig)
			Expect(err).ToNot(HaveOccurred())
			Expect(defaultConfig.Values).To(Equal([]int{4, 5, 6}))
		})
		It("changes the values of a slice", func() {
			bytes2 := []byte("values: [7,8,9]")
			unmarshaller.Reload(bytes2)
			err = unmarshaller.Unmarshal(defaultConfig)
			Expect(err).ToNot(HaveOccurred())
			Expect(defaultConfig.Values).To(Equal([]int{7, 8, 9}))
			defaultConfig.Values = []int{10, 11, 12}
			Expect(defaultConfig.Values).To(Equal([]int{10, 11, 12}))
		})
		It("checks if defaults are loaded correctly in a new config", func() {
			defaultConfig2 := &defaultTestConfig{}
			bytes2 := []byte("")
			unmarshaller, err := KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller(bytes2)
			Expect(err).NotTo(HaveOccurred())
			err = unmarshaller.Unmarshal(defaultConfig2)
			Expect(err).ToNot(HaveOccurred())
			Expect(defaultConfig2.Values).To(Equal([]int{1, 2, 3}))
		})
	})
	Context("when testing a slice of structs by loading a yaml configuration in it", func() {
		test := testStruct{}

		bytes := []byte(`
          configs:
            - Name: test1
              Val: 1
              values:
                - 1
                - 2
                - 3
                - 4
                - 5
            - Name: test2
              Val: 2`)
		It("loads values in the first struct and defaults in the other", func() {
			unmarshaller, err := KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller(bytes)
			Expect(err).NotTo(HaveOccurred())
			err = unmarshaller.Unmarshal(&test)
			Expect(err).ToNot(HaveOccurred())
			Expect(test.Configs[0].Val).To(Equal(int(1)))
			Expect(test.Configs[0].Values).To(Equal([]int{1, 2, 3, 4, 5}))
			Expect(test.Configs[1].Val).To(Equal(int(2)))
			Expect(test.Configs[1].Values).To(Equal([]int{1, 2, 3}))
		})
	})
})

var _ = Describe("TimeStamp", func() {
	Context("when marshalling same timeStamp twice", func() {
		ts := &Time{
			timestamp: timestamppb.Now(),
		}
		b, err := ts.MarshalJSON()
		b2, err2 := ts.MarshalJSON()
		It("should return nil error and match both timeStamps", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
			Expect(b).To(Equal(b2))
		})
	})
	Context("when unmarshalling a timestamp", func() {
		ts := &Time{
			timestamp: timestamppb.Now(),
		}
		ts2 := &Time{}
		b, err := ts.MarshalJSON()
		errBytes := []byte(`{"timestamp":"` + ts.timestamp.String() + `"}`)
		It("should return nil error, match both timestamps and their string values", func() {
			Expect(err).NotTo(HaveOccurred())
			err = ts2.UnmarshalJSON(b)
			Expect(err).NotTo(HaveOccurred())
			Expect(ts2.timestamp).To(Equal(ts.timestamp))
			Expect(ts2.String()).To(Equal(ts.String()))
		})
		It("should return error when trying to parse an unsupported timestamp", func() {
			err = ts2.UnmarshalJSON(errBytes)
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("Duration", func() {
	Context("when marshalling 2 durations", func() {
		dur := MakeDuration(1 * time.Second)
		b, err := dur.MarshalJSON()
		dur2 := MakeDuration(1 * time.Second)
		b2, err2 := dur2.MarshalJSON()
		It("should return nil error and match both durations", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
			Expect(b).To(Equal(b2))
		})
		It("should throw error for two different durations", func() {
			dur2 = MakeDuration(2 * time.Second)
			b2, err2 = dur2.MarshalJSON()
			Expect(err2).NotTo(HaveOccurred())
			Expect(b).NotTo(Equal(b2))
		})
	})

	Context("when unmarshalling a duration", func() {
		dur := MakeDuration(1 * time.Second)
		dur2 := MakeDuration(0)
		b, err := dur.MarshalJSON()
		errBytes := []byte(`{"duration":"` + dur.String() + `"}`)
		It("should return nil error, match both durations and their string values", func() {
			Expect(err).NotTo(HaveOccurred())
			err := dur2.UnmarshalJSON(b)
			Expect(err).NotTo(HaveOccurred())
			Expect(dur2.AsDuration()).To(Equal(dur.AsDuration()))
			Expect(dur2.String()).To(Equal(dur.String()))
		})
		It("should return error when trying to parse an unsupported duration", func() {
			err := dur2.UnmarshalJSON(errBytes)
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("ProtobufUnmarshaller", func() {
	Context("when unmarshalling a protobuf", func() {
		flowSelector := &policylangv1.FlowSelector{
			ServiceSelector: &policylangv1.ServiceSelector{
				AgentGroup: "ag",
				Service:    "s.n.svc.cluster.local",
			},
			FlowMatcher: &policylangv1.FlowMatcher{
				ControlPoint: "egress",
			},
		}
		selectorBytes, err := proto.Marshal(flowSelector)
		Expect(err).NotTo(HaveOccurred())

		unmarshaller, err := NewProtobufUnmarshaller(selectorBytes)
		Expect(err).NotTo(HaveOccurred())

		It("parses selectorBytes content into newSel and matches both contents", func() {
			var newSel policylangv1.FlowSelector
			err := unmarshaller.Unmarshal(&newSel)
			Expect(err).NotTo(HaveOccurred())
			Expect(newSel.String()).To(Equal(flowSelector.String()))
		})

		It("should return an error when unmarshalling a non-protobuf message", func() {
			err := unmarshaller.Unmarshal("test, non-protobuf message")
			Expect(err).To(HaveOccurred())
		})

		It("should return an error when trying to reload with nil bytes", func() {
			err := unmarshaller.Reload(nil)
			Expect(err).To(HaveOccurred())
		})

		It("should recover the panic when calling UnmarshalKey", func() {
			defer func() {
				recover()
			}()
			unmarshaller.UnmarshalKey("panic recover test", nil)
		})

		It("should recover the panic when calling IsSet", func() {
			defer func() {
				recover()
			}()
			unmarshaller.IsSet("panic recover test")
		})

		It("should recover the panic when calling Get", func() {
			defer func() {
				recover()
			}()
			unmarshaller.Get("panic recover test")
		})
	})
})
