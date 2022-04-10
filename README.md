# k8s-example

## test
[x] connect to internal service  
[ ] connect to external service(database)  
[ ] persistent volume  
[x] route requests to services  
[ ] configMap  

## run
### docker-compose
```sh
cd run/docker-compose
docker-compose up
nginx -p ./ -c ./nginx.test.conf
# then visit http://127.0.0.1:8090
```

### k8s
```sh
cd run/k8s
kubectl apply -f ./
# kubectl delete -f ./
```
