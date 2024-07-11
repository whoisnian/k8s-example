# k8s-example-file
persistent file storage

## routes
| method | path             | description   |
| ------ | ---------------- | ------------- |
| GET    | /file/objects    | list files    |
| POST   | /file/objects    | upload files  |
| GET    | /file/object/:id | download file |
| DELETE | /file/object/:id | delete file   |

## config
| env name          | default value                                                                   |
| ----------------- | ------------------------------------------------------------------------------- |
| CFG_DEBUG         | false                                                                           |
| CFG_VERSION       | false                                                                           |
| CFG_LISTENADDR    | 0.0.0.0:8081                                                                    |
| CFG_MYSQLDSN      | root:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=UTC |
| CFG_AUTOMIGRATE   | false                                                                           |
| CFG_STORAGEDRIVER | filesystem                                                                      |
| CFG_STORAGEBUCKET | bucket01                                                                        |
| CFG_ROOTDIRECTORY | ./uploads                                                                       |
| CFG_S3ENDPOINT    | s3.amazonaws.com                                                                |
| CFG_S3ACCESSKEY   | QZH1XZPZLP8DA3GKA3J1                                                            |
| CFG_S3SECRETKEY   | VQyou21kIHVuKLkULNaETFnN7kLstyiX2KEtVbuI                                        |
| CFG_S3SECURE      | true                                                                            |

## start
```sh
# pwd: src/file
export CFG_MYSQLDSN="root:ChFHZ8Jjo9u6F3RxKbiO@tcp(127.0.0.1:3306)/demodb?charset=utf8mb4&parseTime=True&loc=UTC"
export CFG_STORAGEDRIVER=aws-s3
export CFG_S3ENDPOINT=127.0.0.1:9000
export CFG_S3ACCESSKEY=DNtNHG02un
export CFG_S3SECRETKEY=LGoucBTxlsXwhmJ9Q8aS
export CFG_S3SECURE=false
export CFG_EXTERNALSVCUSER="http://127.0.0.1:8080"

./build/build.sh . && CFG_AUTOMIGRATE=true ./output/k8s-example-file
./build/build.sh . && ./output/k8s-example-file
```
