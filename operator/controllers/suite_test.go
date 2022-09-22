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
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/fluxninja/aperture/operator/api"
	agentv1alpha1 "github.com/fluxninja/aperture/operator/api/agent/v1alpha1"
	"github.com/fluxninja/aperture/operator/api/common"
	controllerv1alpha1 "github.com/fluxninja/aperture/operator/api/controller/v1alpha1"
	policyv1alpha1 "github.com/fluxninja/aperture/operator/api/policy/v1alpha1"
	"github.com/fluxninja/aperture/pkg/config"
	etcd "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/prometheus"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	testEnv *envtest.Environment
	cancel  context.CancelFunc
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	Ctx, cancel = context.WithCancel(context.TODO())

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
		CRDInstallOptions: envtest.CRDInstallOptions{
			MaxTime: 60 * time.Second,
		},
	}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = api.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = corev1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	K8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(K8sClient).NotTo(BeNil())

	K8sDynamicClient, err = dynamic.NewForConfig(cfg)
	Expect(err).NotTo(HaveOccurred())
	Expect(K8sDynamicClient).NotTo(BeNil())

	K8sManager, err = ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	err = os.MkdirAll(CertDir, 0o777)
	Expect(err).NotTo(HaveOccurred())

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: AppName,
		},
	}
	Expect(K8sClient.Create(Ctx, ns)).To(BeNil())

	DefaultControllerInstance = &controllerv1alpha1.Controller{
		ObjectMeta: metav1.ObjectMeta{
			Name:      AppName,
			Namespace: AppName,
		},
		Spec: controllerv1alpha1.ControllerSpec{
			ConfigSpec: controllerv1alpha1.ControllerConfigSpec{
				CommonConfigSpec: common.CommonConfigSpec{
					Etcd: etcd.EtcdConfig{
						Endpoints: []string{"10.10.10.10:1010"},
					},
					Prometheus: prometheus.PrometheusConfig{
						Address: "20.20.20.20:2020",
					},
				},
			},
			CommonSpec: common.CommonSpec{
				LivenessProbe: common.Probe{
					FailureThreshold: 1,
					PeriodSeconds:    1,
					SuccessThreshold: 1,
					TimeoutSeconds:   1,
				},
				ReadinessProbe: common.Probe{
					FailureThreshold: 1,
					PeriodSeconds:    1,
					SuccessThreshold: 1,
					TimeoutSeconds:   1,
				},
				ServiceAccountSpec: common.ServiceAccountSpec{
					Create: true,
				},
			},
			Image: common.ControllerImage{
				Image: common.Image{
					PullPolicy: string(corev1.PullAlways),
				},
			},
		},
	}

	err = config.UnmarshalYAML([]byte{}, &DefaultControllerInstance.Spec)
	Expect(err).NotTo(HaveOccurred())

	DefaultAgentInstance = &agentv1alpha1.Agent{
		ObjectMeta: metav1.ObjectMeta{
			Name:      AppName,
			Namespace: AppName,
		},
		Spec: agentv1alpha1.AgentSpec{
			ConfigSpec: agentv1alpha1.AgentConfigSpec{
				CommonConfigSpec: common.CommonConfigSpec{
					Etcd: etcd.EtcdConfig{
						Endpoints: []string{"10.10.10.10:1010"},
					},
					Prometheus: prometheus.PrometheusConfig{
						Address: "20.20.20.20:2020",
					},
				},
			},
			CommonSpec: common.CommonSpec{
				LivenessProbe: common.Probe{
					FailureThreshold: 1,
					PeriodSeconds:    1,
					SuccessThreshold: 1,
					TimeoutSeconds:   1,
				},
				ReadinessProbe: common.Probe{
					FailureThreshold: 1,
					PeriodSeconds:    1,
					SuccessThreshold: 1,
					TimeoutSeconds:   1,
				},
				ServiceAccountSpec: common.ServiceAccountSpec{
					Create: true,
				},
			},
			Image: common.AgentImage{
				Image: common.Image{
					PullPolicy: string(corev1.PullAlways),
				},
			},
		},
	}
	err = config.UnmarshalYAML([]byte{}, &DefaultAgentInstance.Spec)
	Expect(err).NotTo(HaveOccurred())

	DefaultPolicyInstance = &policyv1alpha1.Policy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      AppName,
			Namespace: AppName,
		},
		Spec: runtime.RawExtension{
			Raw: []byte("{\"circuit\":{\"components\":[{\"constant\":{\"out_ports\":{\"output\":{\"signal_name\":\"EMA_LIMIT_MULTIPLIER\"}}}}]},\"resources\":{\"flux_meters\":{\"service1-demo-app\":{\"selector\":{\"agent_group\":\"default\",\"control_point\":{\"traffic\":\"ingress\"},\"service\":\"service1-demo-app.demoapp.svc.cluster.local\"}}},\"classifiers\":[{\"rules\":{\"user_type\":{\"extractor\":{\"from\":\"request.http.headers.user_type\"}}},\"selector\":{\"agent_group\":\"default\",\"control_point\":{\"traffic\":\"ingress\"},\"service\":\"service1-demo-app.demoapp.svc.cluster.local\"}}]}}"),
		},
	}

	PolicyFilePath = PoliciesDir
	err = os.MkdirAll(PolicyFilePath, 0o777)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	os.RemoveAll(CertDir)
	os.RemoveAll(PoliciesDir)
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
