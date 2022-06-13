FROM alpine:3.6 as certs
RUN apk add -U --no-cache ca-certificates make

FROM golang:1.18.2-alpine as build
RUN apk add -U --no-cache make
WORKDIR /workspace
COPY . .
RUN make install

FROM alpine:3.6
EXPOSE 3000
ENV GODEBUG netdns=go
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/bin/drone-promote-auth /go/bin/

ENTRYPOINT ["/go/bin/drone-promote-auth"]
