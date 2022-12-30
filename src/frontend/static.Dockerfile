FROM reg.whoisnian.com/proxy.docker.io/library/alpine:3.17

COPY html /app/html

ENTRYPOINT [ "du", "-hd", "1", "/app" ]