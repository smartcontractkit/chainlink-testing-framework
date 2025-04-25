## ChaosMesh Playground

This deployment is used to debug various [ChaosMesh](https://chaos-mesh.org/) with [Kind](https://kind.sigs.k8s.io/)

Install [DevBox](https://www.jetify.com/devbox) and run your environment
```
devbox run up
```

Check the services
```
devbox shell
k9s
```
Apply experiments (inside devbox shell)

If you running it from any other shell or using `Go` don't forget to apply `kubectl config set-context kind-cm-playground` before!
```
kubectl apply -f manifests/latency.yaml
```
Debug `ChaosMesh` using `k9s`, check daemon logs.

Remove the environment
```
devbox run down
```

## Using Havoc

To use our chaos testing framework and apply a single experiment either locally or remotely use
```
// main.stage
aws sso login --profile=staging-crib
kubectl config use-context main-stage-cluster-crib

// OR local chaosmesh-playground
kubectl config use-context kind-cm-playground

go test -v -run TestChaosSample
```

Open `k9s` and search your namespace for `networkchaos`