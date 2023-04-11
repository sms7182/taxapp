FROM registry.pinsvc.net/mirror/alpine

RUN apk add --no-cache libc6-compat && apk add tzdata && apk add curl && apk add gcompat
ENV TZ Asia/Tehran
COPY app /opt
COPY db/migrations/*.sql /opt/db/migrations/
COPY sign_ara.key /opt
COPY sign_delijan.key /opt
COPY templates /opt/templates
WORKDIR /opt
CMD ["./app"]
