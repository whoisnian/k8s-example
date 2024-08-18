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

https://www.elastic.co/guide/en/observability/current/apm-getting-started-apm-server.html
apm-server 依赖 fleet 的 APM integration 配置 index templates/ILM policies/ingest pipelines  
apm-server 只负责发送数据到 es

https://www.elastic.co/guide/en/fleet/8.15/fleet-server.html
Fleet Server is a subprocess that runs inside a deployed Elastic Agent. This means the deployment steps are similar to any Elastic Agent, except that you enroll the agent in a special Fleet Server policy. Typically—​especially in large-scale deployments—​this agent is dedicated to running Fleet Server as an Elastic Agent communication host and is not configured for data collection.

https://discuss.elastic.co/t/apm-server-installation-without-fleet-server-in-single-node/330600
https://github.com/elastic/apm-server/issues/10361#issue-1600406694
https://discuss.elastic.co/t/how-to-install-apm-server-in-legacy-mode-without-elastic-apm-integration/325724
https://www.elastic.co/guide/en/observability/current/apm-running-on-docker.html#_configure_apm_server_on_docker
https://www.elastic.co/guide/en/observability/current/apm-privileges-to-publish-events.html

```yaml
elastic-agent:
  image: docker.elastic.co/beats/elastic-agent:8.15.0
  restart: always
  environment:
    - KIBANA_FLEET_SETUP=1
    - KIBANA_FLEET_HOST=http://kibana:5601
    - KIBANA_FLEET_USERNAME=elastic
    - KIBANA_FLEET_PASSWORD=ClkmQKTesKG4ozWYf9G6
    - FLEET_SERVER_ENABLE=1
    - FLEET_SERVER_ELASTICSEARCH_HOST=http://elasticsearch:9200
    - FLEET_SERVER_POLICY_ID=fleet-server-policy
  ports:
    - 8220:8220
  depends_on:
    kibana:
      condition: service_healthy

# ./kibana.yml:/usr/share/kibana/config/kibana.yml
xpack.fleet.packages:
  - name: fleet_server
    version: latest
  - name: apm
    version: latest
xpack.fleet.agentPolicies:
  - name: fleet-server-policy
    id: fleet-server-policy
    namespace: default
    monitoring_enabled: []
    package_policies:
      - name: fleet_server-1
        package:
          name: fleet_server
```
