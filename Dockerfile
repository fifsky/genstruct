FROM alpine:latest
RUN apk update && \
  apk add --no-cache ca-certificates \
   caddy \
   tzdata

COPY ./Caddyfile /blog/
COPY ./cert/ /blog/cert/
COPY ./blogapi /blog/
COPY ./api/config/prod/ /blog/config/prod/
COPY ./dist/ /blog/dist/
COPY ./run.sh /blog/

ENV TZ=Asia/Shanghai
EXPOSE 80 443

ENTRYPOINT ["/bin/sh", "/blog/run.sh"]