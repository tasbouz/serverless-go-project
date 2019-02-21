FROM golang:1.11.1

RUN go get -u -v github.com/kardianos/govendor

RUN mkdir -p /go/src/github.com/tasbouz/faas/fibonacci
WORKDIR /go/src/github.com/tasbouz/faas/fibonacci

EXPOSE 8080

CMD govendor init && govendor fetch -v +m && go run /go/src/github.com/tasbouz/faas/fibonacci/fibonaccid.go