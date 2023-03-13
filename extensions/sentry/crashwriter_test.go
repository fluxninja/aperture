package sentry

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Crash-Writer", func() {
	It("gets all the buffered logs", func() {
		crashWriter := newCrashWriter(5)
		data := [][]byte{[]byte("{\"log1\": \"test1\"}"), []byte("{\"log2\": \"test2\"}"), []byte("{\"log3\": \"test3\"}"), []byte("{\"log4\": \"test4\"}")}

		for _, d := range data {
			_, err := crashWriter.Write(d)
			Expect(err).NotTo(HaveOccurred())
		}

		logs := crashWriter.GetCrashLogs()
		Expect(logs).To(Equal(convertToCrashLogs(data)))
	})

	It("gets last 5 buffered logs over limit", func() {
		crashWriter := newCrashWriter(5)
		data := [][]byte{[]byte("{\"log1\": \"test1\"}"), []byte("{\"log2\": \"test2\"}"), []byte("{\"log3\": \"test3\"}"), []byte("{\"log4\": \"test4\"}"), []byte("{\"log5\": \"test5\"}"), []byte("{\"log6\": \"test6\"}")}

		for _, d := range data {
			_, err := crashWriter.Write(d)
			Expect(err).NotTo(HaveOccurred())
		}

		logs := crashWriter.GetCrashLogs()
		Expect(logs).To(Equal(convertToCrashLogs(data[1:])))
	})

	It("gets last buffered log over limit size 1", func() {
		crashWriter := newCrashWriter(1)
		data := [][]byte{[]byte("{\"log1\": \"test1\"}"), []byte("{\"log2\": \"test2\"}"), []byte("{\"log3\": \"test3\"}"), []byte("{\"log4\": \"test4\"}"), []byte("{\"log5\": \"test5\"}"), []byte("{\"log6\": \"test6\"}")}

		for _, d := range data {
			_, err := crashWriter.Write(d)
			Expect(err).NotTo(HaveOccurred())
		}

		logs := crashWriter.GetCrashLogs()
		Expect(logs).To(Equal(convertToCrashLogs(data[len(data)-1:])))
	})

	It("flushes and drains the crash writer buffer", func() {
		crashWriter := newCrashWriter(5)
		data := [][]byte{[]byte("{\"log1\": \"test1\"}"), []byte("{\"log2\": \"test2\"}"), []byte("{\"log3\": \"test3\"}"), []byte("{\"log4\": \"test4\"}")}

		for _, d := range data {
			_, err := crashWriter.Write(d)
			Expect(err).NotTo(HaveOccurred())
		}

		crashWriter.Flush()
		logs := crashWriter.GetCrashLogs()
		Expect(logs).To(BeEmpty())
	})
})

func convertToCrashLogs(data [][]byte) []map[string]interface{} {
	var logs []map[string]interface{}

	for _, d := range data {
		log := make(map[string]interface{})
		_ = json.Unmarshal(d, &log)
		logs = append(logs, log)
	}

	return logs
}
