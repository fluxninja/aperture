package controlplane_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	admissionv1 "k8s.io/api/admission/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"

	policyv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/policy/v1alpha1"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane"
	"github.com/fluxninja/aperture/v2/pkg/webhooks/policyvalidator"
)

var _ = Describe("Validator", Ordered, func() {
	policySpecValidator := &controlplane.PolicySpecValidator{}
	policyValidator := policyvalidator.NewPolicyValidator([]policyvalidator.PolicySpecValidator{policySpecValidator})

	validateExample := func(contents string) {
		os.Setenv("APERTURE_CONTROLLER_NAMESPACE", "aperture-controller")
		jsonPolicy, err := yaml.YAMLToJSON([]byte(contents))
		Expect(err).ToNot(HaveOccurred())
		var policy policyv1alpha1.Policy
		err = json.Unmarshal([]byte(jsonPolicy), &policy)
		Expect(err).NotTo(HaveOccurred())
		request := &admissionv1.AdmissionRequest{
			Name:      policy.Name,
			Namespace: policy.Namespace,
			Kind:      v1.GroupVersionKind(policy.GroupVersionKind()),
			Object:    runtime.RawExtension{Raw: []byte(jsonPolicy)},
		}

		ok, msg, err := policyValidator.ValidateObject(context.TODO(), request)
		Expect(err).NotTo(HaveOccurred())
		// print msg if validation fails
		if !ok {
			fmt.Println(msg)
		}
		Expect(ok).To(BeTrue())
	}

	It("accepts example policy for demoapp", func() {
		validateExample(latencyGradientPolicy)
	})

	It("accepts example policy for rate limit", func() {
		validateExample(rateLimitPolicy)
	})

	It("accepts example policy for classification", func() {
		validateExample(classificationPolicy)
	})

	It("does not accept policy in other namespace than controller", func() {
		os.Setenv("APERTURE_CONTROLLER_NAMESPACE", "")
		jsonPolicy, err := yaml.YAMLToJSON([]byte(rateLimitPolicy))
		Expect(err).ToNot(HaveOccurred())
		var policy policyv1alpha1.Policy
		err = json.Unmarshal([]byte(jsonPolicy), &policy)
		Expect(err).NotTo(HaveOccurred())
		request := &admissionv1.AdmissionRequest{
			Name:      policy.Name,
			Namespace: policy.Namespace,
			Kind:      v1.GroupVersionKind(policy.GroupVersionKind()),
			Object:    runtime.RawExtension{Raw: []byte(jsonPolicy)},
		}

		ok, msg, err := policyValidator.ValidateObject(context.TODO(), request)
		Expect(err).NotTo(HaveOccurred())
		Expect(msg).To(Equal("Policy should be created in the same namespace as Aperture Controller"))
		Expect(ok).To(BeFalse())
	})
})

