FROM golang:alpine as builder

ADD . /go/src/github.com/dominikschulz/promcache
WORKDIR /go/src/github.com/dominikschulz/promcache

RUN go install

FROM alpine:latest

COPY --from=builder /go/bin/promcache /usr/local/bin/promcache
CMD [ "/usr/local/bin/promcache" ]
EXPOSE 9091 9092
