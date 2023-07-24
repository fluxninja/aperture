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

package agent

import (
	"bytes"
	"context"
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

	agent "github.com/fluxninja/aperture/v2/cmd/aperture-agent/config"
	agentv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/agent/v1alpha1"
	"github.com/fluxninja/aperture/v2/operator/api/common"
	. "github.com/fluxninja/aperture/v2/operator/controllers"
	"github.com/fluxninja/aperture/v2/operator/controllers/testutils"
	"github.com/fluxninja/aperture/v2/pkg/config"
	distcacheconfig "github.com/fluxninja/aperture/v2/pkg/dist-cache/config"
	"github.com/fluxninja/aperture/v2/pkg/etcd"
	"github.com/fluxninja/aperture/v2/pkg/net/listener"
	otelconfig "github.com/fluxninja/aperture/v2/pkg/otelcollector/config"
	prometheus "github.com/fluxninja/aperture/v2/pkg/prometheus/config"
)

//go:embed config_test.tpl
var agentConfigYAML string

var _ = Describe("ConfigMap for Agent", func() {
	Context("Instance", func() {
		It("returns correct ConfigMap", func() {
			instance := &agentv1alpha1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: agentv1alpha1.AgentSpec{
					ConfigSpec: agentv1alpha1.AgentConfigSpec{
						CommonConfigSpec: common.CommonConfigSpec{
							Server: common.ServerConfigSpec{
								Listener: listener.ListenerConfig{
									Addr: ":80",
								},
							},
							Log: config.LogConfig{
								PrettyConsole: false,
								NonBlocking:   true,
								LogLevel:      "info",
								Writers: []config.LogWriterConfig{
									{
										File: "stderr",
									},
								},
							},
						},
						Etcd: etcd.EtcdConfig{
							Endpoints: []string{"http://agent-etcd:2379"},
							LeaseTTL:  config.MakeDuration(60 * time.Second),
						},
						Prometheus: prometheus.PrometheusConfig{
							Address: "http://aperture-prometheus-server:80",
						},
						DistCache: distcacheconfig.DistCacheConfig{
							BindAddr:           ":3320",
							MemberlistBindAddr: ":3322",
						},
						OTel: agent.AgentOTelConfig{
							CommonOTelConfig: otelconfig.CommonOTelConfig{
								Ports: otelconfig.PortsConfig{
									DebugPort:       8888,
									HealthCheckPort: 13133,
									PprofPort:       1777,
									ZpagesPort:      55679,
								},
							},
							BatchPrerollup: agent.BatchPrerollupConfig{
								Timeout:          config.MakeDuration(1 * time.Second),
								SendBatchSize:    10000,
								SendBatchMaxSize: 20000,
							},
							BatchPostrollup: agent.BatchPostrollupConfig{
								Timeout:          config.MakeDuration(1 * time.Second),
								SendBatchSize:    100,
								SendBatchMaxSize: 200,
							},
							EnableHighCardinalityPlatformMetrics: false,
						},
					},
				},
			}
			config.SetDefaults(&instance.Spec.ConfigSpec)

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
					Name:      AgentServiceName,
					Namespace: AppName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       AppName,
						"app.kubernetes.io/instance":   AppName,
						"app.kubernetes.io/managed-by": OperatorName,
						"app.kubernetes.io/component":  AgentServiceName,
					},
					Annotations: nil,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Agent",
							Controller:         pointer.Bool(true),
							BlockOwnerDeletion: pointer.Bool(true),
						},
					},
				},
				Data: map[string]string{
					"aperture-agent.yaml": config.String(),
				},
			}

			result, err := configMapForAgentConfig(context.Background(), K8sClient, instance.DeepCopy(), scheme.Scheme)
			Expect(err).NotTo(HaveOccurred())

			testutils.CompareConfigMap(result, expected)
		})
	})
})
