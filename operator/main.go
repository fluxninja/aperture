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

package main

import (
	"flag"
	"fmt"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	"k8s.io/client-go/dynamic"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/fluxninja/aperture/operator/api"
	"github.com/fluxninja/aperture/operator/controllers"
	"github.com/fluxninja/aperture/operator/controllers/agent"
	"github.com/fluxninja/aperture/operator/controllers/controller"
	"github.com/fluxninja/aperture/operator/controllers/mutatingwebhook"
	"github.com/fluxninja/aperture/operator/controllers/namespace"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(api.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var agentManager bool
	var controllerManager bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&agentManager, "agent", false,
		"Enable manager for Aperture Agent. "+
			"Enabling this will ensure that Agent Custom Resource is monitored by the Operator.")
	flag.BoolVar(&controllerManager, "controller", false,
		"Enable manager for Aperture Controller. "+
			"Enabling this will ensure that Controller Custom Resource is monitored by the Operator.")

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	if !agentManager && !controllerManager {
		setupLog.Info("One of the --agent or --controller flag is required.")
		os.Exit(1)
	}

	var leaderElectionID string

	if agentManager && controllerManager {
		leaderElectionID = "a4362587.fluxninja.com"
	} else if agentManager {
		leaderElectionID = "a4362587-agent.fluxninja.com"
	} else {
		leaderElectionID = "a4362587-controller.fluxninja.com"
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       leaderElectionID,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	dynamicClient, err := dynamic.NewForConfig(ctrl.GetConfigOrDie())
	if err != nil {
		setupLog.Error(err, "unable to create Dynamic Client")
		os.Exit(1)
	}

	var server *webhook.Server

	if agentManager || controllerManager {
		if err = controllers.CheckAndGenerateCertForOperator(); err != nil {
			setupLog.Error(err, "unable to manage webhook certificates")
			os.Exit(1)
		}

		server = mgr.GetWebhookServer()
		server.CertDir = os.Getenv("APERTURE_OPERATOR_CERT_DIR")
		server.CertName = os.Getenv("APERTURE_OPERATOR_CERT_NAME")
		server.KeyName = os.Getenv("APERTURE_OPERATOR_KEY_NAME")

		if err = (&mutatingwebhook.MutatingWebhookReconciler{
			Client:            mgr.GetClient(),
			Scheme:            mgr.GetScheme(),
			AgentManager:      agentManager,
			ControllerManager: controllerManager,
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "MutatingWebhook")
			os.Exit(1)
		}
	}

	if agentManager {
		reconciler := &agent.AgentReconciler{
			Client:        mgr.GetClient(),
			DynamicClient: dynamicClient,
			Scheme:        mgr.GetScheme(),
			Recorder:      mgr.GetEventRecorderFor("aperture-agent"),
		}

		if err = reconciler.SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Agent")
			os.Exit(1)
		}

		if err = (&namespace.NamespaceReconciler{
			Client: mgr.GetClient(),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Namespace")
			os.Exit(1)
		}

		apertureInjector := &mutatingwebhook.ApertureInjector{
			Client: mgr.GetClient(),
		}
		reconciler.ApertureInjector = apertureInjector

		server.Register(controllers.MutatingWebhookURI, &webhook.Admission{Handler: apertureInjector})
		server.Register(fmt.Sprintf("/%s", controllers.AgentMutatingWebhookURI), &webhook.Admission{Handler: &agent.AgentHooks{}})
	}

	if controllerManager {
		if err = (&controller.ControllerReconciler{
			Client:        mgr.GetClient(),
			DynamicClient: dynamicClient,
			Scheme:        mgr.GetScheme(),
			Recorder:      mgr.GetEventRecorderFor("aperture-controller"),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Controller")
			os.Exit(1)
		}

		server.Register(fmt.Sprintf("/%s", controllers.ControllerMutatingWebhookURI), &webhook.Admission{Handler: &controller.ControllerHooks{}})
	}

	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
