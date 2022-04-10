FROM reg.whoisnian.com/proxy.docker.io/library/nginx:mainline-alpine

RUN mkdir -p /app/html/

COPY html /app/html
COPY nginx/nginx.prod.conf /app/nginx.conf

EXPOSE 8082

ENTRYPOINT [ "/usr/sbin/nginx", "-c", "/app/nginx.conf" ]