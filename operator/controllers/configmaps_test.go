/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/pointer"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
)

//go:embed agent_config_without_fluxninja_plugin_test.tpl
var agentConfigWithOutFluxNinjaPluginTest string

//go:embed agent_config_with_fluxninja_plugin_test.tpl
var agentConfigWithFluxNinjaPluginTest string

//go:embed controller_config_without_fluxninja_plugin_test.tpl
var controllerConfigWithOutFluxNinjaPluginTest string

//go:embed controller_config_with_fluxninja_plugin_test.tpl
var controllerConfigWithFluxNinjaPluginTest string

var _ = Describe("ConfigMap for Agent", func() {
	Context("Instance with FluxNinja plugin enabled", func() {
		It("returns correct ConfigMap", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						FluxNinjaPlugin: v1alpha1.FluxNinjaPluginSpec{
							Enabled:            true,
							Endpoint:           test,
							HeartbeatsInterval: "10s",
							TLS: v1alpha1.TLSSpec{
								Insecure:           true,
								InsecureSkipVerify: true,
								CAFile:             test,
							},
							APIKeySecret: v1alpha1.APIKeySecret{
								Value: test,
							},
						},
						Etcd: v1alpha1.EtcdSpec{
							Endpoints: []string{"http://agent-etcd:2379"},
							LeaseTTL:  "60s",
						},
						Prometheus: v1alpha1.PrometheusSpec{
							Address: "http://aperture-prometheus-server:80/",
						},
						ServerPort: 80,
						Log: v1alpha1.Log{
							PrettyConsole: false,
							NonBlocking:   true,
							Level:         "info",
							File:          "stderr",
						},
					},
					BatchPrerollup: v1alpha1.Batch{
						Timeout:       time.Second,
						SendBatchSize: 10000,
					},
					BatchPostrollup: v1alpha1.Batch{
						Timeout:       time.Second,
						SendBatchSize: 10000,
					},
					BatchMetricsFast: v1alpha1.Batch{
						Timeout:       time.Second,
						SendBatchSize: 10000,
					},
					DistributedCachePort: 3320,
					MemberListPort:       3322,
				},
			}

			t, err := template.New("config").Parse(agentConfigWithFluxNinjaPluginTest)
			if err != nil {
				panic(fmt.Errorf("failed to parse test config for Agent. error: '%s'", err.Error()))
			}
			var config bytes.Buffer
			if err := t.Execute(&config, struct{}{}); err != nil {
				panic(err)
			}

			expected := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      agentServiceName,
					Namespace: appName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  agentServiceName,
					},
					Annotations: nil,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Agent",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
						},
					},
				},
				Data: map[string]string{
					"aperture-agent.yaml": config.String(),
				},
			}

			result, _ := configMapForAgentConfig(instance.DeepCopy(), scheme.Scheme)
			Expect(result.Data["aperture-agent.yaml"]).To(Equal(expected.Data["aperture-agent.yaml"]))
		})
	})

	Context("Instance without FluxNinja plugin enabled", func() {
		It("returns correct ConfigMap", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					CommonSpec: v1alpha1.CommonSpec{
						ServerPort: 80,
						Log: v1alpha1.Log{
							PrettyConsole: false,
							NonBlocking:   true,
							Level:         "info",
							File:          "stderr",
						},
						Etcd: v1alpha1.EtcdSpec{
							Endpoints: []string{"http://agent-etcd:2379"},
							LeaseTTL:  "60s",
						},
						Prometheus: v1alpha1.PrometheusSpec{
							Address: "http://aperture-prometheus-server:80/",
						},
						FluxNinjaPlugin: v1alpha1.FluxNinjaPluginSpec{
							Endpoint:           test,
							HeartbeatsInterval: "10s",
						},
					},
					DistributedCachePort: 3320,
					MemberListPort:       3322,
				},
			}

			t, err := template.New("config").Parse(agentConfigWithOutFluxNinjaPluginTest)
			if err != nil {
				panic(fmt.Errorf("failed to parse test config for Agent. error: '%s'", err.Error()))
			}
			var config bytes.Buffer
			if err := t.Execute(&config, struct{}{}); err != nil {
				panic(err)
			}

			expected := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      agentServiceName,
					Namespace: appName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  agentServiceName,
					},
					Annotations: nil,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Agent",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
						},
					},
				},
				Data: map[string]string{
					"aperture-agent.yaml": config.String(),
				},
			}

			result, _ := configMapForAgentConfig(instance.DeepCopy(), scheme.Scheme)
			Expect(result.Data["aperture-agent.yaml"]).To(Equal(expected.Data["aperture-agent.yaml"]))
		})
	})
})

