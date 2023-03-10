FROM docker.io/library/golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
RUN go build -o backend-file

FROM docker.io/library/alpine:3.17

LABEL org.opencontainers.image.source https://github.com/whoisnian/k8s-example

WORKDIR /app

COPY --from=build /app/backend-file ./

VOLUME /app/uploads

EXPOSE 8081

ENTRYPOINT [ "/app/backend-file" ]