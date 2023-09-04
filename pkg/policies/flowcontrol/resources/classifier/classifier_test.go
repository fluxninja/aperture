package classifier_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/open-policy-agent/opa/ast"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/alerts"
	"github.com/fluxninja/aperture/v2/pkg/labels"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/status"

	flowlabel "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/label"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/resources/classifier"
	. "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/resources/classifier"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/resources/classifier/compiler"
)

type object = map[string]interface{}

var classifierAttributes = &policysyncv1.ClassifierAttributes{
	PolicyName:      "test",
	PolicyHash:      "test",
	ClassifierIndex: 0,
}

var _ = Describe("Classifier", func() {
	var classifier *ClassificationEngine

	BeforeEach(func() {
		log.SetGlobalLevel(log.WarnLevel)

		alerter := alerts.NewSimpleAlerter(100)
		classifier = NewClassificationEngine(
			agentinfo.NewAgentInfo("testGroup"),
			status.NewRegistry(log.GetGlobalLogger(), alerter))
	})

	It("returns empty slice, when no rules configured", func() {
		Expect(classifier.ActiveRules()).To(BeEmpty())
	})

	Context("configured with some classification rules", func() {
		// Classifier with a simple extractor-based rule
		rs1 := &policylangv1.Classifier{
			Selectors: []*policylangv1.Selector{
				{
					ControlPoint: "ingress",
					Service:      "my-service.default.svc.cluster.local",
					AgentGroup:   "testGroup",
				},
			},
			Rules: map[string]*policylangv1.Rule{
				"foo": {
					Source:    headerExtractor("foo"),
					Telemetry: true,
				},
			},
		}

		// Classifier with Raw-rego rule, additionally gated for just "version one"
		rs2 := &policylangv1.Classifier{
			Selectors: []*policylangv1.Selector{
				{
					ControlPoint: "ingress",
					Service:      "my-service.default.svc.cluster.local",
					AgentGroup:   "testGroup",
					LabelMatcher: &policylangv1.LabelMatcher{
						MatchLabels: map[string]string{"version": "one"},
					},
				},
			},
			Rego: &policylangv1.Rego{
				Labels: map[string]*policylangv1.Rego_LabelProperties{
					"bar_twice": {
						Telemetry: true,
					},
				},
				Module: `
					package my.pkg
					bar_twice := input.attributes.request.http.headers.bar * 2
				`,
			},
		}

		// Classifier with a no service populated
		rs3 := &policylangv1.Classifier{
			Selectors: []*policylangv1.Selector{
				{
					ControlPoint: "ingress",
					Service:      "any",
					AgentGroup:   "testGroup",
				},
			},
			Rules: map[string]*policylangv1.Rule{
				"fuu": {
					Source:    headerExtractor("fuu"),
					Telemetry: true,
				},
			},
		}

		var ars1, ars2, ars3 ActiveRuleset
		BeforeEach(func() {
			var err error
			ars1, err = classifier.AddRules(context.TODO(), "one", &policysyncv1.ClassifierWrapper{
				Classifier:           rs1,
				ClassifierAttributes: classifierAttributes,
			})
			Expect(err).NotTo(HaveOccurred())
			ars2, err = classifier.AddRules(context.TODO(), "two", &policysyncv1.ClassifierWrapper{
				Classifier:           rs2,
				ClassifierAttributes: classifierAttributes,
			})
			Expect(err).NotTo(HaveOccurred())
			ars3, err = classifier.AddRules(context.TODO(), "three", &policysyncv1.ClassifierWrapper{
				Classifier:           rs3,
				ClassifierAttributes: classifierAttributes,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns active rules", func() {
			Expect(classifier.ActiveRules()).To(ConsistOf(
				compiler.ReportedRule{
					RulesetName: "one",
					LabelName:   "foo",
					Rule:        rs1.Rules["foo"],
				},
				compiler.ReportedRule{
					RulesetName: "three",
					LabelName:   "fuu",
					Rule:        rs3.Rules["fuu"],
				},
			))
		})

		It("classifies input by returning flow labels", func() {
			_, labels := classifier.Classify(
				context.TODO(),
				[]string{"my-service.default.svc.cluster.local"},
				"ingress",
				labels.PlainMap{"version": "one", "other": "tag"},
				attributesWithHeaders(object{
					"foo": "hello",
					"bar": int64(21),
				}),
			)
			Expect(labels).To(Equal(flowlabel.FlowLabels{
				"foo":       fl("hello"),
				"bar_twice": fl("42"),
			}))
		})

		It("does not classify if direction does not match", func() {
			_, labels := classifier.Classify(
				context.TODO(),
				[]string{"my-service.default.svc.cluster.local"},
				"egress",
				labels.PlainMap{"version": "one"},
				attributesWithHeaders(object{
					"foo": "hello",
					"bar": int64(21),
				}),
			)
			Expect(labels).To(BeEmpty())
		})

		It("skips rules with non-matching labels", func() {
			_, labels := classifier.Classify(
				context.TODO(),
				[]string{"my-service.default.svc.cluster.local"},
				"ingress",
				labels.PlainMap{"version": "two"},
				attributesWithHeaders(object{
					"foo": "hello",
					"bar": int64(21),
				}),
			)
			Expect(labels).To(Equal(flowlabel.FlowLabels{
				"foo": fl("hello"),
			}))
		})

		Context("when ruleset is dropped", func() {
			BeforeEach(func() { ars1.Drop() })

			It("removes subset of rules", func() {
				_, labels := classifier.Classify(
					context.TODO(),
					[]string{"my-service.default.svc.cluster.local"},
					"ingress",
					labels.PlainMap{"version": "one"},
					attributesWithHeaders(object{
						"foo": "hello",
						"bar": int64(21),
					}),
				)
				Expect(labels).To(Equal(flowlabel.FlowLabels{
					"bar_twice": fl("42"),
				}))
			})
		})

		Context("when all rulesets dropped", func() {
			BeforeEach(func() {
				ars1.Drop()
				ars2.Drop()
				ars3.Drop()
			})

			It("removes all the rules", func() {
				Expect(classifier.ActiveRules()).To(BeEmpty())
			})
		})
	})

	// helper for setting rules with a "default" selector
	setRulesForMyService := func(labelRules map[string]*policylangv1.Rule, rego *policylangv1.Rego) error {
		_, err := classifier.AddRules(context.TODO(), "test", &policysyncv1.ClassifierWrapper{
			Classifier: &policylangv1.Classifier{
				Selectors: []*policylangv1.Selector{
					{
						ControlPoint: "ingress",
						Service:      "my-service.default.svc.cluster.local",
						AgentGroup:   "testGroup",
					},
				},
				Rules: labelRules,
				Rego:  rego,
			},
			ClassifierAttributes: classifierAttributes,
		})
		return err
	}

	Context("configured classification rules with some label flags", func() {
		rules := map[string]*policylangv1.Rule{
			"foo": {
				Source:    headerExtractor("foo"),
				Telemetry: false,
			},
		}

		rego := &policylangv1.Rego{
			Labels: map[string]*policylangv1.Rego_LabelProperties{
				"bar": {
					Telemetry: true,
				},
			},
			Module: `
				package my.pkg
				bar := input.attributes.request.http.headers.bar
			`,
		}

		BeforeEach(func() {
			Expect(setRulesForMyService(rules, rego)).To(Succeed())
		})

		It("marks the returned flow labels with those flags", func() {
			_, labels := classifier.Classify(
				context.TODO(),
				[]string{"my-service.default.svc.cluster.local"},
				"ingress",
				nil,
				attributesWithHeaders(object{
					"foo": "hello",
					"bar": int64(21),
				}),
			)
			Expect(labels).To(Equal(flowlabel.FlowLabels{
				"foo": flowlabel.FlowLabelValue{Value: "hello", Telemetry: false},
				"bar": flowlabel.FlowLabelValue{Value: "21", Telemetry: true},
			}))
		})
	})

	Context("configured with same label for different rules in yaml", func() {
		// Note: we do not support multiple rules for the same label in a single
		// rulesets. But we might add support in the future, eg.:
		// "foo/1": ...
		// "foo/2": ...
		rules1 := map[string]*policylangv1.Rule{
			"foo": {
				Source:    headerExtractor("foo"),
				Telemetry: true,
			},
		}
		rules2 := map[string]*policylangv1.Rule{
			"foo": {
				Source:    headerExtractor("xyz"),
				Telemetry: true,
			},
		}

		BeforeEach(func() {
			Expect(setRulesForMyService(rules1, nil)).To(Succeed())
			Expect(setRulesForMyService(rules2, nil)).To(Succeed())
		})

		It("classifies and returns flow labels (overwrite order not specified)", func() {
			// Perhaps we can specify order by sorting rulesets? (eg. giving
			// them names from filenames)
			_, labels := classifier.Classify(
				context.TODO(),
				[]string{"my-service.default.svc.cluster.local"},
				"ingress",
				nil,
				attributesWithHeaders(object{
					"foo": "hello",
					"xyz": "cos",
					"bar": int64(21),
				}),
			)
			Expect(labels).To(SatisfyAny(
				Equal(flowlabel.FlowLabels{"foo": fl("cos")}),
				Equal(flowlabel.FlowLabels{"foo": fl("hello")}),
			))
		})
	})

	Context("configured with same label for different rules in rego", func() {
		rego1 := &policylangv1.Rego{
			Labels: map[string]*policylangv1.Rego_LabelProperties{
				"bar": {
					Telemetry: true,
				},
			},
			Module: `
				package my.pkg
				bar := input.attributes.request.http.headers.bar * 3
			`,
		}

		rego2 := &policylangv1.Rego{
			Labels: map[string]*policylangv1.Rego_LabelProperties{
				"bar": {
					Telemetry: true,
				},
			},
			Module: `
				package my.pkg
				bar := input.attributes.request.http.headers.bar * 2
			`,
		}

		BeforeEach(func() {
			Expect(setRulesForMyService(nil, rego1)).To(Succeed())
			Expect(setRulesForMyService(nil, rego2)).To(Succeed())
		})

		It("classifies and returns flow labels (overwrite order not specified)", func() {
			// Perhaps we can specify order by sorting rulesets? (eg. giving
			// them names from filenames)
			_, labels := classifier.Classify(
				context.TODO(),
				[]string{"my-service.default.svc.cluster.local"},
				"ingress",
				nil,
				attributesWithHeaders(object{
					"foo": "hello",
					"bar": int64(21),
				}),
			)
			Expect(labels).To(SatisfyAny(
				Equal(flowlabel.FlowLabels{"bar": fl("63")}),
				Equal(flowlabel.FlowLabels{"bar": fl("42")}),
			))
		})
	})

	Context("incorrect rego passed", func() {
		rego := &policylangv1.Rego{
			Labels: map[string]*policylangv1.Rego_LabelProperties{
				"bar_twice": {
					Telemetry: true,
				},
			},
			Module: `
				Package my.pkg
				bar_twice := input.attributes.request.http.headers.bar * 2
				bar_twice := input.attributes.request.http.headers.foo
			`,
		}

		It("fails to compile rego", func() {
			err := setRulesForMyService(nil, rego)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to get package name from rego module"))
		})
	})

	Context("configured with ambiguous rules in rego", func() {
		rego := &policylangv1.Rego{
			Labels: map[string]*policylangv1.Rego_LabelProperties{
				"bar": {
					Telemetry: true,
				},
			},
			Module: `
				package my.pkg
				bar = input.attributes.request.http.headers.bar * 3
				bar = input.attributes.request.http.headers.foo
			`,
		}

		BeforeEach(func() {
			Expect(setRulesForMyService(nil, rego)).To(Succeed())
		})

		It("classifies and returns empty flow labels - could not decide which rego to use", func() {
			_, labels := classifier.Classify(
				context.TODO(),
				[]string{"my-service.default.svc.cluster.local"},
				"ingress",
				nil,
				attributesWithHeaders(object{
					"foo": "hello",
					"bar": int64(21),
				}),
			)
			Expect(labels).To(Equal(flowlabel.FlowLabels{}))
		})
	})

	Context("configured with invalid label name", func() {
		// Classifier with a simple extractor-based rule
		rs := &policylangv1.Classifier{
			Selectors: []*policylangv1.Selector{
				{
					ControlPoint: "ingress",
					Service:      "my-service.default.svc.cluster.local",
					AgentGroup:   "testGroup",
				},
			},
			Rules: map[string]*policylangv1.Rule{
				"user-agent": {
					Source:    headerExtractor("foo"),
					Telemetry: true,
				},
			},
		}

		It("should reject the ruleset", func() {
			_, err := classifier.AddRules(context.TODO(), "one", &policysyncv1.ClassifierWrapper{
				Classifier:           rs,
				ClassifierAttributes: classifierAttributes,
			})
			Expect(err).To(HaveOccurred())
		})
	})
})

func fl(s string) flowlabel.FlowLabelValue {
	return flowlabel.FlowLabelValue{
		Value:     s,
		Telemetry: true,
	}
}

func attributesWithHeaders(headers object) classifier.Input {
	return newTestInput(object{
		"attributes": object{
			"request": object{
				"http": object{
					"headers": headers,
				},
			},
		},
	})
}

func headerExtractor(headerName string) *policylangv1.Rule_Extractor {
	return &policylangv1.Rule_Extractor{
		Extractor: &policylangv1.Extractor{
			Variant: &policylangv1.Extractor_From{
				From: "request.http.headers." + headerName,
			},
		},
	}
}

type testInput struct {
	iface interface{}
	value ast.Value
}

func newTestInput(input interface{}) classifier.Input {
	return testInput{
		iface: input,
		value: ast.MustInterfaceToValue(input),
	}
}

func (i testInput) Value() ast.Value       { return i.value }
func (i testInput) Interface() interface{} { return i.iface }
