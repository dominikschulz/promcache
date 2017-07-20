FROM golang:1.8-alpine3.6 as builder

ADD . /go/src/github.com/dominikschulz/promcache
WORKDIR /go/src/github.com/dominikschulz/promcache

RUN go install

FROM alpine:3.6

COPY --from=builder /go/bin/promcache /usr/local/bin/promcache
CMD [ "/usr/local/bin/promcache" ]
EXPOSE 9091 9092
