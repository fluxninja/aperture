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
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	agentv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/agent/v1alpha1"
	controllerv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/controller/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	apimachineryversion "k8s.io/apimachinery/pkg/util/version"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/fluxninja/aperture/v2/operator/api"
	"github.com/fluxninja/aperture/v2/operator/controllers"
	"github.com/fluxninja/aperture/v2/operator/controllers/agent"
	"github.com/fluxninja/aperture/v2/operator/controllers/controller"
	"github.com/fluxninja/aperture/v2/operator/controllers/mutatingwebhook"
	"github.com/fluxninja/aperture/v2/operator/controllers/namespace"
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
	var multipleControllersEnabled bool
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
	flag.BoolVar(&multipleControllersEnabled, "experimental-multiple-controllers", false,
		"Experimental support for deployment of multiple controllers.")

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	if multipleControllersEnabled {
		setupLog.Info("Experimental support for managing multiple controllers enabled.")
	}

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

	server := webhook.NewServer(webhook.Options{
		CertDir:  os.Getenv("APERTURE_OPERATOR_CERT_DIR"),
		CertName: os.Getenv("APERTURE_OPERATOR_CERT_NAME"),
		KeyName:  os.Getenv("APERTURE_OPERATOR_KEY_NAME"),
	})

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       leaderElectionID,
		WebhookServer:          server,
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

	// Creating the discovery client
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(ctrl.GetConfigOrDie())
	if err != nil {
		setupLog.Error(err, "unable to create Discovery Client")
		os.Exit(1)
	}

	// Querying the local kubernetes version
	serverVersion, err := discoveryClient.ServerVersion()
	if err != nil {
		setupLog.Error(err, "unable to get the local kubernetes version")
		os.Exit(1)
	}
	controllers.CurrentKubernetesVersion = apimachineryversion.MustParseSemantic(serverVersion.String())

	// Checking if the minimum Kubernetes version is satisfied
	controllers.MinimumKubernetesVersionBool = controllers.CurrentKubernetesVersion.AtLeast(apimachineryversion.MustParseSemantic(controllers.MinimumKubernetesVersion))

	if agentManager || controllerManager {
		if err = controllers.CheckAndGenerateCertForOperator(ctrl.GetConfigOrDie()); err != nil {
			setupLog.Error(err, "unable to manage webhook certificates")
			os.Exit(1)
		}

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
	var agentReconciler *agent.AgentReconciler
	if agentManager {
		agentReconciler = &agent.AgentReconciler{
			Client:        mgr.GetClient(),
			DynamicClient: dynamicClient,
			Scheme:        mgr.GetScheme(),
			Recorder:      mgr.GetEventRecorderFor("aperture-agent"),
		}

		if err = agentReconciler.SetupWithManager(mgr); err != nil {
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
			Client:  mgr.GetClient(),
			Decoder: admission.NewDecoder(mgr.GetScheme()),
		}
		agentReconciler.ApertureInjector = apertureInjector

		server.Register(controllers.MutatingWebhookURI, &webhook.Admission{Handler: apertureInjector})
		server.Register(fmt.Sprintf("/%s", controllers.AgentMutatingWebhookURI), &webhook.Admission{Handler: &agent.AgentHooks{}})
	}

	var controllerReconciler *controller.ControllerReconciler
	if controllerManager {
		controllerReconciler = &controller.ControllerReconciler{
			Client:              mgr.GetClient(),
			DynamicClient:       dynamicClient,
			Scheme:              mgr.GetScheme(),
			Recorder:            mgr.GetEventRecorderFor("aperture-controller"),
			MultipleControllers: multipleControllersEnabled,
			ResourcesDeleted:    map[string]bool{},
		}

		if err = controllerReconciler.SetupWithManager(mgr); err != nil {
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

	ctx := setupContext(agentReconciler, controllerReconciler, agentManager)
	setupLog.Info("starting webhook server")
	go func() {
		if err := server.Start(ctx); err != nil {
			setupLog.Error(err, "unable to run webhook server")
			os.Exit(1)
		}
	}()

	setupLog.Info("starting manager")
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func setupContext(agentReconciler *agent.AgentReconciler, controllerReconciler *controller.ControllerReconciler, agentManager bool) context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		removeFinalizer(ctx, agentReconciler, controllerReconciler, agentManager)
		cancel()
	}()
	return ctx
}

func removeFinalizer(ctx context.Context, agentReconciler *agent.AgentReconciler, controllerReconciler *controller.ControllerReconciler, agentManager bool) {
	if agentManager {
		agentList := &agentv1alpha1.AgentList{}
		err := agentReconciler.Client.List(ctx, agentList)
		if err != nil {
			setupLog.Error(err, "Error while getting the agent")
		}
		if agentList.Items != nil && len(agentList.Items) != 0 {
			for _, agentCR := range agentList.Items {
				agentCR := agentCR
				if controllerutil.ContainsFinalizer(&agentCR, controllers.FinalizerName) {
					setupLog.Info("Operator is getting deleted. Removing finalizer from the Agent CR")
					controllerutil.RemoveFinalizer(&agentCR, controllers.FinalizerName)
					if err = agentReconciler.UpdateAgent(ctx, &agentCR); err != nil && !errors.IsNotFound(err) {
						setupLog.Error(err, "Error while removing Finalizer from the agent")
					}
				}
			}
		}
	} else {
		controllerList := &controllerv1alpha1.ControllerList{}

		err := controllerReconciler.Client.List(ctx, controllerList)
		if err != nil {
			setupLog.Error(err, "Error while getting the controller")
		}
		if controllerList.Items != nil && len(controllerList.Items) != 0 {
			for _, controllerCR := range controllerList.Items {
				controllerCR := controllerCR
				if controllerutil.ContainsFinalizer(&controllerCR, controllers.FinalizerName) {
					setupLog.Info("Operator is getting deleted. Removing finalizer from the Controller CR")
					controllerutil.RemoveFinalizer(&controllerCR, controllers.FinalizerName)
					if err = controllerReconciler.UpdateController(ctx, &controllerCR); err != nil && !errors.IsNotFound(err) {
						setupLog.Error(err, "Error while removing Finalizer from the controller")
					}
				}
			}
		}
	}
}
