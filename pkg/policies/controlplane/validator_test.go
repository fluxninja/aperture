package controlplane_test

import (
	"context"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	"github.com/fluxninja/aperture/pkg/policies/controlplane"
	"github.com/fluxninja/aperture/pkg/webhooks/validation"
)

var _ = Describe("Validator", func() {
	cmFileValidator := &controlplane.CMFileValidator{}
	It("rejects non-policy configmap name", func() {
		ok := cmFileValidator.CheckCMName("foo")
		Expect(ok).To(BeFalse())
	})

	It("rejects empty policy file", func() {
		ok, _, err := cmFileValidator.ValidateFile(context.TODO(), "foo", []byte{})
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeFalse())
	})

	cmVatidator := validation.NewCMValidator([]validation.CMFileValidator{cmFileValidator})

	validateExample := func(contents string) {
		var cm corev1.ConfigMap
		err := yaml.Unmarshal([]byte(contents), &cm)
		Expect(err).NotTo(HaveOccurred())

		ok, msg, err := cmVatidator.ValidateConfigMap(context.TODO(), cm)
		Expect(err).NotTo(HaveOccurred())
		Expect(msg).To(BeEmpty())
		Expect(ok).To(BeTrue())
	}

	It("accepts example configmap for demoapp", func() {
		validateExample(latencyGradientPolicy)
	})

	It("accepts example configmap for classification", func() {
		validateExample(rateLimitPolicy)
	})

	It("accepts example configmap for classification", func() {
		validateExample(classificationPolicy)
	})
})

const latencyGradientPolicy = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: policies
  namespace: aperture-controller
  labels:
    fluxninja.com/validate: "true"
