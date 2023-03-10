# feature-test
* yaml/feature-test1.yaml: 测试 NodePort 暴露服务
* yaml/feature-test2.yaml: 测试单 pod 多 containers 时的 pod 状态更新
* yaml/feature-test3.yaml: 测试双 namespace 使用 ingress-nginx.canary 实现灰度上线
* yaml/feature-test4.yaml: 测试 pod 初始化 host 目录及 podAffinity 建立依赖关系
* yaml/feature-test5.yaml: 测试 pod 利用 initContainers 强制设定启动顺序

## routes
| method | path      | description                         | usage                                      |
| ------ | --------- | ----------------------------------- | ------------------------------------------ |
| GET    | /ping     | test pod termination                | watch -en 1 curl -sf "127.0.0.1:8080/ping" |
| GET    | /mem      | test node OOM Killer                | curl -f "127.0.0.1:8080/mem?cnt=1024"      |
| GET    | /upstream | test multiple containers in one pod | curl -f "127.0.0.1:8080/upstream"          |
| GET    | /podname  | test pod env.valueFrom              | curl -f "127.0.0.1:8080/podname"           |

## start
```sh
# pwd: feature-test
go run main.go
```

## build
```sh
# pwd: feature-test
cp ../go.mod ../go.sum ./
TAG=$(cat VERSION)
docker build . -t ghcr.io/whoisnian/feature-test:$TAG
docker push ghcr.io/whoisnian/feature-test:$TAG
rm ./go.mod ./go.sum
```
