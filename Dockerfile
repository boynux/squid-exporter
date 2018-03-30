FROM golang

WORKDIR /go/src/github.com/boynux/squid-exporter
COPY . .

# Compile the binary statically, so it can be run without libraries.
RUN CGO_ENABLED=0 GOOS=linux go install -a -ldflags '-extldflags "-static"' .

FROM alpine
COPY --from=0 /go/bin/squid-exporter /usr/local/bin/squid-exporter
ENTRYPOINT /usr/local/bin/squid-exporter
