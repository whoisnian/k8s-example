# k8s-example

## todo
* [x] communication between internal services
* [x] communication with external service (MySQL)
* [x] mount persistent volume (NFS)
* [x] route requests to services (Ingress)
* [x] inject configuration data (ConfigMap)
* [x] organize yaml with kustomization
* [ ] yaml generator

## run
### with docker-compose
```sh
# pwd: run/docker-compose
docker-compose -p k8s-example up
# then visit http://127.0.0.1:8090
```

### with k8s in [minikube](https://minikube.sigs.k8s.io)
```sh
# start mysql, should be accessed through '192.168.49.1:3306'
docker volume create mysql-data
docker run -d --name mysql \
    -e MYSQL_DATABASE=k8s \
    -e MYSQL_USER=k8s \
    -e MYSQL_PASSWORD=KxY8cSAWz1WJEfs3 \
    -e MYSQL_RANDOM_ROOT_PASSWORD=yes \
    -p 3306:3306 \
    -v mysql-data:/var/lib/mysql \
    mysql:8

# start nfs, should be accessed through '192.168.49.1:2049'
mkdir /tmp/exports
docker run -d --name nfsv4-server \
  --privileged \
  -p 2049:2049 \
  -v /tmp/exports:/exports \
  ghcr.io/whoisnian/nfsv4-server-docker:0.0.2

# pwd: run/k8s
kubectl apply -f ./
# then visit http://192.168.49.2:80
```

### with kustomize in [minikube](https://minikube.sigs.k8s.io)
```sh
# start mysql, should be accessed through '192.168.49.1:3306'
docker volume create mysql-data
docker run -d --name mysql \
    -e MYSQL_DATABASE=k8s \
    -e MYSQL_USER=k8s \
    -e MYSQL_PASSWORD=KxY8cSAWz1WJEfs3 \
    -e MYSQL_RANDOM_ROOT_PASSWORD=yes \
    -p 3306:3306 \
    -v mysql-data:/var/lib/mysql \
    mysql:8

# start nfs, should be accessed through '192.168.49.1:2049'
mkdir /tmp/exports
docker run -d --name nfsv4-server \
  --privileged \
  -p 2049:2049 \
  -v /tmp/exports:/exports \
  ghcr.io/whoisnian/nfsv4-server-docker:0.0.2

# pwd: run/kustomize
kubectl apply -f ./k8s-example-namespace.yaml
kubectl apply -k ./overlays/dev
# then visit http://192.168.49.2:80
```
