# k8s-example-frontend
frontend served by nginx

## routes
| method | path      | description      |
| ------ | --------- | ---------------- |
| GET    | /static/* | static resources |
| GET    | /view/    | list all files   |

## start
```sh
# pwd: src/frontend
nginx -p ./ -c ./nginx/nginx.dev.conf
# then visit http://127.0.0.1:8082
```

## build
```sh
# pwd: src/frontend
TAG=$(cat VERSION)
docker build . -f static.Dockerfile -t ghcr.io/whoisnian/k8s-example-frontend-static:$TAG
docker push ghcr.io/whoisnian/k8s-example-frontend-static:$TAG
docker build . -f nginx.Dockerfile -t ghcr.io/whoisnian/k8s-example-frontend-nginx:$TAG
docker push ghcr.io/whoisnian/k8s-example-frontend-nginx:$TAG
```
