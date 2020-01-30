.PHONY: all fmt vet test build docker clean

all: fmt vet test build

EXE = ./bin/squid-exporter
SRC = $(shell find . -type f -name '*.go')
VERSION ?= $(shell cat VERSION)
REVISION = $(shell git rev-parse HEAD)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
LDFLAGS = -extldflags "-s -w -static" \
		  -X github.com/prometheus/common/version.Version=$(VERSION) \
		  -X github.com/prometheus/common/version.Revision=$(REVISION) \
		  -X github.com/prometheus/common/version.Branch=$(BRANCH)

$(EXE): $(SRC)
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '$(LDFLAGS)' -o $(EXE) .

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test -v ./...

build: $(EXE)

docker:
	docker build -t squid-exporter .

clean:
	rm -f $(EXE)

