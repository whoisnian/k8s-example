FROM nginx:1.27-alpine

COPY html /app/html
COPY nginx/nginx.prod.conf /app/nginx.conf

EXPOSE 8082

ENTRYPOINT [ "/usr/sbin/nginx", "-p", "/app/", "-c", "/app/nginx.conf" ]