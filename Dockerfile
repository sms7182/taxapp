FROM registry.pinsvc.net/mirror/alpine

RUN apk add --no-cache libc6-compat && apk add tzdata && apk add curl
ENV TZ Asia/Tehran
COPY app /opt
COPY db/migrations/*.sql /opt/db/migrations/
WORKDIR /opt
CMD ["./app"]
