# k8s-example-test
test k8s features

## routes
| method | path  | description          | usage                                      |
| ------ | ----- | -------------------- | ------------------------------------------ |
| GET    | /ping | test pod termination | watch -en 1 curl -sf "127.0.0.1:8080/ping" |
| GET    | /mem  | test node OOM Killer | curl -f "127.0.0.1:8080/mem?cnt=1024"      |

## start
```sh
cd ./test
go run main.go
```

## build
```sh
cd ./test
cp ../go.mod ../go.sum ./
TAG=$(cat VERSION)
docker build . -t reg.whoisnian.com/k8s-example/k8s-example-test:$TAG
docker push reg.whoisnian.com/k8s-example/k8s-example-test:$TAG
rm ./go.mod ./go.sum
```
