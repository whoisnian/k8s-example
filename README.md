# k8s-example

## todo
* [x] crud user snippet from redis
* [x] log with otel tracing
* [x] manually docker build
* [x] run with docker-compose
* [x] fluentbit in docker-compose
* [ ] run with raw-yaml
* [ ] fluentbit as daemonset
* [ ] run with kustomize
* [ ] compare jaegertracing with elastic-apm
* [ ] change jaeger storage backend
* [ ] prometheus metrics with grafana
* [ ] nfs shared persistent volume

## src
* fe: frontend served by nginx
* user: user info and authentication
* file: persistent file storage

## deps
`docker compose --env-file ./deps/compose.main.env --file ./deps/compose.main.yaml up`
* mysql: primary relational data storage
* redis: user sessions and temporary cache
* nfs: persistent shared file storage
* minio: self-hosted s3-compatible object storage
* jaeger: open-source distributed tracing platform (traces)
* prom: open-source systems monitoring and alerting (metrics)
* efk: open-source log aggregation, analysis, visualization (logs)

## run
* docker-compose
  ```sh
  # pwd: run/docker-compose
  docker-compose up
  # then visit http://127.0.0.1:8090
  ```
* raw-yaml (todo)
  ```sh
  # pwd: run/raw-yaml
  kubectl apply -f ./
  ```
* kustomize (todo)
  ```sh
  # pwd: run/kustomize
  kubectl create namespace k8s-example
  kubectl apply -k ./overlays/dev
  ```
