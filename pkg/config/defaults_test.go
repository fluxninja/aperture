package config

import (
	"reflect"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gmeasure "github.com/onsi/gomega/gmeasure"

	"google.golang.org/protobuf/types/known/durationpb"
)

type Parent struct {
	Children []Child
}

type Child struct {
	Name string
	Age  int `default:"10"`
}

type ExampleBasic struct {
	Bool       bool    `default:"true"`
	Integer    int     `default:"33"`
	Integer8   int8    `default:"8"`
	Integer16  int16   `default:"16"`
	Integer32  int32   `default:"32"`
	Integer64  int64   `default:"64"`
	UInteger   uint    `default:"11"`
	UInteger8  uint8   `default:"18"`
	UInteger16 uint16  `default:"116"`
	UInteger32 uint32  `default:"132"`
	UInteger64 uint64  `default:"164"`
	String     string  `default:"foo"`
	Bytes      []byte  `default:"bar"`
	Float32    float32 `default:"3.2"`
	Float64    float64 `default:"6.4"`
	Struct     struct {
		Bool    bool `default:"true"`
		Integer int  `default:"33"`
	}
	Duration         time.Duration `default:"1s"`
	Children         []Child
	Second           time.Duration `default:"1s"`
	StringSlice      []string      `default:"[1,2,3,4]"`
	IntSlice         []int         `default:"[1,2,3,4]"`
	IntSliceSlice    [][]int       `default:"[[1],[2],[3],[4]]"`
	StringSliceSlice [][]string    `default:"[[1],[]]"`

	ConfigDuration      Duration             `default:"10s"`
	ConfigDurationSlice []Duration           `default:"[1s,2s,3s,4s]"`
	PBDuration          *durationpb.Duration `default:"3s"`
	ExamplePtr          *ExamplePtr
	ExamplePtrPtr       **ExamplePtr

	StringStringMap   map[string]string        `default:"{foo:bar}"`
	StringIntMap      map[string]int           `default:"{foo:1}"`
	StringDurationMap map[string]time.Duration `default:"{foo:1s}"`
}

type ExamplePtr struct {
	FloatSlice []float64 `default:"[1.1,2.2,3.3,4.4]"`
	Integer    int       `default:"33"`
}

type ExampleNested struct {
	Struct ExampleBasic
}

type ExampleMap struct {
	Values map[string]ExampleBasic
	Ptrs   map[string]*ExampleBasic
}

type FixtureTypeInt int

type ExampleStruct struct {
	X int
	T Time
	D Duration
}

var _ = Describe("Defaults", func() {
	Context("SetDefaults", func() {
		It("is threadsafe", func() {
			// Test-suite should be run with `-race` for the test to be meaningful
			done := make(chan struct{})
			go func() {
				var s ExampleStruct
				SetDefaults(&s)
				close(done)
			}()
			var s ExampleStruct
			SetDefaults(&s)
			<-done

			Eventually(done).Should(BeClosed())
		})
	})
})

