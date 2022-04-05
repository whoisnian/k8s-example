# backend-file
persistent file storage

## routes
| method | path            | description           |
| ------ | --------------- | --------------------- |
| DELETE | /self/file/data | delete file(internal) |
| POST   | /file/data      | upload file           |
| GET    | /file/data      | download file         |

## start
```sh
cd src/backend-file
go run main.go
```
