FROM golang:alpine

ENV GO111MODULE on
ENV GOFLAGS -mod=vendor

WORKDIR /go/src/srv
COPY . .

RUN go build -o $GOPATH/bin/app

HEALTHCHECK \
  CMD app healthcheck

EXPOSE 8080
ENTRYPOINT ["app"]
