# Debugging K8s Chaos Tests

This deployment is used to debug various [ChaosMesh](https://chaos-mesh.org/) with [Kind](https://kind.sigs.k8s.io/)

Install [DevBox](https://www.jetify.com/devbox) and run your environment
```
cd infra/chaosmesh-playground
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
