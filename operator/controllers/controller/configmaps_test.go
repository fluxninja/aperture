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

package controller

import (
	"bytes"
	_ "embed"
	"fmt"
	"path"
	"text/template"
	"time"

	controller "github.com/fluxninja/aperture/cmd/aperture-controller/config"
	. "github.com/fluxninja/aperture/operator/controllers"

	"github.com/fluxninja/aperture/operator/api/common"
	controllerv1alpha1 "github.com/fluxninja/aperture/operator/api/controller/v1alpha1"
	"github.com/fluxninja/aperture/pkg/config"
	etcd "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/net/listener"
	"github.com/fluxninja/aperture/pkg/net/tlsconfig"
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
	"github.com/fluxninja/aperture/pkg/plugins"
	"github.com/fluxninja/aperture/pkg/prometheus"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/pointer"
)

//go:embed config_test.tpl
var controllerConfigYAML string

var _ = Describe("ConfigMap for Controller", func() {
	Context("Instance without FluxNinja plugin enabled", func() {
		It("returns correct ConfigMap", func() {
			instance := &controllerv1alpha1.Controller{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AppName,
					Namespace: AppName,
				},
				Spec: controllerv1alpha1.ControllerSpec{
					ConfigSpec: controllerv1alpha1.ControllerConfigSpec{
						CommonConfigSpec: common.CommonConfigSpec{
							Server: common.ServerConfigSpec{
								Listener: listener.ListenerConfig{
									Addr: ":80",
								},
								TLS: tlsconfig.ServerTLSConfig{
									CertFile: path.Join(ControllerCertPath, ControllerCertName),
									KeyFile:  path.Join(ControllerCertPath, ControllerCertKeyName),
									Enabled:  true,
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
							Plugins: plugins.PluginsConfig{
								DisablePlugins:  false,
								DisabledPlugins: []string{"aperture-plugin-fluxninja"},
							},
							Etcd: etcd.EtcdConfig{
								Endpoints: []string{"http://agent-etcd:2379"},
								LeaseTTL:  config.MakeDuration(60 * time.Second),
							},
							Prometheus: prometheus.PrometheusConfig{
								Address: "http://aperture-prometheus-server:80",
							},
						},
						OTEL: controller.ControllerOTELConfig{
							CommonOTELConfig: otelconfig.CommonOTELConfig{
								Ports: otelconfig.PortsConfig{
									DebugPort:       8888,
									HealthCheckPort: 13133,
									PprofPort:       1777,
									ZpagesPort:      55679,
								},
							},
						},
					},
				},
			}
			config.SetDefaults(&instance.Spec.ConfigSpec)

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
					Name:      ControllerServiceName,
					Namespace: AppName,
					Labels: map[string]string{
						"app.kubernetes.io/name":       AppName,
						"app.kubernetes.io/instance":   AppName,
						"app.kubernetes.io/managed-by": OperatorName,
						"app.kubernetes.io/component":  ControllerServiceName,
					},
					Annotations: nil,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "fluxninja.com/v1alpha1",
							Name:               instance.GetName(),
							Kind:               "Controller",
							Controller:         pointer.Bool(true),
							BlockOwnerDeletion: pointer.Bool(true),
						},
					},
				},
				Data: map[string]string{
					"aperture-controller.yaml": config.String(),
				},
			}

			result, err := configMapForControllerConfig(instance.DeepCopy(), scheme.Scheme)
			Expect(err).NotTo(HaveOccurred())

			CompareComfigMap(result, expected)
		})
	})
})
