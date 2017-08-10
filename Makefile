.PHONY: all clean

all: test build

BUILD_PATH = ./cmd/squid-exporter
EXE = ./bin/squid-exporter
SRC = $(shell find ./ -type f -name '*.go')

$(EXE): $(SRC)
	go build -o $(EXE) $(BUILD_PATH)

test:
	go test -v ./...

build: $(EXE)

clean:
	rm -f $(EXE)

