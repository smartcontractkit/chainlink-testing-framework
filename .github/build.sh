#!/bin/sh

export PATH=$PATH:$(go env GOPATH)/bin
make compile_soak_all
make compile_smoke_all