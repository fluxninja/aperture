---
title: Percentage Rollouts
keywords:
  - tutorial
  - feature-flags
sidebar_position: 3
sidebar_label: Percentage Rollouts
---

## Overview

Percentage rollout is a systematic procedure of gradually introducing a new
feature to a segment of users. This practice mitigates potential risks
associated with a full-scale release, allowing the team to gather user feedback
and implement requisite modifications before launching the feature universally.

With Aperture, teams are empowered to control the launch of new features,
gradually ramp up the load, and adjust in real-time based on service health
feedback. This approach ensures that the introduction of new features does not
cause any performance regressions, thereby reducing operational risk.

<Zoom>

```mermaid
{@include: ./assets/percentage-rollouts/percentage-rollouts.mmd}
```

</Zoom>

This graph portrays the role of the Controller in actuating the rollout policy.
The rollout policy defines the 'forward' and 'reset' criteria, which dictate
when to roll forward or reset the feature, respectively. Based on these
criteria, the Controller sends the defined load ramp percentage to the Agent.
Subsequently, the Agent redirects traffic based on this percentage to the new
feature of the service.

## Example Scenario

Consider a scenario where a video streaming platform is planning to introduce a
new interface design. Instead of releasing the new design to all users at once,
the platform could use Aperture to facilitate a controlled rollout. With
Aperture's load ramping policy, the new design could be gradually introduced to
a select group of users, allowing the platform to gather feedback and make
necessary adjustments. This strategy would enable the platform to continue
releasing the new design to an increasingly larger user base, until it is
eventually launched universally. Throughout this process, Aperture would ensure
that the rollout doesn't impact the overall system performance, thereby
minimizing operational risk.

```mdx-code-block
import DocCardList from '@theme/DocCardList';
```

<DocCardList />
