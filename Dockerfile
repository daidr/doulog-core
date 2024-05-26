FROM alpine:latest

WORKDIR /doulog_core
COPY doulog_core ./doulog_core

RUN apk update \
    && apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && chmod +x ./doulog_core

EXPOSE 5212
VOLUME ["/doulog_core/data"]

ENTRYPOINT ["./doulog_core"]