const latencyGradientPolicy = `
apiVersion: fluxninja.com/v1alpha1
kind: Policy
metadata:
  name: policies
  namespace: aperture-controller
  labels:
    fluxninja.com/validate: "true"
  name: service1-demo-app
spec:
  circuit:
    components:
    - flow_control:
        adaptive_load_scheduler:
          dry_run: false
          dry_run_config_key: load_scheduler
          in_ports:
            overload_confirmation:
              constant_signal:
                value: 1
            setpoint:
              signal_name: SETPOINT
            signal:
              signal_name: SIGNAL
          out_ports:
            desired_load_multiplier:
              signal_name: DESIRED_LOAD_MULTIPLIER
            observed_load_multiplier:
              signal_name: OBSERVED_LOAD_MULTIPLIER
          parameters:
            alerter:
              alert_name: Load Throttling Event
            gradient:
              max_gradient: 1
              min_gradient: 0.1
              slope: -1
            load_multiplier_linear_increment: 0.0025
            load_scheduler:
              scheduler:
                workloads:
                - label_matcher:
                    match_labels:
                      user_type: guest
                  parameters:
                    priority: 50
                - label_matcher:
                    match_labels:
                      http.request.header.user_type: subscriber
                  parameters:
                    priority: 200
              selectors:
              - control_point: ingress
                service: service1-demo-app.demoapp.svc.cluster.local
            max_load_multiplier: 2
    - decider:
        in_ports:
          lhs:
            signal_name: DESIRED_LOAD_MULTIPLIER
          rhs:
            constant_signal:
              value: 1
        operator: lt
        out_ports:
          output:
            signal_name: IS_CRAWLER_ESCALATION
        true_for: 30s
    - switcher:
        in_ports:
          off_signal:
            constant_signal:
              value: 10
          on_signal:
            constant_signal:
              value: 0
          switch:
            signal_name: IS_CRAWLER_ESCALATION
        out_ports:
          output:
            signal_name: RATE_LIMIT
    - flow_control:
        rate_limiter:
          in_ports:
            bucket_capacity:
              signal_name: RATE_LIMIT
            fill_amount:
              signal_name: RATE_LIMIT
          parameters:
            label_key: http.request.header.user_id
            interval: 1s
          selectors:
          - control_point: ingress
            label_matcher:
              match_labels:
                http.request.header.user_type: crawler
            service: service1-demo-app.demoapp.svc.cluster.local
    - query:
        promql:
          evaluation_interval: 10s
          out_ports:
            output:
              signal_name: SIGNAL
          query_string: sum(increase(flux_meter_sum{flow_status="OK", flux_meter_name="service1-demo-app"}[30s]))/sum(increase(flux_meter_count{flow_status="OK",
            flux_meter_name="service1-demo-app"}[30s]))
    - arithmetic_combinator:
        in_ports:
          lhs:
            signal_name: SIGNAL
          rhs:
            constant_signal:
              value: 2
        operator: mul
        out_ports:
          output:
            signal_name: MAX_EMA
    - arithmetic_combinator:
        in_ports:
          lhs:
            signal_name: SIGNAL_EMA
          rhs:
            constant_signal:
              value: 1.1
        operator: mul
        out_ports:
          output:
            signal_name: SETPOINT
    - ema:
        in_ports:
          input:
            signal_name: SIGNAL
          max_envelope:
            signal_name: MAX_EMA
        out_ports:
          output:
            signal_name: SIGNAL_EMA
        parameters:
          correction_factor_on_max_envelope_violation: 0.95
          ema_window: 1500s
          warmup_window: 60s
    evaluation_interval: 1s
  resources:
    flow_control:
      classifiers:
      - rules:
          user_type:
            extractor:
              from: request.http.headers.user-type
        selectors:
        - control_point: ingress
          service: service1-demo-app.demoapp.svc.cluster.local
      flux_meters:
        service1-demo-app:
          selectors:
          - control_point: ingress
            service: service3-demo-app.demoapp.svc.cluster.local
`

const classificationPolicy = `
apiVersion: fluxninja.com/v1alpha1
kind: Policy
metadata:
  name: policies
  namespace: aperture-controller
  labels:
    fluxninja.com/validate: "true"
spec:
  resources:
    flow_control:
      classifiers:
        - selectors:
            - control_point: ingress
              service: productpage.bookinfo.svc.cluster.local
          rego:
            labels:
              user:
                telemetry: true
            module: |
              package user_from_cookie
              cookies := split(input.attributes.request.http.headers.cookie, "; ")
              user := user {
                  cookie := cookies[_]
                  startswith(cookie, "session=")
                  session := substring(cookie, count("session="), -1)
                  parts := split(session, ".")
                  object := json.unmarshal(base64url.decode(parts[0]))
                  user := object.user
              }
          rules:
            ua:
              extractor:
                from: request.http.headers.user-agent
`

const rateLimitPolicy = `
apiVersion: fluxninja.com/v1alpha1
kind: Policy
metadata:
  name: policies
  namespace: aperture-controller
  labels:
    fluxninja.com/validate: "true"
spec:
  circuit:
    evaluation_interval: "10s"
    components:
      - variable:
          constant_output:
            value: 250.0
          out_ports:
            output:
              signal_name: "RATE_LIMIT"
      - flow_control:
          rate_limiter:
            in_ports:
              bucket_capacity:
                signal_name: "RATE_LIMIT"
              fill_amount:
                signal_name: "RATE_LIMIT"
            selectors:
              - control_point: "ingress"
                service: "service1-demo-app.demoapp.svc.cluster.local"
            parameters:
              label_key: "http.request.header.user_type"
              interval: "1s"
`
