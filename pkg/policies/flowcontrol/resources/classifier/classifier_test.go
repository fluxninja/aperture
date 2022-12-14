package classifier_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/open-policy-agent/opa/ast"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/status"

	flowlabel "github.com/fluxninja/aperture/pkg/policies/flowcontrol/label"
	. "github.com/fluxninja/aperture/pkg/policies/flowcontrol/resources/classifier"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/resources/classifier/compiler"
)

type object = map[string]interface{}

var commonAttributes = &policysyncv1.CommonAttributes{
	PolicyName:     "test",
	PolicyHash:     "test",
	ComponentIndex: 0,
}

var _ = Describe("Classifier", func() {
	var classifier *ClassificationEngine

	BeforeEach(func() {
		log.SetGlobalLevel(log.WarnLevel)

		classifier = NewClassificationEngine(status.NewRegistry(log.GetGlobalLogger()))
	})

	It("returns empty slice, when no rules configured", func() {
		Expect(classifier.ActiveRules()).To(BeEmpty())
	})

	Context("configured with some classification rules", func() {
		// Classifier with a simple extractor-based rule
		rs1 := &policylangv1.Classifier{
			FlowSelector: &policylangv1.FlowSelector{
				ServiceSelector: &policylangv1.ServiceSelector{
					Service: "my-service.default.svc.cluster.local",
				},
				FlowMatcher: &policylangv1.FlowMatcher{
					ControlPoint: "ingress",
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
			FlowSelector: &policylangv1.FlowSelector{
				ServiceSelector: &policylangv1.ServiceSelector{
					Service: "my-service.default.svc.cluster.local",
				},
				FlowMatcher: &policylangv1.FlowMatcher{
					LabelMatcher: &policylangv1.LabelMatcher{
						MatchLabels: map[string]string{"version": "one"},
					},
					ControlPoint: "ingress",
				},
			},
			Rules: map[string]*policylangv1.Rule{
				"bar-twice": {
					Source: &policylangv1.Rule_Rego_{
						Rego: &policylangv1.Rule_Rego{
							Source: `
								package my.pkg
								answer := input.attributes.request.http.headers.bar * 2
							`,
							Query: "data.my.pkg.answer",
						},
					},
					Telemetry: true,
				},
			},
		}

		// Classifier with a no service populated
		rs3 := &policylangv1.Classifier{
			FlowSelector: &policylangv1.FlowSelector{
				FlowMatcher: &policylangv1.FlowMatcher{
					ControlPoint: "ingress",
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
				Classifier:       rs1,
				CommonAttributes: commonAttributes,
			})
			Expect(err).NotTo(HaveOccurred())
			ars2, err = classifier.AddRules(context.TODO(), "two", &policysyncv1.ClassifierWrapper{
				Classifier:       rs2,
				CommonAttributes: commonAttributes,
			})
			Expect(err).NotTo(HaveOccurred())
			ars3, err = classifier.AddRules(context.TODO(), "three", &policysyncv1.ClassifierWrapper{
				Classifier:       rs3,
				CommonAttributes: commonAttributes,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns active rules", func() {
			Expect(classifier.ActiveRules()).To(ConsistOf(
				compiler.ReportedRule{
					RulesetName:  "one",
					LabelName:    "foo",
					Rule:         rs1.Rules["foo"],
					FlowSelector: rs1.FlowSelector,
				},
				compiler.ReportedRule{
					RulesetName:  "two",
					LabelName:    "bar-twice",
					Rule:         rs2.Rules["bar-twice"],
					FlowSelector: rs2.FlowSelector,
				},
				compiler.ReportedRule{
					RulesetName:  "three",
					LabelName:    "fuu",
					Rule:         rs3.Rules["fuu"],
					FlowSelector: rs3.FlowSelector,
				},
			))
		})

		It("classifies input by returning flow labels", func() {
			_, labels := classifier.Classify(
				context.TODO(),
				[]string{"my-service.default.svc.cluster.local"},
				"ingress",
				map[string]string{"version": "one", "other": "tag"},
				attributesWithHeaders(object{
					"foo": "hello",
					"bar": 21,
				}),
			)
			Expect(labels).To(Equal(flowlabel.FlowLabels{
				"foo":       fl("hello"),
				"bar-twice": fl("42"),
			}))
		})

		It("doesn't classify if direction doesn't match", func() {
			_, labels := classifier.Classify(
				context.TODO(),
				[]string{"my-service.default.svc.cluster.local"},
				"egress",
				map[string]string{"version": "one"},
				attributesWithHeaders(object{
					"foo": "hello",
					"bar": 21,
				}),
			)
			Expect(labels).To(BeEmpty())
		})

		It("skips rules with non-matching labels", func() {
			_, labels := classifier.Classify(
				context.TODO(),
				[]string{"my-service.default.svc.cluster.local"},
				"ingress",
				map[string]string{"version": "two"},
				attributesWithHeaders(object{
					"foo": "hello",
					"bar": 21,
				}),
			)
			Expect(labels).To(Equal(flowlabel.FlowLabels{
				"foo": fl("hello"),
			}))
		})

		Context("when ruleset is dropped", func() {
			BeforeEach(func() { ars1.Drop() })

			It("removes removes subset of rules", func() {
				_, labels := classifier.Classify(
					context.TODO(),
					[]string{"my-service.default.svc.cluster.local"},
					"ingress",
					map[string]string{"version": "one"},
					attributesWithHeaders(object{
						"foo": "hello",
						"bar": 21,
					}),
				)
				Expect(labels).To(Equal(flowlabel.FlowLabels{
					"bar-twice": fl("42"),
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
	setRulesForMyService := func(labelRules map[string]*policylangv1.Rule) error {
		_, err := classifier.AddRules(context.TODO(), "test", &policysyncv1.ClassifierWrapper{
			Classifier: &policylangv1.Classifier{
				FlowSelector: &policylangv1.FlowSelector{
					ServiceSelector: &policylangv1.ServiceSelector{
						Service: "my-service.default.svc.cluster.local",
					},
					FlowMatcher: &policylangv1.FlowMatcher{
						ControlPoint: "ingress",
					},
				},
				Rules: labelRules,
			},
			CommonAttributes: commonAttributes,
		})
		return err
	}

	Context("configured classification rules with some label flags", func() {
		rules := map[string]*policylangv1.Rule{
			"foo": {
				Source:    headerExtractor("foo"),
				Telemetry: false,
			},
			"bar": {
				Source: &policylangv1.Rule_Rego_{
					Rego: &policylangv1.Rule_Rego{
						Source: `
							package my.pkg
							answer := input.attributes.request.http.headers.bar
							`,
						Query: "data.my.pkg.answer",
					},
				},
				Telemetry: true,
			},
		}

		BeforeEach(func() {
			Expect(setRulesForMyService(rules)).To(Succeed())
		})

		It("marks the returned flow labels with those flags", func() {
			_, labels := classifier.Classify(
				context.TODO(),
				[]string{"my-service.default.svc.cluster.local"},
				"ingress",
				nil,
				attributesWithHeaders(object{
					"foo": "hello",
					"bar": 21,
				}),
			)
			Expect(labels).To(Equal(flowlabel.FlowLabels{
				"foo": flowlabel.FlowLabelValue{Value: "hello", Telemetry: false},
				"bar": flowlabel.FlowLabelValue{Value: "21", Telemetry: true},
			}))
		})
	})

	Context("configured with same label for different rules in yaml", func() {
		// Note: we don't support multiple rules for the same label in a single
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
			Expect(setRulesForMyService(rules1)).To(Succeed())
			Expect(setRulesForMyService(rules2)).To(Succeed())
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
					"bar": 21,
				}),
			)
			Expect(labels).To(SatisfyAny(
				Equal(flowlabel.FlowLabels{"foo": fl("cos")}),
				Equal(flowlabel.FlowLabels{"foo": fl("hello")}),
			))
		})
	})

	Context("configured with same label for different rules in rego", func() {
		rules1 := map[string]*policylangv1.Rule{
			"bar": {
				Source: &policylangv1.Rule_Rego_{
					Rego: &policylangv1.Rule_Rego{
						Source: `
							package my.pkg
							answer := input.attributes.request.http.headers.bar * 3
						`,
						Query: "data.my.pkg.answer",
					},
				},
				Telemetry: true,
			},
		}
		rules2 := map[string]*policylangv1.Rule{
			"bar": {
				Source: &policylangv1.Rule_Rego_{
					Rego: &policylangv1.Rule_Rego{
						Source: `
							package my.pkg
							answer2 := input.attributes.request.http.headers.bar * 2
						`,
						Query: "data.my.pkg.answer2",
					},
				},
				Telemetry: true,
			},
		}

		BeforeEach(func() {
			Expect(setRulesForMyService(rules1)).To(Succeed())
			Expect(setRulesForMyService(rules2)).To(Succeed())
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
					"bar": 21,
				}),
			)
			Expect(labels).To(SatisfyAny(
				Equal(flowlabel.FlowLabels{"bar": fl("63")}),
				Equal(flowlabel.FlowLabels{"bar": fl("42")}),
			))
		})
	})

	Context("incorrect rego passed", func() {
		rules := map[string]*policylangv1.Rule{
			"bar-twice": {
				Source: &policylangv1.Rule_Rego_{
					Rego: &policylangv1.Rule_Rego{
						Source: `
							Package my.pkg
							bar := input.attributes.request.http.headers.bar * 2
							bar := input.attributes.request.http.headers.foo
						`,
						Query: "data.my.pkg.bar",
					},
				},
				Telemetry: true,
			},
		}

		It("fails to compile rego", func() {
			err := setRulesForMyService(rules)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(compiler.BadRego))
		})
	})

	Context("configured with ambiguous rules in rego", func() {
		rules := map[string]*policylangv1.Rule{
			"bar": {
				Source: &policylangv1.Rule_Rego_{
					Rego: &policylangv1.Rule_Rego{
						Source: `
							package my.pkg
							answer = input.attributes.request.http.headers.bar * 3
							answer = input.attributes.request.http.headers.foo
						`,
						Query: "data.my.pkg.answer",
					},
				},
				Telemetry: true,
			},
		}

		BeforeEach(func() {
			Expect(setRulesForMyService(rules)).To(Succeed())
		})

		It("classifies and returns empty flow labels - could not decide which rego to use", func() {
			_, labels := classifier.Classify(
				context.TODO(),
				[]string{"my-service.default.svc.cluster.local"},
				"ingress",
				nil,
				attributesWithHeaders(object{
					"foo": "hello",
					"bar": 21,
				}),
			)
			Expect(labels).To(Equal(flowlabel.FlowLabels{}))
		})
	})

	Context("configured with invalid label name", func() {
		// Classifier with a simple extractor-based rule
		rs := &policylangv1.Classifier{
			FlowSelector: &policylangv1.FlowSelector{
				ServiceSelector: &policylangv1.ServiceSelector{
					Service: "my-service.default.svc.cluster.local",
				},
				FlowMatcher: &policylangv1.FlowMatcher{
					ControlPoint: "ingress",
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
				Classifier:       rs,
				CommonAttributes: commonAttributes,
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

func attributesWithHeaders(headers object) ast.Value {
	return ast.MustInterfaceToValue(
		object{
			"attributes": object{
				"request": object{
					"http": object{
						"headers": headers,
					},
				},
			},
		},
	)
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
