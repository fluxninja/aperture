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
	_ "embed"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const agentConfigYAML = `dist_cache:
  bind_addr: :3320
  memberlist_bind_addr: :3322
etcd:
  endpoints:
  - http://agent-etcd:2379
  lease_ttl: 60s
log:
  file: stderr
  level: info
  non_blocking: true
otel:
  batch_postrollup:
    send_batch_size: 15000
    timeout: 1s
  batch_prerollup:
    send_batch_size: 15000
    timeout: 1s
  grpc_addr: :4317
  http_addr: :4318
plugins:
  disable_plugins: false
  disabled_plugins:
  - aperture-plugin-fluxninja
prometheus:
  address: http://aperture-prometheus-server:80
server:
  addr: :80
`

const controllerConfigYAML = `etcd:
  endpoints:
  - http://agent-etcd:2379
  lease_ttl: 60s
log:
  file: stderr
  level: info
  non_blocking: true
otel:
  batch_postrollup:
    send_batch_size: 15000
    timeout: 1s
  batch_prerollup:
    send_batch_size: 15000
    timeout: 1s
  grpc_addr: :4317
  http_addr: :4318
plugins:
  disable_plugins: false
  disabled_plugins:
  - aperture-plugin-fluxninja
prometheus:
  address: http://aperture-prometheus-server:80
server:
  addr: :80
  tls:
    certs_path: /etc/aperture/aperture-controller/certs
    enabled: true
    server_cert: crt.pem
    server_key: key.pem
`

/*var _ = Describe("ConfigMap for Agent", func() {
	Context("Instance without FluxNinja plugin enabled", func() {
		It("returns correct ConfigMap", func() {
			instance := &v1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appName,
					Namespace: appName,
				},
				Spec: v1alpha1.AgentSpec{
					ConfigSpec: v1alpha1.AgentConfigSpec{
						CommonConfigSpec: v1alpha1.CommonConfigSpec{
							Server: v1alpha1.ServerConfigSpec{
								ListenerConfig: listener.ListenerConfig{
									Addr: ":80",
								},
							},
							Log: config.LogConfig{
								PrettyConsole: false,
								NonBlocking:   true,
								LogLevel:      "info",
								LogWriterConfig: config.LogWriterConfig{
									File: "stderr",
								},
							},
							Plugins: plugins.PluginsConfig{
								DisablePlugins:  false,
								DisabledPlugins: []string{"aperture-plugin-fluxninja"},
							},
							Otel: otel.OtelConfig{
								GRPCAddr: ":4317",
								HTTPAddr: ":4318",
								BatchPrerollup: otel.BatchConfig{
									Timeout:       config.MakeDuration(1 * time.Second),
									SendBatchSize: 15000,
								},
								BatchPostrollup: otel.BatchConfig{
									Timeout:       config.MakeDuration(1 * time.Second),
									SendBatchSize: 15000,
								},
							},
							Etcd: etcd.EtcdConfig{
								Endpoints: []string{"http://agent-etcd:2379"},
								LeaseTTL:  config.MakeDuration(60 * time.Second),
							},
							Prometheus: prometheus.PrometheusConfig{
								Address: "http://aperture-prometheus-server:80",
							},
						},
						DistCache: distcache.DistCacheConfig{
							BindAddr:           ":3320",
							MemberlistBindAddr: ":3322",
						},
					},
				},
			}

			t, err := template.New("config").Parse(agentConfigYAML)
			if err != nil {
				panic(fmt.Errorf("failed to parse test config for Agent. error: '%s'", err.Error()))
			}
			var config bytes.Buffer
			if err = t.Execute(&config, struct{}{}); err != nil {
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

			result, err := configMapForAgentConfig(instance.DeepCopy(), scheme.Scheme)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Data).To(Equal(expected.Data))
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
					ConfigSpec: v1alpha1.ControllerConfigSpec{
						CommonConfigSpec: v1alpha1.CommonConfigSpec{
							Server: v1alpha1.ServerConfigSpec{
								ListenerConfig: listener.ListenerConfig{
									Addr: ":80",
								},
								TLS: tlsconfig.ServerTLSConfig{
									CertsPath:  controllerCertPath,
									ServerCert: controllerCertName,
									ServerKey:  controllerCertKeyName,
									Enabled:    true,
								},
							},
							Log: config.LogConfig{
								PrettyConsole: false,
								NonBlocking:   true,
								LogLevel:      "info",
								LogWriterConfig: config.LogWriterConfig{
									File: "stderr",
								},
							},
							Plugins: plugins.PluginsConfig{
								DisablePlugins:  false,
								DisabledPlugins: []string{"aperture-plugin-fluxninja"},
							},
							Otel: otel.OtelConfig{
								GRPCAddr: ":4317",
								HTTPAddr: ":4318",
								BatchPrerollup: otel.BatchConfig{
									Timeout:       config.MakeDuration(1 * time.Second),
									SendBatchSize: 15000,
								},
								BatchPostrollup: otel.BatchConfig{
									Timeout:       config.MakeDuration(1 * time.Second),
									SendBatchSize: 15000,
								},
							},
							Etcd: etcd.EtcdConfig{
								Endpoints: []string{"http://agent-etcd:2379"},
								LeaseTTL:  config.MakeDuration(60 * time.Second),
							},
							Prometheus: prometheus.PrometheusConfig{
								Address: "http://aperture-prometheus-server:80",
							},
						},
					},
				},
			}

			t, err := template.New("config").Parse(controllerConfigYAML)
			if err != nil {
				panic(fmt.Errorf("failed to parse test config for Controller. error: '%s'", err.Error()))
			}
			var config bytes.Buffer
			if err = t.Execute(&config, struct{}{}); err != nil {
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

			result, err := configMapForControllerConfig(instance.DeepCopy(), scheme.Scheme)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Data).To(Equal(expected.Data))
		})
	})
})*/

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
