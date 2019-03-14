# Not using alpine because it requires gcc to build
# TODO: Use multi stage to lighten up the final image
FROM golang:latest

ENV GO111MODULE on
ENV GOFLAGS -mod=vendor

WORKDIR /go/src/srv
COPY . .

RUN go build -o $GOPATH/bin/app

HEALTHCHECK \
  CMD app healthcheck

EXPOSE 8080
ENTRYPOINT ["app"]
