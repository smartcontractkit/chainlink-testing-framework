#!/bin/bash

# Inputs:
#   performance: 0 or 1
#     Indicates whether this is a performance test run or not.
#     - 0: No performance tests
#     - 1: Performance tests
performance=$1

if [[ $performance == 0 ]]; then
  echo "Running smoke tests"
  ginkgo -r -keepGoing --trace --randomizeAllSpecs --randomizeSuites --progress -nodes=10 -skipPackage=./suite/performance,./suite/chaos ./suite/...
else
  echo "Running performance and chaos tests"
  ginkgo -r -keepGoing --trace --randomizeAllSpecs --randomizeSuites --progress ./suite/performance ./suite/chaos
fi