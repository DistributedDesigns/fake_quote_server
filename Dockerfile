FROM golang:1.7.4-alpine

ADD server.go $GOPATH/src/fake-quoteserv/server.go

RUN go build \
  -o $GOPATH/bin/fake-quoteserv \
  $GOPATH/src/fake-quoteserv/server.go

EXPOSE 4443

CMD $GOPATH/bin/fake-quoteserv
