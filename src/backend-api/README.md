# backend-api
file meta data

## routes
| method | path           | description                 |
| ------ | -------------- | --------------------------- |
| POST   | /self/api/file | create file entry(internal) |
| GET    | /api/files     | list all files              |
| DELETE | /api/file      | delete file                 |

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
