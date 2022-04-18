---
layout: default
title: Kubernetes
nav_order: 1
parent: Setup
---

# Kubernetes Setup

In order to use this framework, you must have a connection to an actively running [Kubernetes cluster](https://kubernetes.io/docs/tasks/access-application-cluster/access-cluster/) and an install of [kubectl](https://kubernetes.io/releases/download/). If you don't have a Kubernetes cluster handy, check out our quickstart guide on setting up a [local cluster](https://smartcontractkit.github.io/integrations-framework/quickstart/local-k8s-setup.html#local-kubernetes-cluster) which should work fine for smaller tests. Larger tests, or many tests run in parallel will likely render these local solutions inadequate. A Kubernetes cluster with 4 vCPU and 10 GB RAM is a good starting point.

**The framework will use whatever your current KUBECONFIG context is**, see how to set a context [here](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/). Learn more about setting up a Kubernetes cluster [here](https://kubernetes.io/docs/setup/). Once you have the cluster setup and are running tests, you can monitor the test architecture with the [kubectl](https://kubernetes.io/docs/reference/kubectl/kubectl/) CLI, or check out [Lens](https://k8slens.dev/) for a handy GUI.

## Why?

There's a lot of different components to bring up for each test, most of which involve:

* A simulated blockchain
* Some number of Chainlink nodes
* An equal number of postgres DBs to support the Chainlink nodes
* At least one external adapter

Following the good testing practice of having clean, non-dependent test environments means we're creating a lot of these components for each test, and tearing them down soon after. In order to organize these test environments, and after finding `docker compose` to be woefully inadequate after a certain point, Kubernetes was the obvious choice.

<div class="note note-purple">
The Kubernetes setup process, and the resources needed to run a K8s cluster for these tests is a common pain point. We're exploring ways to lessen the resources needed, and possibly expand to other systems as well as Kubernetes.
</div>
