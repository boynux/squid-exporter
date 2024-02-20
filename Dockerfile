FROM golang:1.19.1-alpine as build

ARG TARGETPLATFORM
RUN echo "Building for ${TARGETPLATFORM}"

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-s -w -static"' -o /squid-exporter .


FROM gcr.io/distroless/static:nonroot as final

LABEL org.opencontainers.image.title="Squid Exporter"
LABEL org.opencontainers.image.description="This is a Docker image for Squid Prometheus Exporter."
LABEL org.opencontainers.image.source="https://github.com/boynux/squid-exporter/"
LABEL org.opencontainers.image.licenses="MIT"

ENV SQUID_EXPORTER_LISTEN="0.0.0.0:9301"

COPY --from=build /squid-exporter /usr/local/bin/squid-exporter
# Allow /etc/hosts to be used for DNS
COPY --from=build /etc/nsswitch.conf /etc/nsswitch.conf

EXPOSE 9301

ENTRYPOINT ["/usr/local/bin/squid-exporter"]
