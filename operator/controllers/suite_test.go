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
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/fluxninja/aperture/operator/api/v1alpha1"
	"github.com/fluxninja/aperture/pkg/config"
	etcd "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/prometheus"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	k8sClient                 client.Client
	k8sManager                ctrl.Manager
	testEnv                   *envtest.Environment
	ctx                       context.Context
	cancel                    context.CancelFunc
	defaultAgentInstance      *v1alpha1.Agent
	defaultControllerInstance *v1alpha1.Controller
	namespaceReconciler       *NamespaceReconciler
	mutatingWebhookReconciler *MutatingWebhookReconciler
	certDir                   = filepath.Join(".", "certs")
	test                      = "test"
	testTwo                   = "test2"
	testArray                 = []string{test}
	testArrayTwo              = []string{testTwo, test}
	testMap                   = map[string]string{
		test: test,
	}
	testMapTwo = map[string]string{
		test:    test,
		testTwo: testTwo,
	}
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	ctx, cancel = context.WithCancel(context.TODO())

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

	err = v1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = corev1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	namespaceReconciler = &NamespaceReconciler{
		Client: k8sClient,
		Scheme: k8sClient.Scheme(),
	}

	mutatingWebhookReconciler = &MutatingWebhookReconciler{
		Client:            k8sClient,
		Scheme:            k8sClient.Scheme(),
		AgentManager:      true,
		ControllerManager: true,
	}

	k8sManager, err = ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	err = os.MkdirAll(certDir, 0o777)
	Expect(err).NotTo(HaveOccurred())

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: appName,
		},
	}
	Expect(k8sClient.Create(ctx, ns)).To(BeNil())

	unmarshaller, err := config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller([]byte{})
	Expect(err).NotTo(HaveOccurred())

	defaultControllerInstance = &v1alpha1.Controller{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appName,
			Namespace: appName,
		},
		Spec: v1alpha1.ControllerSpec{
			ConfigSpec: v1alpha1.ControllerConfigSpec{
				CommonConfigSpec: v1alpha1.CommonConfigSpec{
					Etcd: etcd.EtcdConfig{
						Endpoints: []string{"10.10.10.10:1010"},
					},
					Prometheus: prometheus.PrometheusConfig{
						Address: "20.20.20.20:2020",
					},
				},
			},
			CommonSpec: v1alpha1.CommonSpec{
				LivenessProbe: v1alpha1.Probe{
					FailureThreshold: 1,
					PeriodSeconds:    1,
					SuccessThreshold: 1,
					TimeoutSeconds:   1,
				},
				ReadinessProbe: v1alpha1.Probe{
					FailureThreshold: 1,
					PeriodSeconds:    1,
					SuccessThreshold: 1,
					TimeoutSeconds:   1,
				},
				ServiceAccountSpec: v1alpha1.ServiceAccountSpec{
					Create: true,
				},
			},
			Image: v1alpha1.ControllerImage{
				Image: v1alpha1.Image{
					PullPolicy: string(corev1.PullAlways),
				},
			},
		},
	}
	unmarshaller.Unmarshal(&defaultControllerInstance.Spec)

	defaultAgentInstance = &v1alpha1.Agent{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appName,
			Namespace: appName,
		},
		Spec: v1alpha1.AgentSpec{
			ConfigSpec: v1alpha1.AgentConfigSpec{
				CommonConfigSpec: v1alpha1.CommonConfigSpec{
					Etcd: etcd.EtcdConfig{
						Endpoints: []string{"10.10.10.10:1010"},
					},
					Prometheus: prometheus.PrometheusConfig{
						Address: "20.20.20.20:2020",
					},
				},
			},
			CommonSpec: v1alpha1.CommonSpec{
				LivenessProbe: v1alpha1.Probe{
					FailureThreshold: 1,
					PeriodSeconds:    1,
					SuccessThreshold: 1,
					TimeoutSeconds:   1,
				},
				ReadinessProbe: v1alpha1.Probe{
					FailureThreshold: 1,
					PeriodSeconds:    1,
					SuccessThreshold: 1,
					TimeoutSeconds:   1,
				},
				ServiceAccountSpec: v1alpha1.ServiceAccountSpec{
					Create: true,
				},
			},
			Image: v1alpha1.AgentImage{
				Image: v1alpha1.Image{
					PullPolicy: string(corev1.PullAlways),
				},
			},
		},
	}
	unmarshaller.Unmarshal(&defaultAgentInstance.Spec)
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	os.RemoveAll(certDir)
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
