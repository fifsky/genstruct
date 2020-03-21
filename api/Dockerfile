FROM alpine:latest
RUN apk update && \
  apk add --no-cache ca-certificates \
   tzdata

WORKDIR /app

COPY ./app ./
COPY .env ./

ENV TZ=Asia/Shanghai
EXPOSE 80

ENTRYPOINT ["./app",":80"]