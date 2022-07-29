package panic_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/FluxNinja/aperture/pkg/panic"
)

var _ = Describe("Crash-Writer", func() {
	It("flushes logs within limit, writes all the logs", func() {
		crashWriter := panic.NewCrashWriter(5)
		data := [][]byte{[]byte("log 0 "), []byte("log 1 "), []byte("log 2 "), []byte("log 3 ")}
		for _, d := range data {
			_, err := crashWriter.Write(d)
			Expect(err).NotTo(HaveOccurred())
		}

		file := bytes.Buffer{}
		crashWriter.Flush(&file)
		Expect(file.String()).To(Equal(convertToString(data)))
	})

	It("flushes logs over limit, writes last 5 logs", func() {
		crashWriter := panic.NewCrashWriter(5)
		data := [][]byte{[]byte("log 0 "), []byte("log 1 "), []byte("log 2 "), []byte("log 3 "), []byte("log 4 "), []byte("log 5 ")}
		for _, d := range data {
			_, err := crashWriter.Write(d)
			Expect(err).NotTo(HaveOccurred())
		}

		file := bytes.Buffer{}
		crashWriter.Flush(&file)
		Expect(file.String()).To(Equal(convertToString(data[1:])))
	})

	It("flushes logs over small limit, writes only the last log", func() {
		crashWriter := panic.NewCrashWriter(1)
		data := [][]byte{[]byte("log 0 "), []byte("log 1 "), []byte("log 2 "), []byte("log 3 "), []byte("log 4 "), []byte("log 5 ")}
		for _, d := range data {
			_, err := crashWriter.Write(d)
			Expect(err).NotTo(HaveOccurred())
		}

		file := bytes.Buffer{}
		crashWriter.Flush(&file)
		Expect(file.String()).To(Equal(convertToString(data[len(data)-1:])))
	})
})

func convertToString(data [][]byte) string {
	if data == nil {
		return ""
	}

	var str string
	for _, d := range data {
		str += string(d)
	}
	return str
}