data:
  latency-gradient.yaml: |
    resources:
      flux_meters:
        "service_latency":
          selector:
            service: "service1-demo-app.demoapp.svc.cluster.local"
            control_point:
              traffic: "ingress"
      classifiers:
        - selector:
            service: service1-demo-app.demoapp.svc.cluster.local
            control_point: { traffic: ingress }
          rules:
            # An example rule using extractor.
            # See following RFC for list of available extractors and their syntax.
            ua:
              extractor:
                from: request.http.headers.user-agent
            # The same rule using raw rego. Requires specifying rego source code and a query
            also-ua:
              rego:
                source: |
                  package my.rego.pkg
                  import input.attributes.request.http
                  ua = http.headers["user-agent"]
                query: data.my.rego.pkg.ua
            user-type:
              extractor:
                from: request.http.headers.user-type
    circuit:
      evaluation_interval: "0.5s"
      components:
        - promql:
            query_string: "sum(increase(flux_meter_sum{decision_type!=\"DECISION_TYPE_REJECTED\", policy_name=\"latency-gradient\", flux_meter_name=\"service_latency\"}[5s]))/sum(increase(flux_meter_count{decision_type!=\"DECISION_TYPE_REJECTED\", policy_name=\"latency-gradient\", flux_meter_name=\"service_latency\"}[5s]))"
            evaluation_interval: "1s"
            out_ports:
              output:
                signal_name: "LATENCY"
        - constant:
            value: "2.0"
            out_ports:
              output:
                signal_name: "EMA_LIMIT_MULTIPLIER"
        - arithmetic_combinator:
            operator: "mul"
            in_ports:
              lhs:
                signal_name: "LATENCY"
              rhs:
                signal_name: "EMA_LIMIT_MULTIPLIER"
            out_ports:
              output:
                signal_name: "MAX_EMA"
        - ema:
            ema_window: "300s"
            warm_up_window: "10s"
            correction_factor_on_max_envelope_violation: "0.95"
            in_ports:
              input:
                signal_name: "LATENCY"
              max_envelope:
                signal_name: "MAX_EMA"
            out_ports:
              output:
                signal_name: "LATENCY_EMA"
        - constant:
            value: "1.1"
            out_ports:
              output:
                signal_name: "EMA_SETPOINT_MULTIPLIER"
        - arithmetic_combinator:
            operator: "mul"
            in_ports:
              lhs:
                signal_name: "LATENCY_EMA"
              rhs:
                signal_name: "EMA_SETPOINT_MULTIPLIER"
            out_ports:
              output:
                signal_name: "LATENCY_SETPOINT"
        - gradient_controller:
            slope: -1
            min_gradient: "0.1"
            max_gradient: "1.0"
            in_ports:
              signal:
                signal_name: "LATENCY"
              setpoint:
                signal_name: "LATENCY_SETPOINT"
              max:
                signal_name: "MAX_CONCURRENCY"
              control_variable:
                signal_name: "ACCEPTED_CONCURRENCY"
              optimize:
                signal_name: "CONCURRENCY_INCREMENT"
            out_ports:
              output:
                signal_name: "DESIRED_CONCURRENCY"
        - arithmetic_combinator:
            operator: "sub"
            in_ports:
              lhs:
                signal_name: "INCOMING_CONCURRENCY"
              rhs:
                signal_name: "DESIRED_CONCURRENCY"
            out_ports:
              output:
                signal_name: "DELTA_CONCURRENCY"
        - arithmetic_combinator:
            operator: "div"
            in_ports:
              lhs:
                signal_name: "DELTA_CONCURRENCY"
              rhs:
                signal_name: "INCOMING_CONCURRENCY"
            out_ports:
              output:
                signal_name: "LSF"
        - concurrency_limiter:
            scheduler:
              selector:
                service: "service1-demo-app.demoapp.svc.cluster.local"
                control_point:
                  traffic: "ingress"
              auto_tokens: true
              default_workload:
                priority: 20
              workloads:
                - workload:
                    priority: 50
                  label_matcher:
                    match_labels:
                      user-type: "guest"
                - workload:
                    priority: 200
                  label_matcher:
                    match_labels:
                      request_header_user-type: "subscriber"
              out_ports:
                accepted_concurrency:
                  signal_name: "ACCEPTED_CONCURRENCY"
                incoming_concurrency:
                  signal_name: "INCOMING_CONCURRENCY"
            load_shed_actuator:
              in_ports:
                load_shed_factor:
                  signal_name: "LSF"
        - constant:
            value: "2.0"
            out_ports:
              output:
                signal_name: "CONCURRENCY_LIMIT_MULTIPLIER"
        - constant:
            value: "10.0"
            out_ports:
              output:
                signal_name: "MIN_CONCURRENCY"
        - constant:
            value: "5.0"
            out_ports:
              output:
                signal_name: "LINEAR_CONCURRENCY_INCREMENT"
        - arithmetic_combinator:
            operator: "mul"
            in_ports:
              lhs:
                signal_name: "CONCURRENCY_LIMIT_MULTIPLIER"
              rhs:
                signal_name: "ACCEPTED_CONCURRENCY"
            out_ports:
              output:
                signal_name: "UPPER_CONCURRENCY_LIMIT"
        - max:
            in_ports:
              inputs:
                - signal_name: "UPPER_CONCURRENCY_LIMIT"
                - signal_name: "MIN_CONCURRENCY"
            out_ports:
              output:
                signal_name: "MAX_CONCURRENCY"
        - sqrt:
            scale: "0.5"
            in_ports:
              input:
                signal_name: "ACCEPTED_CONCURRENCY"
            out_ports:
              output:
                signal_name: "SQRT_CONCURRENCY_INCREMENT"
        - arithmetic_combinator:
            operator: "add"
            in_ports:
              lhs:
                signal_name: "LINEAR_CONCURRENCY_INCREMENT"
              rhs:
                signal_name: "SQRT_CONCURRENCY_INCREMENT"
            out_ports:
              output:
                signal_name: "CONCURRENCY_INCREMENT_NORMAL"
        - constant:
            value: "1.2"
            out_ports:
              output:
                signal_name: "OVERLOAD_MULTIPLIER"
        - arithmetic_combinator:
            operator: "mul"
            in_ports:
              lhs:
                signal_name: "LATENCY_EMA"
              rhs:
                signal_name: "OVERLOAD_MULTIPLIER"
            out_ports:
              output:
                signal_name: "LATENCY_OVERLOAD"
        - constant:
            value: "10.0"
            out_ports:
              output:
                signal_name: "CONCURRENCY_INCREMENT_OVERLOAD"
        - decider:
            operator: "gt"
            in_ports:
              lhs:
                signal_name: "LATENCY"
              rhs:
                signal_name: "LATENCY_OVERLOAD"
              on_true:
                signal_name: "CONCURRENCY_INCREMENT_OVERLOAD"
              on_false:
                signal_name: "CONCURRENCY_INCREMENT_NORMAL"
            out_ports:
              output:
                signal_name: "CONCURRENCY_INCREMENT"
  `

const classificationPolicy = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: policies
  namespace: aperture-controller
  labels:
    fluxninja.com/validate: "true"
data:
  productpage.yaml: |
    resources:
      classifiers:
        - selector:
            service: productpage.bookinfo.svc.cluster.local
            control_point: { traffic: ingress }
          rules:
            ua:
              extractor:
                from: request.http.headers.user-agent
            user:
              rego:
                query: data.user_from_cookie.user
                source: |
                  package user_from_cookie
                  cookies := split(input.attributes.request.http.headers.cookie, "; ")
                  cookie := cookies[_]
                  cookie.startswith("session=")
                  session := substring(cookie, count("session="), -1)
                  parts := split(session, ".")
                  object := json.unmarshal(base64url.decode(parts[0]))
                  user := object.user
`

const rateLimitPolicy = `
  rate-limit-policy.yaml: |
    circuit:
      evaluation_interval: "0.5s"
      components:
        - constant:
            value: "250.0"
            out_ports:
              output:
                signal_name: "RATE_LIMIT"
        - rate_limiter:
            in_ports:
              limit:
                signal_name: "RATE_LIMIT"
            selector:
              service: "service1-demo-app.demoapp.svc.cluster.local"
              control_point:
                traffic: "ingress"
            label_key: "request_header_user-type"
            limit_reset_interval: "1s"
`
