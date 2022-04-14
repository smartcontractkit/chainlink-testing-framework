---
layout: default
title: Writing a Test
nav_order: 7
has_children: true
---

# Writing a Test

Here we walk through writing a couple different types of tests. More examples can be found in our [test suite directory](https://github.com/smartcontractkit/integrations-framework/tree/main/suite).

We’ve been using [Ginkgo](https://github.com/onsi/ginkgo), a BDD testing framework for Go that we’ve found handy for organizing and running tests. You should be able to use any other testing framework you like, including Go’s built-in `testing` package, but the examples you find here will be in Ginkgo and its accompanying assertions library, [Gomega](https://github.com/onsi/gomega). Both of which are easily human readable however (part of why we like it) and shouldn’t hurt your understanding of the framework at all.
