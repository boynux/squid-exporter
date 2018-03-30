.PHONY: all clean

all: test build

EXE = ./bin/squid-exporter
SRC = $(shell find ./ -type f -name '*.go')

$(EXE): $(SRC)
	go build -o $(EXE) .

test:
	go test -v ./...

build: $(EXE)

docker:
	docker build -t squid-exporter .

clean:
	rm -f $(EXE)

