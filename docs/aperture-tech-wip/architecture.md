---
title: Aperture Architecture
slug: /architecture
description: Aperture Architecture description page
keywords:
  - architecture
  - design
  - charts
---

# Architecture Diagrams

![Aperture Architecture](../assets/img/architecture.png "Aperture Architecture")

![Aperture Flowchart](../assets/img/flowchart.png "Aperture Flowchart")

Here we describe key modules and flowchart that are shown above in the Aperture
Architecture.

## [Aperture Controller](../Introduction/introduction.md#aperture-controller)

### Policy Control Loop

The objective of control theory is to develop a model or algorithm governing the
application of system inputs to drive the system to the desired state. To do
this, a controller with the requisite corrective behavior is required. This
controller monitors the controlled process variable (PV) and compares it with
the reference or set point (SP). The difference between the actual and desired
value of the process variable called the error signal or SP-PV error is applied
as feedback to generate a control action to bring the controlled process
variable to the same value as the set point.

For every API request that occurs in the system, the Aperture plugin or library
layer communicates with the
[Aperture Agent](#aperture-agent)
twice. The first call’s main purpose is to check if the API request should be
let through or not via schedulers (explained in more detail below).
Additionally, in this step two auxiliary features are performed – execution of
flow classification rules (described in more detail below) and matching
[FluxMeters](#fluxmeter) that apply to this request.

The second point of communication between plugins/libraries is at the end of the
request. At this point we know the full duration of a request and all its
metadata, thus we can send a trace event for metrics processing. The important
part of event metadata is the actuators and [FluxMeters](#fluxmeter) that were matched for
this request. They will be used to create signals in the short-timescale-metrics
database for the consumption of the control loop circuit. Note that trace
reporting is not required to be real-time and can run in the background/be
batched within a reasonable margin.

#### Flow Classification Rule

All user workflows in the FluxNinja solution are service centric. Automatic
metadata collection to help discover the service labels are the key building
blocks. labels allows us to show how nodes are related to one
another(hierarchy), what traffic is flowing between any given nodes and helps
easily answer the question of what all nodes+traffic flows make up an
application. Through the use of Group by and filters, you can offer even more
details about those relationships.

To perform abuse prevention of a service, we need to actuate a dynamic rate
limit controller applied to flows labels. Aperture provides the ability to
implement flow classification rules.One of the foundational pieces that lead to
flexible observability and control is the ability to dynamically identify flows
at the level of users, locations, JWT tokens, URLs and so on by extracting
fields of interest from HTTP requests. This capability also provides the ability
to control the traffic (e.g. rate limit, concurrency control) based on
classification key(s) in a prioritized manner. This capability is to be built
using open policy agent/Rego.

As part of the classification step, FluxNinja library also inserts tracing
baggage headers so as to watermark the original flow and child flows – E.g. For
each classification rule that is a match we insert <flow_class = flow_key>
header in the API. These classification headers are used for fine-grained
observability into flows and flow control.

#### Flow Control Policies

Overload management techniques are based on an adaptive concurrency control
mechanism that attempts to apply WFQ algorithm on traffic flowing through a
service. This is accomplished by internally monitoring the performance of the
service, which is decomposed into a set of event-driven components connected
with signals. By calculating the rate(desired concurrency) at which requests
should be admitted, the service can perform focused overload management, for
example, by blocking or dropping only those requests that lead to resource
bottlenecks.

The flow control policy is defined as a circuit.

- The circuit defines a signal processing network as a list of components.
- Components are computational blocks that form the circuit.
- Signals flow into the components via input ports and results are emitted on
  output ports. Signals can be external (e.g., coming from FluxMeters via
  Prometheus) or internal (coming from other components).
- Components are wired to each other based on signal names forming an execution
  graph of the circuit.
- Loops are broken by the runtime at the earliest component index that is part
  of the loop. The looped signals are saved in the tick they are generated and
  served in the subsequent tick.

Execution of flow control policy can be separated into three areas: data
collection, control loop itself, and actuation. The control loop itself runs on
Aperture Controller, but the data collection and actuation need help from the
components that the actual traffic is flowing through – namely proxies and
libraries. These need to communicate through the Aperture Agent.

circuit provides a way to use any controller (gradient controller (AIMD, MIMD),
PID controller, etc.) as part of the flow control policy circuit. There are
actuators as part of our data-plane components in the flow control policy
circuit that get installed using Selector that is based on infra or flow labels.

![Example Circuit Policy](../assets/img/circuit.png "Example Circuit Policy")

Referring to the above E.g., in a particular embodiment, using such a circuit
model, we showcase an example policy that protects the service trying to
maintain latency objectives. There’s a signal, latency, that pulls the measured
request latency via the FluxMeter. Note that the FluxMeter creates a histogram,
but the signal is a scalar value. For the conversion, a query language (E.g.,
Not limiting to PromQL - Prometheus query language) could extract e.g., 95th
percentile as a signal. This latency signal is then used to power the Gradient
Controller. Note that the same signal is used as a signal and setpoint inputs to
the Gradient Controller. The setpoint is routed via the Exponential Moving
Average filter to achieve a result of “learning” the idle latency. One of the
other notable inputs to the Controller is the Control Variable – accepted
concurrency in this case. Based on all the inputs, for every circuit evaluation
tick, the gradient controller computes what’s the desired concurrency that aims
to bring the input signal to the setpoint. This is then subtracted from and
divided by incoming concurrency (using “arithmetic” components), resulting in a
Load Shed Factor signal, that’s the input of the Concurrency Actuator. This
circuit forms a simple policy that enables protection from traffic spikes and
helps achieve service level objectives of maintaining latency.

The same circuit building blocks that constitute this policy could also be used
to create more complex policies – they could use multiple input signals (e.g.,
combining multiple latency measurements or incorporating CPU and memory usage)
and multiple controllers.

## [Aperture Agent](../Setup/aperture-setup.md#aperture-agent)

### Flow Control Service

It provides a new innovative scheduler that operates on labels gathered from
flows (traffic or code). This is a very different usage and implementation of
the Weighted Fair Queuing (WFQ) scheduler than what is seen in the network layer
adopted by networking companies like Cisco, Arista, etc.

The scheduler’s job is to shape the traffic (by blocking chosen requests) to
fulfill the goal of the decision message (e.g., to achieve a certain load shed
factor). The scheduler classifies a request as belonging to a workload and
returns a decision on whether a request was accepted. Scheduler considers
workload priorities, fairness, estimated weight of a request, and current load
shed factor. Schedulers can accept or reject a request immediately or after some
short period (on the order of milliseconds). If all matching schedulers accept a
request, the request is considered accepted, and such information is reported
back to the Aperture plugin/library, which then can either pass the request
forward or block it.

#### Concurrency Control

Little’s law states that to compute the concurrency, we need to know the average
number of requests per second, but also, the average latency. We assign each
request several tokens, where the number of tokens is proportional to the
processing time of a request. To calculate the concurrency precisely, we assign
a request to a different workload, and for each workload, we estimate the
expected number of tokens used by a request (auto tokens feature). The scheduler
calculates the available number of tokens per second, based on current
concurrency and load shed factor, and then for each request, it decides if the
request may go through, based on the number of estimated tokens, number of
already used tokens, and priority of a request. Each request has its timeout
(configured per workload) in which it should be scheduled. If it’s not scheduled
within that timeout (because of a lack of available tokens in the scheduler),
it’s considered dropped.

##### Actuator

The process of actuation is split between two components. The actuator is
firstly a component in the control loop circuit, which knows what’s the desired
Load Shed Factor (as calculated by the circuit). To perform the load shed, the
actuator needs to drop some traffic. Because of this, the actuator has also a
second part running in Aperture Agent that handles the actual execution of flow
control policy. The Controller’s part of the actuator communicates with the
Agent’s part via periodically updating a “decision message”. The decision
message could contain information such as the current load shed factor for an
agent, but it could also be more complex. Note that this channel of
communication is directional – from Controller to Agents. The opposite direction
of communication is realized via signals. This “decision message” is then used
by the Scheduler. Every Agent has a copy of a Scheduler corresponding to an
Actuator from a policy (multiple policies or multiple actuators in a policy
could result in multiple schedulers). For every request, Agent executes every
Scheduler that matches the request. Matching is based on selectors – each
actuator is defined in the policy to act on a specific part of a traffic –
selector is defined as a service, control point and optionally set of flow
labels.

##### FluxMeter

Its a metrics meter for used-defined slice of traffic, defined by selector on
flow labels. In other words, matching metrics from database for specific labels
created by classification rule

##### Selector

Its a definition to where a fluxmeter/actuator or a flow classification rule may
apply to. Selector points into a particular control point of a particular
service. Optionally a selector may apply only to a subset of traffic, defined by
flow label selector

##### API Quotas

It provides a way to do high performance distributed rate limiting based on flow
labels.

For protection from malicious users, bots, or limits imposed by external
services, static quotas can be imposed on a service. Static quotas require a
specified rate limit key, for which the quota will be counted. The extraction of
the rate limit key can be powered by classification rules and can leverage
baggage propagation to allow using labels created earlier in the request chain.
Quota counters are implemented via a distributed counter, synced between agents.
Note that the static quotas system can be thought as different kind of actuator
and plugs into the policy-as-circuit model, which powers the quotas system with
flexibility, e.g., precise limit could be a predefined constant, but could also
be derived from some signals instead. Thanks to flow label matching, different
quotas can also be applied to different kinds of traffic.

The rate limit value can be dynamic coming from the circuit to prevent abuse
prevention. limit_key identifies the flows which are to be limited. Example,
this key could be used. The counters are distributed and updated in real time.
So, even small violations in limits are enforced via this mechanism. Customers
can implement Rate Limits per user via this. An example, AWS limits API requests
per API key per minute. This helps enforce limits in the application. Another
example is rate limit escalation policy for bots and changing the rate limit
dynamically: if service is overloaded for 2min (Concurrency Limit Load
Shed >50%), stop BoT traffic.

Distributed rate limiter is designed mostly for speed rather than accuracy. In
the CAP-theorem spectrum, we sacrifice consistency in favor of availability and
partition tolerance. Rate limit counters are stored locally (in memory) on
agents and are always available, counter synchronization may be running in “lazy
mode.” This compromise allows for the rate limit feature not to introduce
additional latency, as compared to the central rate limit system. The cost of
compromise here is the possibility that in some rare occurrences a few
additional requests may get through, but that’s a right compromise in order to
prevent the bigger issues of abuse prevention of assigned quotas.
