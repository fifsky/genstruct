FROM alpine:latest
RUN apk update && \
  apk add --no-cache ca-certificates \
   tzdata

WORKDIR /app

COPY ./app ./

ENV TZ=Asia/Shanghai
EXPOSE 8989

ENTRYPOINT ["./app"]