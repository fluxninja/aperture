## Overview

The Aperture architecture consists of two major components – the controller and
agent. Together, these entities collect and analyze data as well as execute
decisions to efficiently shape traffic and manage overload conditions.

## Requirements

Review the following requirements before you proceed with the installation: [Requirements](./requirements.md)

## Aperture Controller

### Overview

The Aperture Controller functions as the primary decision maker of the system.
Leveraging our advanced control loop, the controller routinely analyzes polled
metrics and indicators to determine how traffic should be shaped as defined by
set policies. Once determined, these decisions are then exported to all agents
in order to effectively handle workloads. Only one controller is needed to
effectively manage each cluster.

The closed feedback loop functions primarily by monitoring the variables
reflecting stability conditions (i.e. process variables) and compares them
against set points. The difference in the variable values against these points
is referred to as the error signal. The feedback loop then works to minimize
these error signals by determining and distributing specific decisions, or
control actions, that adjust these process variables and maintain their values
within the optimal range.

### Configuration

The Aperture Controller related configurations are stored in a configmap which
is created during the installation using Helm. All the configuration parameters
are listed on the [README](link_to_chart_readme_file) file of the Helm chart.

### Installation

The Aperture Controller can be run in Kubernetes environment and only a single
instance of the Controller is required per cluster. Please follow the
installation instructions of [The Aperture Agent](#agent-installation) as that
will install the Aperture Controller as well.

## Aperture Agent

### Overview

The Aperture Agent is the decision executor. In addition to gathering data, the
Aperture Agent functions as a gatekeeper, acting on traffic based on decisions
made by the controller. Specifically, depending on feedback from the controller,
the agent will effectively allow or drop incoming requests. Further supporting
the controller, the agent works to inject information into traffic, including
the specific traffic-shaping decisions made and classification labels, which can
later be used in policing. One agent is deployed per node.

### Configuration

The Aperture Agent related configurations are stored in a configmap which is
created during the installation using Helm. All the configuration parameters are
listed on the [README](link_to_chart_readme_file) file of the Helm chart.

### Installation {#agent-installation}

(Consult [Supported Platforms](./supported-platforms) before installing.)

Below are the steps to install or upgrade the Aperture Agent and Controller into
your setup using the [Aperture Agent Helm chart](link_to_helm_chart).

By following these instructions, you will have deployed the Aperture Agent and
Controller into your cluster.

1. Add the Helm chart repo in your environment:

   ```bash
   helm repo add aperture LINK_TO_THE_HELM_CHART
   helm repo update
   ```

   You may need to install the missing dependencies

   ```bash
   helm dependency build aperture/agent
   ```

2. Install or upgrade the chart:

   ```bash
   helm upgrade --install agent aperture/agent
   ```

3. [**OPTIONAL**] If you want to connect the Aperture Agent and Controller with
   the FluxNinja cloud, create a `values.yaml` file with below parameters:

   ```yaml
   cloudIntegration: true
   agent:
     config:
       ingestion:
         address: "INGESTION_SERVICE_ADDRESS"
         port: 443
         insecure: false
     apiKeySecret:
       value: "AGENT_API_KEY"
   agentController:
     apiKeySecret:
       value: "CONTROLLER_API_KEY"
   heartbeats:
     serverAddress: "AGENT_SERVICE_ADDRESS"
     port: 443
     tls:
       insecureSkipVerify: true
   ```

   To generate the `AGENT_API_KEY` and `CONTROLLER_API_KEY`, please follow
   instructions on
   [FluxNinja Sign-Up](https://docs.dev.fluxninja.com/docs/Getting%20started/sign_up)
   to get an account and
   [Generate API Keys](https://docs.dev.fluxninja.com/docs/Agent/Agent%20Management).

4. The chart installs Istio, Prometheus and Etcd instances by default. If you
   don't want to install and use your existing instances of Istio, Prometheus or
   Etcd, configure below values in the `values.yaml` file and pass it with
   `helm upgrade`:

   ```yaml
   etcd:
     enabled: false

   prometheus:
     enabled: false

   istio:
     enabled: false

   agent:
     config:
       etcd:
         endpoints: ["ETCD_INSTANCE_ENDPOINT"]
       prometheus:
         address: "PROMETHEUS_INSTANCE_ADDRESS"

   agentController:
     config:
       etcd:
         endpoints: ["ETCD_INSTANCE_ENDPOINT"]
       prometheus:
         address: "PROMETHEUS_INSTANCE_ADDRESS"
   ```

   ```bash
   helm upgrade --install agent aperture/agent -f values.yaml
   ```

   A list of other configurable parameters for Istio, Etcd and Prometheus can be
   found in the [README](link_to_chart_readme_file).

   **Note**: Please make sure that the flag `web.enable-remote-write-receiver`
   is enabled on your existing Prometheus instance as it is required by the
   Agent.

5. The chart also installs a `EnvoyFilter` resource for collecting data from the
   running applications. The details about the configuration and what details
   are being collected are available at [Envoy Filter](./istio#envoy-filter). If
   you do not want to install the embedded Envoy Filter and want to install by
   yourself, configure below value in the `values.yaml` file and pass it with
   `helm upgrade`:

   ```yaml
   istio:
     envoyFilter:
       install: false
   ```

   ```bash
   helm upgrade --install agent aperture/agent -f values.yaml
   ```

6. If you want to modify the default parameters, you can create or update the
   `values.yaml` file and pass it with `helm upgrade`:

   ```bash
   helm upgrade --install agent aperture/agent -f values.yaml
   ```

   A list of configurable parameters can be found in the
   [README](link_to_chart_readme_file).

7. If you want to deploy the Aperture Agent and Controller into a namespace
   other than `default`, use the `-n` flag:

   ```bash
   NAMESPACE="aperture-system"; helm upgrade --install agent aperture/agent -f values.yaml --set global.istioNamespace=$NAMESPACE -n $NAMESPACE --create-namespace
   ```

8. Once you have successfully deployed the Helm release, confirm that the
   Aperture Agent and Controller are up and running:

   ```bash
   kubectl get pod -A
   ```

   You should see pods for Prometheus, Etcd, Aperture Agent (per node), and
   Aperture Controller in `RUNNING` state.

## Aperture API Gateway

The API endpoints available from the Aperture Agent and Controller are listed on
the [Aperture APIs](./aperture-api).
