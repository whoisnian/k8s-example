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
