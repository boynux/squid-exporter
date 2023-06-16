FROM golang:alpine as builder

WORKDIR /go/src/github.com/boynux/squid-exporter
COPY . .

# Compile the binary statically, so it can be run without libraries.
RUN CGO_ENABLED=0 GOOS=linux go install -a -ldflags '-extldflags "-s -w -static"' .

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /go/bin/squid-exporter /usr/local/bin/squid-exporter

# Allow /etc/hosts to be used for DNS
COPY --from=builder /etc/nsswitch.conf /etc/nsswitch.conf

EXPOSE 9301

ENTRYPOINT ["/usr/local/bin/squid-exporter"]
