package extractors_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/lithammer/dedent"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	classificationv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/resources/classifier/extractors"
	"github.com/fluxninja/aperture/v2/pkg/utils"
)

func TestExtractors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Extractors Suite")
}

type LabelExtractors struct {
	Labels map[string]*classificationv1.Extractor `json:"labels"`
}

func checkOk(yamlString string, expectedRego string) {
	// Note: The map[string]Extractor format does not occur in policy, but is
	// helpful to test multiple extractors in a single test
	var labelExtractors LabelExtractors
	yamlString = dedent.Dedent(yamlString)
	err := config.UnmarshalYAML([]byte(yamlString), &labelExtractors)
	Expect(err).ToNot(HaveOccurred())

	rego, err := extractors.CompileToRego("pkgname", labelExtractors.Labels)
	Expect(err).NotTo(HaveOccurred())
	regoMeat := strings.TrimPrefix(rego, "package pkgname")
	expectedRego = dedent.Dedent(expectedRego)
	Expect(regoMeat).To(Equal(expectedRego))

	// Also, check if survives serialization roundtrip (we have some custom
	// marshal/unmarshal so it is worth checking)
	jsonBytes, err := json.Marshal(labelExtractors)
	Expect(err).NotTo(HaveOccurred())
	var labelExtractors2 LabelExtractors
	Expect(json.Unmarshal(jsonBytes, &labelExtractors2)).To(Succeed())
	Expect(labelExtractors2).To(Equal(labelExtractors))
}

var l *utils.GoLeakDetector

var _ = BeforeSuite(func() {
	l = utils.NewGoLeakDetector()
})

var _ = AfterSuite(func() {
	err := l.FindLeaks()
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("Extractor", func() {
	It("parses and compiles simple extractors", func() {
		checkOk(
			`
      labels:
        method:
          from: request.http.method

        path:
          from: request.http.path

        protocol:
          from: request.http.protocol
			`,
			`
			method := input.attributes.request.http.method
			path := input.attributes.request.http.path
			protocol := input.attributes.request.http.protocol
			`,
		)
	})

	It("parses and compiles advanced extractors", func() {
		checkOk(
			// Note: use spaces inside yaml!
			`
      labels:
        source:
          address:
            from: source.address

        destination:
          address:
            from: destination.address

        endpoint:
          path_templates:
            template_values:
              /poultry/*: animals
              /pets/*: animals
              /cart-db/*: cart
              /cart-ui/*: cart

        query_type:
          json:
            from: request.http.body
            pointer: /query

        user_agent:
          from: request.http.headers.user-agent

        user:
          jwt:
            from: request.http.bearer
            json_pointer: /name
			`,
			`
			destination = result {
			  value := input.attributes.destination.address
			  result := concat(":", [value.socketAddress.address, format_int(value.socketAddress.portValue, 10)])
			}
			destination = input.attributes.destination.address.pipe.path
			endpoint = "cart" {
			  count(_ninja_segments) >= 1
			  _ninja_segments[0] == "cart-db"
			} else = "cart" {
			  count(_ninja_segments) >= 1
			  _ninja_segments[0] == "cart-ui"
			} else = "animals" {
			  count(_ninja_segments) >= 1
			  _ninja_segments[0] == "pets"
			} else = "animals" {
			  count(_ninja_segments) >= 1
			  _ninja_segments[0] == "poultry"
			}
			query_type := json.unmarshal(input.attributes.request.http.body).query
			source = result {
			  value := input.attributes.source.address
			  result := concat(":", [value.socketAddress.address, format_int(value.socketAddress.portValue, 10)])
			}
			source = input.attributes.source.address.pipe.path
			user := payload.name {
			  [_, payload, _] := io.jwt.decode(_ninja_bearer)
			}
			user_agent := input.attributes.request.http.headers["user-agent"]

			_ninja_components := split(input.attributes.request.http.path, "?")
			_ninja_path := _ninja_components[0]
			_ninja_segments := split(trim(_ninja_path, "/"), "/")

			_ninja_bearer := value {
				header := input.attributes.request.http.headers.authorization
				startswith(header, "Bearer ")
				value := substring(header, count("Bearer "), -1)
			}
			`,
		)
	})

	It("handles edge cases of json pointer", func() {
		checkOk(
			`
      labels:
        foo:
          json:
            from: request.http.body
            pointer: /foo/-bar-/~1etc/2
			`,
			`
			foo := json.unmarshal(input.attributes.request.http.body).foo["-bar-"]["/etc"][{"2", 2}[_]]
			`,
		)
	})

	Context("path templates extractor", func() {
		It("parses and compiles", func() {
			checkOk(
				`
        labels:
          endpoint:
            path_templates:
              template_values:
                /users/{userId}: users
                /register: register
                /static/*: static
                /*: other
				`,
				`
				endpoint = "users" {
				  count(_ninja_segments) == 2
				  _ninja_segments[0] == "users"
				} else = "register" {
				  count(_ninja_segments) == 1
				  _ninja_segments[0] == "register"
				} else = "static" {
				  count(_ninja_segments) >= 1
				  _ninja_segments[0] == "static"
				} else = "other" {
				  count(_ninja_segments) >= 0
				}

				_ninja_components := split(input.attributes.request.http.path, "?")
				_ninja_path := _ninja_components[0]
				_ninja_segments := split(trim(_ninja_path, "/"), "/")
				`,
			)
		})

		It("orders matches from most to least specific", func() {
			checkOk(
				`
        labels:
          endpoint:
            path_templates:
              template_values:
                /foo/bar/{}/{}: a
                /foo/bar/{}/{}/*: b
                /foo/bar/{}/*: c
                /foo/bar/*: d
                /foo: e
                /foo/*: f
				`,
				`
				endpoint = "a" {
				  count(_ninja_segments) == 4
				  _ninja_segments[0] == "foo"
				  _ninja_segments[1] == "bar"
				} else = "b" {
				  count(_ninja_segments) >= 4
				  _ninja_segments[0] == "foo"
				  _ninja_segments[1] == "bar"
				} else = "c" {
				  count(_ninja_segments) >= 3
				  _ninja_segments[0] == "foo"
				  _ninja_segments[1] == "bar"
				} else = "d" {
				  count(_ninja_segments) >= 2
				  _ninja_segments[0] == "foo"
				  _ninja_segments[1] == "bar"
				} else = "e" {
				  count(_ninja_segments) == 1
				  _ninja_segments[0] == "foo"
				} else = "f" {
				  count(_ninja_segments) >= 1
				  _ninja_segments[0] == "foo"
				}

				_ninja_components := split(input.attributes.request.http.path, "?")
				_ninja_path := _ninja_components[0]
				_ninja_segments := split(trim(_ninja_path, "/"), "/")
				`,
			)
		})
	})
})
