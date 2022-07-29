package component_test

/*import (
	"reflect"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/durationpb"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/component"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/reading"
)

var _ = Describe("EMA filter", func() {
	var (
		emaFilter *component.EMA
		prevInput reading.Reading
	)

	BeforeEach(func() {
		emaProto := &policylangv1.EMA{
			EmaWindow:    durationpb.New(2000 * time.Millisecond),
			WarmUpWindow: durationpb.New(500 * time.Millisecond),
		}
		prevInput = reading.New(1.0)

		emptyEma := &component.EMA{}
		emaFilter, options, err := component.NewEMAAndOptions(emaProto, 0, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(reflect.TypeOf(emaFilter)).To(Equal(reflect.TypeOf(emptyEma)))
		Expect(options).To(BeNil())
	})
})*/
