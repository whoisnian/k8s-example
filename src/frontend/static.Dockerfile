FROM alpine:3.17

LABEL org.opencontainers.image.source https://github.com/whoisnian/k8s-example

COPY html /app/html

ENTRYPOINT [ "du", "-hd", "1", "/app" ]