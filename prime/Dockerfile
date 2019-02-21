FROM golang:1.11.1

RUN go get -u -v github.com/kardianos/govendor

RUN mkdir -p /go/src/github.com/tasbouz/faas/prime
WORKDIR /go/src/github.com/tasbouz/faas/prime

EXPOSE 8080

CMD govendor init && govendor fetch -v +m && go run /go/src/github.com/tasbouz/faas/prime/primed.go