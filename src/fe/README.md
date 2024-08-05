# k8s-example-fe
frontend served by nginx

## routes
| method | path      | description      |
| ------ | --------- | ---------------- |
| GET    | /static/* | static resources |
| GET    | /view/    | website homepage |

## start
```sh
# pwd: src/fe
nginx -p ./ -c ./nginx/nginx.dev.conf
# then visit http://127.0.0.1:8082
```

## build
```sh
# pwd: src/fe
VERSION=$(git describe --tags 2> /dev/null || echo unknown)

docker build . --file ./build/static.Dockerfile --tag ghcr.io/whoisnian/k8s-example-fe-static:$VERSION
docker push ghcr.io/whoisnian/k8s-example-fe-static:$VERSION
docker build . --file ./build/nginx.Dockerfile --tag ghcr.io/whoisnian/k8s-example-fe-nginx:$VERSION
docker push ghcr.io/whoisnian/k8s-example-fe-nginx:$VERSION
```
