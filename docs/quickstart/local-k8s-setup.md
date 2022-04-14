---
layout: default
title: Local Kubernetes Setup
nav_order: 1
parent: Quick Start
---

# Local Kubernetes Setup

In order to run tests, we need a [Kubernetes](https://kubernetes.io/) (often abbreviated as *K8s*) cluster, and a way to connect to it (see [why Kubernetes](https://smartcontractkit.github.io/integrations-framework/setup/kubernetes.html#why)).

In order to connect to a K8s cluster, we need to [download kubectl](https://kubernetes.io/releases/download/).

Next you'll need to get a Kubernetes cluster running. If you're lucky, you'll have one running in the cloud somewhere that you can utilize. But if you don't, or just want to get something running locally, you can use [K3D](https://k3d.io/v5.4.1/#installation) to launch a local Kubernetes cluster.