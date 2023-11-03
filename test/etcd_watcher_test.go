package test

import (
	"context"
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	etcdnotifier "github.com/fluxninja/aperture/v2/pkg/etcd/notifier"
	filesystem "github.com/fluxninja/aperture/v2/pkg/file-system"
	fsnotifier "github.com/fluxninja/aperture/v2/pkg/file-system/notifier"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
)

var _ = Describe("Etcd Watcher", func() {
	ctx := context.Background()

	When("Adding etcd key notifier", func() {
		var etcdKeyNotifier *etcdnotifier.KeyToEtcdNotifier
		var err error
		const notifierPrefix = "fuu"
		const notifierKey = "baz"
		notifierPath := fmt.Sprintf("%s/%s", notifierPrefix, notifierKey)
		trackedPath := fmt.Sprintf("foo/%s", notifierKey)
		const val = "val"

		BeforeEach(func() {
			k := notifiers.Key(notifierKey)
			etcdKeyNotifier, err = etcdnotifier.NewKeyToEtcdNotifier(k, notifierPrefix, etcdClient, false)
			Expect(err).NotTo(HaveOccurred())
			etcdKeyNotifier.Start()
			err := etcdWatcher.AddKeyNotifier(etcdKeyNotifier)
			Expect(err).NotTo(HaveOccurred())
			time.Sleep(1 * time.Second)
		})

		AfterEach(func() {
			err := etcdWatcher.RemoveKeyNotifier(etcdKeyNotifier)
			time.Sleep(1 * time.Second)
			etcdKeyNotifier.Stop()
			etcdKeyNotifier = nil
			Expect(err).NotTo(HaveOccurred())
		})

		It("tracks etcd key even if it does not exist at first", func() {
			_, err := etcdClient.Get(ctx, notifierPath)
			Expect(err).ToNot(HaveOccurred())

			val := "val"
			etcdClient.Put(trackedPath, val)
			time.Sleep(1 * time.Second)

			resp, err := etcdClient.Get(ctx, notifierPath)
			Expect(err).ToNot(HaveOccurred())

			found := false
			for _, kv := range resp.Kvs {
				if string(kv.Value) == val {
					found = true
				}
			}
			Expect(found).To(BeTrue())
		})

		It("deletes key notifier correctly", func() {
			etcdClient.Put(trackedPath, val)
			time.Sleep(1 * time.Second)

			resp, err := etcdClient.Get(ctx, notifierPath)
			Expect(err).ToNot(HaveOccurred())

			found := false
			for _, kv := range resp.Kvs {
				if string(kv.Value) == val {
					found = true
				}
			}
			Expect(found).To(BeTrue())

			etcdClient.Delete(trackedPath)
			time.Sleep(1 * time.Second)

			resp, err = etcdClient.Get(ctx, notifierPath)
			Expect(err).ToNot(HaveOccurred())

			found = false
			for _, kv := range resp.Kvs {
				if string(kv.Key) == notifierPath {
					found = true
				}
			}
			Expect(found).To(BeFalse())
		})

		It("deletes key when key notifier is removed", func() {
			etcdClient.Put(trackedPath, val)
			time.Sleep(1 * time.Second)

			resp, err := etcdClient.Get(ctx, notifierPath)
			Expect(err).ToNot(HaveOccurred())

			found := false
			for _, kv := range resp.Kvs {
				if string(kv.Value) == val {
					found = true
				}
			}
			Expect(found).To(BeTrue())

			err = etcdWatcher.RemoveKeyNotifier(etcdKeyNotifier)
			Expect(err).NotTo(HaveOccurred())

			newVal := "newVal"
			etcdClient.Put(trackedPath, newVal)
			time.Sleep(1 * time.Second)

			resp, err = etcdClient.Get(ctx, notifierPath)
			Expect(err).ToNot(HaveOccurred())

			found = false
			for _, kv := range resp.Kvs {
				if string(kv.Value) == val {
					found = true
				}
			}
			Expect(found).To(BeFalse())

			resp, err = etcdClient.Get(ctx, trackedPath)
			Expect(err).ToNot(HaveOccurred())

			found = false
			for _, kv := range resp.Kvs {
				if string(kv.Value) == newVal {
					found = true
				}
			}
			Expect(found).To(BeTrue())
		})
	})

	When("Adding etcd prefix notifier", func() {
		var etcdNotifier *etcdnotifier.PrefixToEtcdNotifier

		BeforeEach(func() {
			etcdClient.Put("foo/key1", "val1")
			etcdClient.Put("foo/key2", "val2")
			etcdClient.Put("foo/key3", "val3")
			etcdNotifier = etcdnotifier.NewPrefixToEtcdNotifier("bar/", etcdClient)
			etcdNotifier.Start()
			err := etcdWatcher.AddPrefixNotifier(etcdNotifier)
			Expect(err).ToNot(HaveOccurred())
			time.Sleep(1 * time.Second)
		})

		AfterEach(func() {
			err := etcdWatcher.RemovePrefixNotifier(etcdNotifier)
			time.Sleep(1 * time.Second)
			etcdNotifier.Stop()
			etcdNotifier = nil
			Expect(err).ToNot(HaveOccurred())
		})

		It("tracks the prefix notifier key correctly", func() {
			val1 := "val1"

			resp, err := etcdClient.Get(ctx, "bar/key1")
			Expect(err).ToNot(HaveOccurred())

			found := false
			for _, kv := range resp.Kvs {
				if string(kv.Value) == val1 {
					found = true
				}
			}
			Expect(found).To(BeTrue())
		})

		It("writes new prefix notifier key correctly", func() {
			key4 := "key4"
			fookey4 := fmt.Sprintf("foo/%s", key4)
			barkey4 := fmt.Sprintf("bar/%s", key4)
			val4 := "val4"

			etcdClient.Put(fookey4, val4)

			time.Sleep(1 * time.Second)

			resp, err := etcdClient.Get(ctx, barkey4)
			Expect(err).ToNot(HaveOccurred())

			found := false
			for _, kv := range resp.Kvs {
				if string(kv.Value) == val4 {
					found = true
				}
			}
			Expect(found).To(BeTrue())
		})

		It("updates the prefix notifier key correctly", func() {
			newVal := "v1"
			etcdClient.Put("foo/key2", newVal)

			time.Sleep(1 * time.Second)

			resp, err := etcdClient.Get(ctx, "bar/key2")
			Expect(err).ToNot(HaveOccurred())

			found := false
			for _, kv := range resp.Kvs {
				if string(kv.Value) == newVal {
					found = true
				}
			}
			Expect(found).To(BeTrue())
		})

		It("deletes the prefix notifier key correctly", func() {
			etcdClient.Delete("foo/key3")

			resp, err := etcdClient.Get(ctx, "bar/")
			Expect(err).ToNot(HaveOccurred())

			found := false
			for _, kv := range resp.Kvs {
				if string(kv.Key) == "bar/key3" {
					found = true
				}
			}
			Expect(found).To(BeFalse())
		})
	})

	When("Adding fs prefix notifier", func() {
		var fsNotifier *fsnotifier.PrefixToFSNotifier
		var tempDir string
		var err error

		BeforeEach(func() {
			tempDir, err = os.MkdirTemp("", "etcdwatcher")
			Expect(err).ToNot(HaveOccurred())

			fsNotifier = fsnotifier.NewPrefixToFSNotifier(tempDir, ".yaml")
			fsNotifier.Start()
			err = etcdWatcher.AddPrefixNotifier(fsNotifier)
			Expect(err).ToNot(HaveOccurred())
			time.Sleep(1 * time.Second)
		})

		AfterEach(func() {
			err = etcdWatcher.RemovePrefixNotifier(fsNotifier)
			time.Sleep(1 * time.Second)
			fsNotifier.Stop()
			fsNotifier = nil
			Expect(err).ToNot(HaveOccurred())

			err = os.RemoveAll(tempDir)
			Expect(err).ToNot(HaveOccurred())
		})

		It("tracks the prefix notifier key correctly", func() {
			val1 := "val1"
			fileInfo := filesystem.NewFileInfo(tempDir, "key1", ".yaml")

			b, err := fileInfo.ReadAsByteBufferFromFile()
			Expect(err).NotTo(HaveOccurred())

			Expect(string(b)).To(Equal(val1))
		})

		It("writes new prefix notifier key correctly", func() {
			fileInfo := filesystem.NewFileInfo(tempDir, "key4", ".yaml")
			val4 := "val4"

			etcdClient.Put("foo/key4", val4)

			time.Sleep(1 * time.Second)

			b, err := fileInfo.ReadAsByteBufferFromFile()
			Expect(err).NotTo(HaveOccurred())

			Expect(string(b)).To(Equal(val4))
		})

		It("updates the prefix notifier key correctly", func() {
			newVal := "v1"
			etcdClient.Put("foo/key2", newVal)

			time.Sleep(1 * time.Second)

			fileInfo := filesystem.NewFileInfo(tempDir, "key2", ".yaml")

			b, err := fileInfo.ReadAsByteBufferFromFile()
			Expect(err).NotTo(HaveOccurred())

			Expect(string(b)).To(Equal(newVal))
		})

		It("deletes the prefix notifier key correctly", func() {
			etcdClient.Delete("foo/key3")

			fileInfo := filesystem.NewFileInfo(tempDir, "key3", ".yaml")

			exists, err := fileInfo.ExistsFile()
			Expect(err).NotTo(HaveOccurred())

			Expect(exists).To(BeFalse())
		})
	})
})
