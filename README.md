# k8s-example

## todo
* [x] communication between internal services
* [x] communication with external service (MySQL)
* [x] mount persistent volume (NFS)
* [x] route requests to services (Ingress)
* [x] inject configuration data (ConfigMap)
* [x] organize yaml with kustomization

## run
### with docker-compose
```sh
cd run/docker-compose
docker-compose up
# then visit http://127.0.0.1:8090
```

### with k8s
```sh
cd run/k8s
kubectl apply -f ./
# then visit http://192.168.122.221:30080
```

### with kustomize
```sh
cd run/kustomize
kubectl apply -f ./k8s-example-namespace.yaml
kubectl apply -k ./overlays/dev
# then visit http://192.168.122.221:30080
```
