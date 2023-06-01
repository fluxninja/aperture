---
title: Feature Rollout
keywords:
  - tutorial
  - feature-flags
sidebar_position: 1
sidebar_label: Feature Rollout
---

## Overview

Feature Rollout is a systematic procedure of gradually introducing a new feature
to a segment of users. This practice mitigates potential risks associated with a
full-scale release, allowing the team to gather user feedback and implement
requisite modifications before launching the feature universally.

<Zoom>

```mermaid
{@include: ../assets/feature-rollout.mmd}
```

</Zoom>

This graph portrays the role of the Controller in actuating the rollout policy.
Leveraging forward and reset parameters, the Controller sends the defined load
ramp percentage to the Agent. Subsequently, the Agent redirects traffic based on
this percentage to the new feature of the service.

## Real World Example

In terms of Feature Rollout, imagine a social media platform planning to
introduce a new interface design. Instead of releasing the new design to all
users at once, they gradually roll it out to select users, gather feedback, make
necessary adjustments if needed, or keep releasing it to an increasing number of
users until it is released to all users.

```mdx-code-block
import DocCardList from '@theme/DocCardList';
```

<DocCardList />
