# Running in Kubernetes

Set up Kubernetes context for `main.stage` and add namespace variable to run `devenv`:
```bash
kubectl config use-context staging-eks-admin
export KUBERNETES_NAMESPACE="devenv-1"
```