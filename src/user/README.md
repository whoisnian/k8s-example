# k8s-example-user
user authentication

## routes
| method | path         | description                     |
| ------ | ------------ | ------------------------------- |
| POST   | /user/signup | register a new account          |
| POST   | /user/signin | login using an existing account |

## config
| env name        | default value                                                                   |
| --------------- | ------------------------------------------------------------------------------- |
| CFG_DEBUG       | false                                                                           |
| CFG_VERSION     | false                                                                           |
| CFG_LISTENADDR  | 0.0.0.0:8080                                                                    |
| CFG_MYSQLDSN    | root:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=UTC |
| CFG_AUTOMIGRATE | false                                                                           |

## start
```sh
# pwd: src/user
export CFG_MYSQLDSN="root:PD3lfKSxoXVPCdvriSHv@tcp(127.0.0.1:3306)/db_xfrcw?charset=utf8mb4&parseTime=True&loc=UTC"

CFG_AUTOMIGRATE=true go run main.go
go run main.go
```
