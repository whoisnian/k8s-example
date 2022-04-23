# frontend
frontend served by nginx

## routes
| method | path      | description      |
| ------ | --------- | ---------------- |
| GET    | /static/* | static resources |
| GET    | /view/    | list all files   |

## start
```sh
cd src/frontend
nginx -p ./ -c ./nginx/nginx.dev.conf
# then visit http://127.0.0.1:8082
```

## build
```sh
cd src/frontend
TAG=$(cat VERSION)
docker build . -t reg.whoisnian.com/k8s-example/frontend:$TAG
docker push reg.whoisnian.com/k8s-example/frontend:$TAG
```
