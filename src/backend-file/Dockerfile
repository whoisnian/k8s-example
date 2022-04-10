FROM reg.whoisnian.com/proxy.docker.io/library/golang:1.18-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
RUN go build -o backend-file

FROM reg.whoisnian.com/proxy.docker.io/library/alpine:3.15

WORKDIR /app

COPY --from=build /app/backend-file ./

VOLUME /app/uploads

EXPOSE 8081

ENTRYPOINT [ "/app/backend-file" ]