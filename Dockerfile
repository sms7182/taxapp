FROM golang:alpine
COPY ./ /usr/src/app/taxapp

WORKDIR /usr/src/app/taxapp/cmd

COPY ./conf/dev-conf.yaml ./conf
COPY ./taxDep  ./taxDep
COPY ./pkg ./pkg
COPY ./db/migrations/*.sql ./db/migrations/
COPY ./go.mod ./
COPY ./go.sum ./
COPY ./notify ./notify
COPY ./external/ ./external
RUN go mod download

COPY ./cmd/*.go  ./
RUN go build -o /taxapp
EXPOSE 1401
CMD ["/taxapp"]