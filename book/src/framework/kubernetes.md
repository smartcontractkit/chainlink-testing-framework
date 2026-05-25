# Running in Kubernetes

Set up Kubernetes context for `main.stage` and add namespace variable to run `devenv`:
```bash
kubectl config use-context staging-eks-admin
export KUBERNETES_NAMESPACE="devenv-1"
```

Each Pod always has a service, you can found programmatic connection data in `.out`.

To forward all the services locally use [kubefwd]()
```bash
sudo -E kubefwd svc -n $KUBERNETES_NAMESPACE
```