FROM registry.pinsvc.net/mirror/alpine

RUN apk add --no-cache libc6-compat && apk add tzdata && apk add curl
ENV TZ Asia/Tehran
COPY app /opt
EXPOSE 9000
WORKDIR /opt
CMD ["./app"]
