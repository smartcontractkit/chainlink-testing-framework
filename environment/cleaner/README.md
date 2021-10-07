### Env cleaner
A tool to clean up stale test namespaces by applying different policies

#### Usage
Add some policy to apply for namespace, rebuild image
```
./update_image.sh
```
Run it inside the cluster
```
./start.sh
./stop.sh
```
In order for service to work you need `clusterrolebinding` set to have `remove` verb, example:
```
kubectl create clusterrolebinding cleaner-delete-ns --clusterrole=${some_role_with_remove_access} --serviceaccount=default:default```