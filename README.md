# k8s-example

## todo
* [x] communication between internal services
* [x] communication with external service (MySQL)
* [x] mount persistent volume (NFS)
* [x] route requests to services (Ingress)
* [x] inject configuration data (ConfigMap)
* [x] organize yaml with kustomization

## test
* test/k8s-example-test.yaml:  测试 NodePort 暴露服务
* test/k8s-example-test2.yaml: 测试单 pod 多 containers 时的 pod 状态更新
* test/k8s-example-test3.yaml: 测试 juicefs-csi-driver 可用性
* test/k8s-example-test4.yaml: 测试双 namespace 使用 ingress-nginx.canary 实现灰度上线
* test/k8s-example-test5.yaml: 测试 pod 初始化 host 目录及 podAffinity 建立依赖关系
* test/k8s-example-test6.yaml: 测试 pod 利用 initContainers 强制设定启动顺序

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
