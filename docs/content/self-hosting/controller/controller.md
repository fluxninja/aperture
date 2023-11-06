---
title: Self-Hosted Aperture Controller
description: Install Aperture Controller
keywords:
  - install
  - setup
  - controller
sidebar_position: 3
sidebar_label: Controller
---

```mdx-code-block
import CodeBlock from '@theme/CodeBlock';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import {apertureVersion, apertureVersionWithOutV} from '../../apertureVersion.js';
```

## Overview

The Aperture Controller functions as the brain of the Aperture system.
Leveraging an advanced control loop, the Controller routinely analyzes polled
metrics and indicators to determine how traffic should be shaped as defined by
set policies. Once determined, these decisions are then exported to all Aperture
Agents to effectively handle workloads.

The closed feedback loop functions primarily by monitoring the variables
reflecting stability conditions (process variables) and compares them against
setpoints. The difference in the variable values against these points is
referred to as the error signal. The feedback loop then works to minimize these
error signals by determining and distributing control actions, that adjust these
process variables and maintain their values within the optimal range.

## Installation

The Aperture Controller can be installed using the below options:

1. [**Kubernetes**](kubernetes/kubernetes.md)

   The Aperture Controller can be installed on Kubernetes using the Kubernetes
   Operator available for it, or using namespace-scoped resources.

2. [**Docker**](docker.md)

   The Aperture Controller can also be installed on Docker as containers.
