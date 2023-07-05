---
title: Aperture Controller
description: Install Aperture Controller
keywords:
  - install
  - setup
  - controller
sidebar_position: 2
---

```mdx-code-block
import CodeBlock from '@theme/CodeBlock';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import {apertureVersion, apertureVersionWithOutV} from '../../../apertureVersion.js';
```

## Overview

The Aperture controller functions as the brain of the Aperture system.
Leveraging an advanced control loop, the controller routinely analyzes polled
metrics and indicators to determine how traffic should be shaped as defined by
set policies. Once determined, these decisions are then exported to all Aperture
agents to effectively handle workloads.

The closed feedback loop functions primarily by monitoring the variables
reflecting stability conditions (process variables) and compares them against
setpoints. The difference in the variable values against these points is
referred to as the error signal. The feedback loop then works to minimize these
error signals by determining and distributing control actions, that adjust these
process variables and maintain their values within the optimal range.

## Installation

The Aperture controller can be installed using the below options:

1. [**Kubernetes**](kubernetes/kubernetes.md)

   The Aperture controller can be installed on Kubernetes using the Kubernetes
   Operator available for it, or using namespace-scoped resources.

2. [**Docker**](docker.md)

   The Aperture controller can also be installed on Docker as containers.

## Applying Policies

Once the
[application is set up](/get-started/setting-up-application/setting-up-application.md)
and both the Aperture Controller and Agents are installed, the next crucial step
is to create and apply policies.

[Your first policy](/get-started/policies/policies.md) section provides
step-by-step instructions on customizing, creating, and applying policies within
Aperture. Additionally, the [use-cases](/use-cases/use-cases.md) section serves
as a valuable resource for tailoring policies to meet specific requirements.
