FROM nginx:mainline-alpine

LABEL org.opencontainers.image.source https://github.com/whoisnian/k8s-example

COPY html /app/html
COPY nginx/nginx.prod.conf /app/nginx.conf

EXPOSE 8082

ENTRYPOINT [ "/usr/sbin/nginx", "-p", "/app/", "-c", "/app/nginx.conf" ]