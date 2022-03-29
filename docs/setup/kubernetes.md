---
layout: default
title: Kubernetes
nav_order: 1
parent: Setup
---

# Kubernetes Setup

In order to use this framework, you must have a connection to an actively running [Kubernetes cluster](https://kubernetes.io/docs/tasks/access-application-cluster/access-cluster/). If you don't have one handy, check out [k3d](https://k3d.io/) or [minikube](https://minikube.sigs.k8s.io/docs/start/) which should work fine for smaller tests, but if you write tests that make use of multiple Chainlink nodes, or try to run many tests in parallel, you'll likely find these local solutions inadequate. A Kubernetes cluster with 4 vCPU and 10 GB RAM is a good starting point for when you start to notice issues.

**The framework will use whatever your current KUBECONFIG context is**, see how to set a context [here](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/). Learn more about setting up a Kubernetes cluster [here](https://kubernetes.io/docs/setup/). Once you have the cluster setup and are running tests, you can monitor the test architecture with the [kubectl](https://kubernetes.io/docs/reference/kubectl/kubectl/) CLI, or check out [Lens](https://k8slens.dev/) for a handy GUI.

## Why?

There's a lot of different components to bring up for each test, most of which involve:

* A simulated blockchain
* Some number of Chainlink nodes
* An equal number of postgres DBs to support the Chainlink nodes
* At least one external adapter

Following the good testing practice of having clean, non-dependent test environments means we're creating a lot of these components for each test, and tearing them down soon after. In order to organize these test environments, and after finding `docker compose` to be woefully inadequate after a certain point, Kubernetes was the obvious choice.
