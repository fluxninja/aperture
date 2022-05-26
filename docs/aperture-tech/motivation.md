Aperture provides Flow Control as a Service for reliably running modern
web-scale cloud applications. Aperture’s out-of-the-box adaptive API flow
control policies are orders of magnitude more effective in preventing cascading
failures and ensuring end-user experience than traditional Site Reliability
Engineering (SRE) workflows.

### Why Aperture?

For microservices interacting over the network, recognizing overload's role as a leading cause of cascading failures is a seemingly inevitable issue that works to degrade the quality of service of an application. For large-scale enterprises, it is a quick translation into downtime and revenue losses.
Considering a large, connected system, the downfall of a single service due to overload means increasing queues elsewhere in the system and increasing that service’s likelihood of failing, as well. As such, cascading failures may progressively take down multiple services in a chain, slowing the entire system and making it difficult to identify the actual root cause. Hence, modern web-scale businesses will survive based on how their teams gracefully handle overloads, not relying only on auto-scaling the infrastructure or rate-limiting the traffic at the gateway level.

Here are some of the major factors causing overloads in modern cloud applications:

#### Modern web-scale apps are vulnerable to sudden load

The graceful handling of traffic spikes (e.g., Holiday events such as Black Friday or travel booking for summer holidays, etc.) is a hard problem despite advanced capacity planning (e.g., the modeling demand and capacity is highly inaccurate). To scale microservices, there is an added level of complexity. While dealing with a single application or a load-balanced application, there might be elements of an application written in different programming languages, running on different infrastructure, running on different virtualization hypervisors, and deployed across public cloud or on-premises. When demand increases for the application, all the underlying components must be coordinated to scale, or you must be able to identify which individual elements need to scale to address the surge in demand. When we think about auto-scaling microservices for an application, we are looking at improving on two major points:

- Make sure the deployment can scale up fast in the case of a rapid load increase, so users don't see timeout errors
- Lowering the cost of the infrastructure

#### Peer services inadvertently cause DoS

For many applications, the majority of the outages are caused by changes in a live system. When you change something in your service – you deploy a new version of your feature or change some code – there is always a chance for failure or the introduction of a new bug. No software or code is perfect. Slowness and/or bugs in upstream services can cause overloads in downstream services leading to cascading failures. To prevent this domino effect, modeling and limiting concurrency on each instance of a service is key.

Most enterprise applications have a system-wide toggle switch for enabling or disabling a new feature/experience for their end customers. A more complex implementation could have a toggle switch down to individual users/end customers or even based on actions performed by the user. A dynamic feature flag controller is needed, that maintains the service level objective of an application by refreshing feature toggles when the application is under heavy load. This helps in the graceful degradation of a service without having to restart the overload application.

#### Honoring quota/tenant limits with 3rd party services

Throttling enabled by external services to manage speed or enforce varying
subscription models may lead to timed-out transactions that cause a chain
reaction. Services such as Twilio or cloud providers leverage such mechanisms at
the cost of pesky errors seen by end-users. As such, the effective solution
works backward from the 3rd party hard-limits and adjusts app functionality.

#### Failover is often messy

In another instance of the domino effect, when the primary site(s) become
unavailable, failover is pushed onto a backup, causing sudden congestion in this
new location that may mean cascading failures. During the warming period, the
backup site must see a graceful ramping-up of incoming traffic.

#### Keeping bots and user abuse in check

The final cause of sudden, difficult-to-predict overloads includes automated scrapers and bot users. Such dummy users unfairly consume precious bandwidth and
push services to the brink, heavily driving up costs. Then, rapid identification
of bot-like behavior and traffic scrubbing via distributed rate-limiting become key in minimizing such noise in distributed microservices.

#### Infrastructure/site limitation at Edge locations

For microservices applications deployed in remote edge locations, A heavy investment in reliability is needed, as the sites have the hard capacity and bandwidth limits due to space constraints. These applications cannot depend on auto-scaling and need the ability to do label-based high-performance distributed rate limiting and concurrency control.

Understanding these factors about underlying cascading failures in web-scale
cloud applications, Aperture serves as the solution

## What is Aperture?

```{toggle} Click the button to reveal!
:show:
Below Section will be completed once IP is submitted. Waiting for the claims section done by lawyer
```

Aperture leverages four pillars of service to prevent cascading failures and quota
abuse, and provide service protection via concurrency control:

#### A high-fidelity, high-frequency automated collection of only essential flow telemetry that translates into greater control.

While there may be existing and developing solutions in this space, each fails
to meet expectations or incur unnecessary costs in supporting its design. At
this time, it is common to find solutions that are built on simple conditional
reactivity, looking toward a single signal to indicate violations and inform
decisions. These solutions then fail to handle complex situations and
may lead to incorrect handling. On the other hand, some solutions rely on
collecting as much flow telemetry as possible, multiplying the amount of space
required and storage costs. Further, solutions that take on a more reactive
approach result in delays and customer dissatisfaction, counting on an SRE to
identify the blast radius of a violation after it is noticed and take action.
Noticing these weaknesses, Aperture is engineered to track only the ‘golden’
signals that are truly necessary to monitor and ensure service level
objectives (SLO).

#### An advanced control loop removes the need for user intervention in modeling activity to determine the performance of an application and make related decisions.

Using information gathered from tracking golden signals, Aperture’s unique
control loop then accomplishes automatic decision-making without the need of
site reliability engineer(SRE) oversight, allowing them to then focus time and
energy on critical tasks.

#### During overload situations, prioritized load shedding optimizes the overall user experience, allowing businesses to define the proper priority of API services over non-critical traffic.

Aperture’s third pillar is especially distinctive in its approach. While
prioritized load shedding has already been widely acknowledged as an advantage
for achieving related goals, Aperture implements this technology at the app
layer, taking a proactive approach to ensure such problems are efficiently
handled without the cost of quality or revenue.

#### Aperture effectively tracks long-term behavior in the cloud, understanding patterns to automatically detect abuse caused by bots and malicious attackers and prevent any detrimental effects.

Acknowledging that peer services that become DDoS attackers and bot traffic lend
fuel to overloads and seriously degrade customer experience, Aperture offers a
solution that can help indicate potential malicious activity by tracking cloud
history over time, identifying patterns to inform AI decision making and reporting trend insights. Together, these four pillars constitute an out-of-box active
management service that stands out from the crowd.
