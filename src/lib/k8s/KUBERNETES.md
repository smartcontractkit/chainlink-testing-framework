# Kubernetes


<div class="warning">

Managing k8s is challenging, so we've decided to separate `k8s` deployments here - [CRIB](https://github.com/smartcontractkit/crib)

This documentation is outdated, and we are using it only internally to run our soak tests. For `v2` tests please check [this example](../crib.md) and read [CRIB docs](https://github.com/smartcontractkit/crib)
</div>

We run our software in Kubernetes.

### Local k3d setup

1. `make install`
2. (Optional) Install `Lens` from [here](https://k8slens.dev/) or use `k9s` as a low resource consumption alternative from [here](https://k9scli.io/topics/install/)
   or from source [here](https://github.com/smartcontractkit/helmenv)
3. Setup your docker resources, 6vCPU/10Gb RAM are enough for most CL related tasks
4. `make create_cluster`
5. `make install_monitoring` Note: this will be actively connected to the server, the final log when it is ready is`Forwarding from [::1]:3000 -> 3000` and you can continue with the steps below in another terminal.
6. Check your contexts with `kubectl config get-contexts`
7. Switch context `kubectl config use-context k3d-local`
8. Read [here](README.md) and do some deployments
9. Open Grafana on `localhost:3000` with `admin/sdkfh26!@bHasdZ2` login/password and check the default dashboard
10. `make stop_cluster`
11. `make delete_cluster`

### Typical problems

1. Not enough memory/CPU or cluster is slow
   Recommended settings for Docker are (Docker -> Preferences -> Resources):
   - 6 CPU
   - 10Gb MEM
   - 50-150Gb Disk
2. `NodeHasDiskPressure` errors, pods get evicted
   Use `make docker_prune` to clean up all pods and volumes
