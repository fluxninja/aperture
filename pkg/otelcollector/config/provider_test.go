package otelconfig_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/confmap"
	"golang.org/x/exp/maps"

	otelconfig "github.com/fluxninja/aperture/v2/pkg/otelcollector/config"
)

var _ = Describe("Provider", func() {
	It("triggers collector config update", func() {
		triggered := false
		onUpdate := func(*confmap.ChangeEvent) {
			triggered = true
		}

		provider := otelconfig.NewProvider("foo")
		provider.Retrieve(context.TODO(), "xxx", onUpdate)
		Expect(triggered).To(BeFalse())
		provider.UpdateConfig()
		Expect(triggered).To(BeTrue())
	})

	It("ignore updates after shutdown", func() {
		triggered := false
		onUpdate := func(*confmap.ChangeEvent) {
			triggered = true
		}

		provider := otelconfig.NewProvider("foo")
		provider.Retrieve(context.TODO(), "xxx", onUpdate)
		provider.Shutdown(context.TODO())
		provider.UpdateConfig()
		Expect(triggered).To(BeFalse())
	})

	It("handles hooks", func() {
		triggered := false
		onUpdate := func(*confmap.ChangeEvent) {
			triggered = true
		}

		provider := otelconfig.NewProvider("foo")
		provider.AddMutatingHook(func(cfg *otelconfig.Config) {
			cfg.AddReceiver("base", map[string]any{})
		})

		Expect(retrieveReceivers(provider, onUpdate)).To(ConsistOf([]string{
			"base",
		}))

		By("Adding a hook")
		Expect(triggered).To(BeFalse())
		provider.AddMutatingHook(func(cfg *otelconfig.Config) {
			// Make sure we don't rerun the same hook
			Expect(cfg.Receivers).NotTo(HaveKey("ext1"))
			cfg.AddReceiver("ext1", map[string]any{})
		})
		provider.UpdateConfig()
		Expect(triggered).To(BeTrue())
		Expect(retrieveReceivers(provider, onUpdate)).To(ConsistOf([]string{
			"base",
			"ext1",
		}))

		By("Adding a second hook")
		provider.AddMutatingHook(func(cfg *otelconfig.Config) {
			cfg.AddReceiver("ext2", map[string]any{})
		})
		Expect(retrieveReceivers(provider, onUpdate)).To(ConsistOf([]string{
			"base",
			"ext1",
			"ext2",
		}))
	})
})

func retrieveReceivers(
	provider *otelconfig.Provider,
	onUpdate confmap.WatcherFunc,
) []string {
	retrieved, err := provider.Retrieve(
		context.TODO(),
		"this providers ignores exact uri",
		onUpdate,
	)
	Expect(err).NotTo(HaveOccurred())
	raw, err := retrieved.AsRaw()
	Expect(err).NotTo(HaveOccurred())

	return maps.Keys(raw.(map[string]any)["receivers"].(map[string]any))
}
