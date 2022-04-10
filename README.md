# k8s-example

## test
[ ] connect to internal service  
[ ] connect to external service(database)  
[ ] persistent volume  
[ ] route requests to services  
[ ] configMap  

## run
### docker-compose
```sh
cd run/docker-compose
docker-compose up
nginx -p ./ -c ./nginx.test.conf
# then visit http://127.0.0.1:8090
```
