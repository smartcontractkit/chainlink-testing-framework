## ChaosMesh Playground

This deployment is used to debug various [ChaosMesh](https://chaos-mesh.org/) with [Kind](https://kind.sigs.k8s.io/)

Install [DevBox](https://www.jetify.com/devbox) and run your environment
```
devbox run up
```

Overview the services
```
devbox shell
k9s
```
Apply experiments (inside devbox shell)
```
kubectl apply -f manifests/latency.yaml
```
Debug `ChaosMesh` using k9s, check daemon logs.

Remove the environment
```
devbox run down
```