FROM reg.whoisnian.com/proxy.docker.io/library/golang:1.18-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
RUN go build -o k8s-example-test

FROM reg.whoisnian.com/proxy.docker.io/library/alpine:3.15

WORKDIR /app

COPY --from=build /app/k8s-example-test ./

EXPOSE 8080

ENTRYPOINT [ "/app/k8s-example-test" ]