var _ = Describe("ConfigMap for Controller", func() {
	Context("Instance without FluxNinja plugin enabled", func() {
		It("returns correct ConfigMap", func() {
			instance := &v1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{
					CommonSpec: v1alpha1.CommonSpec{
						ServerPort: 80,
						Log: v1alpha1.Log{
							PrettyConsole: false,
							NonBlocking:   true,
							Level:         "info",
							File:          "stderr",
						},
						Etcd: v1alpha1.EtcdSpec{
							Endpoints: []string{"http://agent-etcd:2379"},
							LeaseTTL:  "60s",
						},
						Prometheus: v1alpha1.PrometheusSpec{
							Address: "http://aperture-prometheus-server:80",
						},
						FluxNinjaPlugin: v1alpha1.FluxNinjaPluginSpec{
							Endpoint:           test,
							HeartbeatsInterval: "10s",
						},
					},
				},
			}

			t, err := template.New("config").Parse(controllerConfigWithOutFluxNinjaPluginTest)
			if err != nil {
				panic(fmt.Errorf("failed to parse test config for Controller. error: '%s'", err.Error()))
			}
			var config bytes.Buffer
			if err := t.Execute(&config, struct{}{}); err != nil {
				panic(err)
			}

			expected := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      controllerServiceName,
					Namespace: appName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  controllerServiceName,
					},
					Annotations: nil,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Controller",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
						},
					},
				},
				Data: map[string]string{
					"aperture-controller.yaml": config.String(),
				},
			}

			result, _ := configMapForControllerConfig(instance.DeepCopy(), scheme.Scheme)
			Expect(result.Data).To(Equal(expected.Data))
		})
	})

	Context("Instance with FluxNinja plugin enabled", func() {
		It("returns correct ConfigMap", func() {
			instance := &v1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.ControllerSpec{
					CommonSpec: v1alpha1.CommonSpec{
						ServerPort: 80,
						Log: v1alpha1.Log{
							PrettyConsole: false,
							NonBlocking:   true,
							Level:         "info",
							File:          "stderr",
						},
						Etcd: v1alpha1.EtcdSpec{
							Endpoints: []string{"http://agent-etcd:2379"},
							LeaseTTL:  "60s",
						},
						Prometheus: v1alpha1.PrometheusSpec{
							Address: "http://aperture-prometheus-server:80",
						},
						FluxNinjaPlugin: v1alpha1.FluxNinjaPluginSpec{
							Enabled:            true,
							Endpoint:           test,
							HeartbeatsInterval: "10s",
							TLS: v1alpha1.TLSSpec{
								Insecure:           true,
								InsecureSkipVerify: true,
								CAFile:             test,
							},
							APIKeySecret: v1alpha1.APIKeySecret{
								Value: test,
							},
						},
					},
				},
			}

			t, err := template.New("config").Parse(controllerConfigWithFluxNinjaPluginTest)
			if err != nil {
				panic(fmt.Errorf("failed to parse test config for Controller. error: '%s'", err.Error()))
			}
			var config bytes.Buffer
			if err := t.Execute(&config, struct{}{}); err != nil {
				panic(err)
			}

			expected := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      controllerServiceName,
					Namespace: appName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   appName,
						"app.kubernetes.io/managed-by": operatorName,
						"app.kubernetes.io/component":  controllerServiceName,
					},
					Annotations: nil,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Controller",
							Controller:         pointer.BoolPtr(true),
							BlockOwnerDeletion: pointer.BoolPtr(true),
						},
					},
				},
				Data: map[string]string{
					"aperture-controller.yaml": config.String(),
				},
			}

			result, _ := configMapForControllerConfig(instance.DeepCopy(), scheme.Scheme)
			Expect(result.Data).To(Equal(expected.Data))
		})
	})
})

var _ = Describe("Test ConfigMap Mutate", func() {
	It("Mutate should update required fields only", func() {
		expected := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{},
			Data:       testMap,
		}

		cm := &corev1.ConfigMap{}
		err := configMapMutate(cm, expected.Data)()
		Expect(err).NotTo(HaveOccurred())
		Expect(cm).To(Equal(expected))
	})
})
