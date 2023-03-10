FROM golang:1.20-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
RUN go build -o feature-test

FROM alpine:3.17

LABEL org.opencontainers.image.source https://github.com/whoisnian/k8s-example

WORKDIR /app

COPY --from=build /app/feature-test ./

EXPOSE 8080

ENTRYPOINT [ "/app/feature-test" ]