var _ = Describe("DefaultTest", func() {
	var exampleBasic *ExampleBasic
	Context("set defaults with some basic example", func() {
		BeforeEach(func() {
			exampleBasic = &ExampleBasic{}
			exampleBasic.ExamplePtr = &ExamplePtr{}
		})

		It("should set defaults for basic example", func() {
			SetDefaults(exampleBasic)
			testDefaultTypes(exampleBasic)
		})

		It("should set defaults with values", func() {
			exampleBasic = &ExampleBasic{
				Integer:     55,
				UInteger:    22,
				Float32:     9.9,
				String:      "bar",
				Bytes:       []byte("foo"),
				Children:    []Child{{Name: "alice"}, {Name: "bob", Age: 2}},
				Duration:    2 * time.Second,
				PBDuration:  durationpb.New(2 * time.Second),
				StringSlice: []string{"11", "22"},
			}
			SetDefaults(exampleBasic)

			Expect(exampleBasic.Integer).To(Equal(55))
			Expect(exampleBasic.Integer8).To(Equal(int8(8)))
			Expect(exampleBasic.UInteger).To(Equal(uint(22)))
			Expect(exampleBasic.Float32).To(Equal(float32(9.9)))
			Expect(exampleBasic.String).To(Equal("bar"))
			Expect(string(exampleBasic.Bytes)).To(Equal("foo"))
			Expect(exampleBasic.Children).To(HaveLen(2))
			Expect(exampleBasic.Children[0].Name).To(Equal("alice"))
			Expect(exampleBasic.Children[0].Age).To(Equal(10))
			Expect(exampleBasic.Children[1].Name).To(Equal("bob"))
			Expect(exampleBasic.Children[1].Age).To(Equal(2))
			Expect(exampleBasic.Duration).To(Equal(2 * time.Second))
			Expect(exampleBasic.PBDuration.AsDuration()).To(Equal(2 * time.Second))
			Expect(exampleBasic.StringSlice).To(Equal([]string{"11", "22"}))
		})

		It("should set defaults efficiently", Serial, Label("benchmark measurement"), func() {
			experiment := gmeasure.NewExperiment("setDefaults performance")
			AddReportEntry(experiment.Name, experiment)

			experiment.Sample(func(i int) {
				exampleBasic = &ExampleBasic{}
				experiment.MeasureDuration("setDefaults", func() {
					SetDefaults(exampleBasic)
				})
			}, gmeasure.SamplingConfig{N: 20, Duration: 10 * time.Second})
		})
	})

	Context("set defaults with nested structs", func() {
		It("should set defaults for nested structs", func() {
			exampleNested := &ExampleNested{}
			exampleNested.Struct.ExamplePtr = &ExamplePtr{}
			SetDefaults(exampleNested)
			testDefaultTypes(&exampleNested.Struct)
		})

		It("should set defaults for structs in maps", func() {
			m := ExampleMap{
				Values: map[string]ExampleBasic{
					"foo": {},
				},
			}
			SetDefaults(&m)
			Expect(m.Values["foo"].Bool).To(BeTrue())
		})

		It("should set defaults for pointers to structs in maps", func() {
			m := ExampleMap{
				Ptrs: map[string]*ExampleBasic{
					"foo": {},
				},
			}
			SetDefaults(&m)
			Expect(m.Ptrs["foo"].Bool).To(BeTrue())
		})
	})

	var calledA bool
	var calledB bool
	Context("filler func by name and type", func() {
		BeforeEach(func() {
			calledA = false
			calledB = false
		})

		It("should fill the struct with the filler func by name", func() {
			f := &filler{
				FuncByName: map[string]fillerFunc{
					"Foo": func(field *fieldData) {
						calledA = true
					},
				},
				FuncByKind: map[reflect.Kind]fillerFunc{
					reflect.Int: func(field *fieldData) {
						calledB = true
					},
				},
			}

			f.fill(&struct{ Foo int }{})
			Expect(calledA).To(BeTrue())
			Expect(calledB).To(BeFalse())
		})

		It("should fill the struct with the filler func by type", func() {
			t := getTypeHash(reflect.TypeOf(new(FixtureTypeInt)))
			f := &filler{
				FuncByType: map[typeHash]fillerFunc{
					t: func(field *fieldData) {
						calledA = true
					},
				},
				FuncByKind: map[reflect.Kind]fillerFunc{
					reflect.Int: func(field *fieldData) {
						calledB = true
					},
				},
			}

			f.fill(&struct{ Foo FixtureTypeInt }{})
			Expect(calledA).To(BeTrue())
			Expect(calledB).To(BeFalse())
			// Test FuncByKind Slice
			Expect(getTypeHash(reflect.TypeOf(new([]string)))).To(Equal(typeHash(".")))
		})
	})

	var called string
	Context("filler func by kind", func() {
		BeforeEach(func() {
			calledA = false
		})

		It("should fill the struct with the filler func by kind", func() {
			f := &filler{
				FuncByKind: map[reflect.Kind]fillerFunc{
					reflect.Int: func(field *fieldData) {
						calledA = true
					},
				},
			}

			f.fill(&struct{ Foo int }{})
			Expect(calledA).To(BeTrue())
		})

		It("should fill the struct with the filler func by kind tag", func() {
			f := &filler{
				Tag: "foo",
				FuncByKind: map[reflect.Kind]fillerFunc{
					reflect.Int: func(field *fieldData) {
						called = field.TagValue
					},
				},
			}

			f.fill(&struct {
				Foo int `foo:"qux"`
			}{})
			Expect(called).To(Equal("qux"))
		})
	})
})

func testDefaultTypes(foo *ExampleBasic) {
	Expect(foo.Bool).To(BeTrue())
	Expect(foo.Integer).To(Equal(33))
	Expect(foo.Integer8).To(Equal(int8(8)))
	Expect(foo.Integer16).To(Equal(int16(16)))
	Expect(foo.Integer32).To(Equal(int32(32)))
	Expect(foo.Integer64).To(Equal(int64(64)))
	Expect(foo.UInteger).To(Equal(uint(11)))
	Expect(foo.UInteger8).To(Equal(uint8(18)))
	Expect(foo.UInteger16).To(Equal(uint16(116)))
	Expect(foo.UInteger32).To(Equal(uint32(132)))
	Expect(foo.UInteger64).To(Equal(uint64(164)))
	Expect(foo.String).To(Equal("foo"))
	Expect(string(foo.Bytes)).To(Equal("bar"))
	Expect(foo.Float32).To(Equal(float32(3.2)))
	Expect(foo.Float64).To(Equal(6.4))
	Expect(foo.Struct.Bool).To(BeTrue())
	Expect(foo.Duration).To(Equal(time.Second))
	Expect(foo.Children).To(BeNil())
	Expect(foo.Second).To(Equal(time.Second))
	Expect(foo.StringSlice).To(Equal([]string{"1", "2", "3", "4"}))
	Expect(foo.IntSlice).To(Equal([]int{1, 2, 3, 4}))
	Expect(foo.IntSliceSlice).To(Equal([][]int{{1}, {2}, {3}, {4}}))
	Expect(foo.StringSliceSlice).To(Equal([][]string{{"1"}, {}}))
	Expect(foo.ConfigDuration.AsDuration()).To(Equal(time.Second * 10))
	Expect(foo.ConfigDurationSlice).To(HaveLen(4))
	Expect(foo.ConfigDurationSlice[0].AsDuration()).To(Equal(time.Second * 1))
	Expect(foo.ConfigDurationSlice[1].AsDuration()).To(Equal(time.Second * 2))
	Expect(foo.ConfigDurationSlice[2].AsDuration()).To(Equal(time.Second * 3))
	Expect(foo.ConfigDurationSlice[3].AsDuration()).To(Equal(time.Second * 4))
	Expect(foo.PBDuration.AsDuration()).To(Equal(time.Second * 3))
	Expect(foo.ExamplePtr).NotTo(BeNil())
	Expect(foo.ExamplePtr.Integer).To(Equal(33))
	Expect(foo.ExamplePtr.FloatSlice).To(Equal([]float64{1.1, 2.2, 3.3, 4.4}))
	Expect(foo.StringStringMap).To(Equal(map[string]string{"foo": "bar"}))
	Expect(foo.StringIntMap).To(Equal(map[string]int{"foo": 1}))
	Expect(foo.StringDurationMap).To(Equal(map[string]time.Duration{"foo": time.Second}))
}
