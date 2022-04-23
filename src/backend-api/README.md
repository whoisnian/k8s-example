# backend-api
file meta data

## routes
| method | path           | description                  |
| ------ | -------------- | ---------------------------- |
| POST   | /self/api/file | create file entry (internal) |
| GET    | /api/files     | list all files               |
| DELETE | /api/file      | delete file                  |

## config
| env name    | default value                                |
| ----------- | -------------------------------------------- |
| LISTEN_ADDR | 0.0.0.0:8080                                 |
| FILE_PREFIX | http://127.0.0.1:8081                        |
| MYSQL_DSN   | k8s:KxY8cSAWz1WJEfs3@tcp(127.0.0.1:3306)/k8s |

## start
```sh
cd src/backend-api
go run main.go
```

## devDB
```sh
docker volume create mysql-data
docker run --net host -d --name mysql-dev \
    -e MYSQL_DATABASE=k8s \
    -e MYSQL_USER=k8s \
    -e MYSQL_PASSWORD=KxY8cSAWz1WJEfs3 \
    -e MYSQL_RANDOM_ROOT_PASSWORD=yes \
    -v mysql-data:/var/lib/mysql \
    mysql:8
```

## build
```sh
cd src/backend-api
cp ../../go.mod ../../go.sum ./
TAG=$(cat VERSION)
docker build . -t reg.whoisnian.com/k8s-example/backend-api:$TAG
docker push reg.whoisnian.com/k8s-example/backend-api:$TAG
rm ./go.mod ./go.sum
```
