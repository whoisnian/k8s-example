# k8s-example-user
user authentication

## routes
| method | path         | description                     |
| ------ | ------------ | ------------------------------- |
| POST   | /user/signup | register a new account          |
| POST   | /user/signin | login using an existing account |

## config
| env name                | default value                                                                   |
| ----------------------- | ------------------------------------------------------------------------------- |
| CFG_DEBUG               | false                                                                           |
| CFG_VERSION             | false                                                                           |
| CFG_LISTENADDR          | 0.0.0.0:8080                                                                    |
| CFG_MYSQLDSN            | root:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=UTC |
| CFG_AUTOMIGRATE         | false                                                                           |
| CFG_REDISURI            | redis://default:password@127.0.0.1:6379/0                                       |
| CFG_DISABLEREGISTRATION | false                                                                           |
| CFG_TRACEENDPOINTURL    | http://127.0.0.1:4318                                                           |

## start
```sh
# pwd: src/user
export CFG_MYSQLDSN="root:ChFHZ8Jjo9u6F3RxKbiO@tcp(127.0.0.1:3306)/demodb?charset=utf8mb4&parseTime=True&loc=UTC"
export CFG_REDISURI="redis://default:R5NjwH9uKh8vuZY1R2fd@127.0.0.1:6379/0"

./build/build.sh . && CFG_AUTOMIGRATE=true ./output/k8s-example-user
./build/build.sh . && ./output/k8s-example-user
```

## build
```sh
# pwd: src/user
MODULE_NAME=$(go mod edit -fmt -print | grep -Po '(?<=^module ).*$')
APP_NAME="k8s-example-user"
BUILDTIME=$(date --iso-8601=seconds)
VERSION=$(git describe --tags 2> /dev/null || echo unknown)

DOCKER_BUILDKIT=1 docker build \
  --file ./build/Dockerfile \
  --progress=plain \
  --platform=linux/amd64 \
  --build-arg MODULE_NAME="$MODULE_NAME" \
  --build-arg APP_NAME="$APP_NAME" \
  --build-arg VERSION="$VERSION" \
  --build-arg BUILDTIME="$BUILDTIME" \
  --tag ghcr.io/whoisnian/k8s-example-file:$VERSION \
  .
docker push ghcr.io/whoisnian/k8s-example-file:$VERSION
```
