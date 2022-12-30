# backend-file
persistent file storage

## routes
| method | path            | description            |
| ------ | --------------- | ---------------------- |
| DELETE | /self/file/data | delete file (internal) |
| POST   | /file/data      | upload file            |
| GET    | /file/data      | download file          |

## config
| env name       | default value         |
| -------------- | --------------------- |
| CFG_LISTENADDR | 0.0.0.0:8081          |
| CFG_APIPREFIX  | http://127.0.0.1:8080 |
| CFG_ROOTPATH   | ./uploads             |

## start
```sh
cd src/backend-file
go run main.go
```

## build
```sh
cd src/backend-file
cp ../../go.mod ../../go.sum ./
TAG=$(cat VERSION)
docker build . -t reg.whoisnian.com/k8s-example/backend-file:$TAG
docker push reg.whoisnian.com/k8s-example/backend-file:$TAG
rm ./go.mod ./go.sum
```
