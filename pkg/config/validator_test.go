package config

import (
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validator", func() {
	Context("ValidateStruct", func() {
		It("is threadsafe", func() {
			// Test-suite should be run with `-race` for the test to be meaningful
			done := make(chan struct{})
			go func() {
				var s ExampleStruct
				ValidateStruct(&s)
				close(done)
			}()
			var s ExampleStruct
			ValidateStruct(&s)
			<-done
		})
	})

	Context("getValidate returns a validator", func() {
		validator := getValidate()
		It("should not be nil and it should validate an empty struct", func() {
			Expect(validator).ShouldNot(BeNil())
			err := validator.Struct(&ExampleStruct{})
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("validator custom functions", func() {
		It("it returns nil whenever an invalid data type is passed", func() {
			mockVal := "test"
			val := durationCustomTypeFunc(reflect.ValueOf(mockVal))
			Expect(val).To(BeNil())
			val = durationpbCustomTypeFunc(reflect.ValueOf(mockVal))
			Expect(val).To(BeNil())
			val = timestampCustomTypeFunc(reflect.ValueOf(mockVal))
			Expect(val).To(BeNil())
			val = timestamppbCustomTypeFunc(reflect.ValueOf(mockVal))
			Expect(val).To(BeNil())
		})
	})
})
