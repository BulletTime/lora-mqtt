.PHONY: default all macos linux windows build clean
VERSION := 0.1.2
COMMIT := $(shell git describe --always)
GOOS ?= darwin
GOARCH ?= amd64
GOPATH ?= $(HOME)/go/
BUILD_DATE = `date -u +%Y-%m-%dT%H:%M.%SZ`
BUILD_NAME = lora-mqtt
MAIN_FILE = main.go

.SILENT:
default: clean build

all: clean macos linux windows

macos: build

linux:
	env GOOS=linux GOARCH=amd64 $(MAKE) build
	env GOOS=linux GOARCH=arm $(MAKE) build

windows:
	env GOOS=windows GOARCH=amd64 BINEXT=.exe $(MAKE) build

build:
	echo "[===] Build for $(GOOS) $(GOARCH) [===]"
	mkdir -p build
	echo "[GO BUILD] $(MAIN_FILE)"
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -ldflags "-X main.version=$(VERSION) -X main.build=$(COMMIT) -X main.buildDate=$(BUILD_DATE)" -o build/$(BUILD_NAME)-$(GOOS)-$(GOARCH)$(BINEXT) $(MAIN_FILE)

clean:
	echo "[===] Cleaning up workspace [===]"
	rm -rf build
	rm -rf lora-mqtt.log
