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
