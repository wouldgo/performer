SHELL := /bin/sh
OUT := $(shell pwd)/_out
BUILDARCH := $(shell uname -m)
GCC := $(OUT)/$(BUILDARCH)-linux-musl-cross/bin/$(BUILDARCH)-linux-musl-gcc
LD := $(OUT)/$(BUILDARCH)-linux-musl-cross/bin/$(BUILDARCH)-linux-musl-ld
VERSION := 0.0.1

test: deps
	rm -Rf _out/.coverage;
	go test -run W -timeout 120s -coverprofile=_out/.coverage -cover -v ./...;
	go tool cover -html=_out/.coverage;

performer: deps
	go run cmd/performer/*.go

compile-performer: deps
	CGO_ENABLED=0 \
	go build \
		-ldflags='-extldflags=-static' \
		-a -o _out/performer cmd/performer/*.go

deps: musl
	go mod tidy -v
	go mod download

musl:
	if [ ! -d "$(OUT)/$(BUILDARCH)-linux-musl-cross" ]; then \
		(cd $(OUT); curl -LOk https://musl.cc/$(BUILDARCH)-linux-musl-cross.tgz) && \
		tar zxf $(OUT)/$(BUILDARCH)-linux-musl-cross.tgz -C $(OUT); \
	fi

clean:
	rm -Rf $(OUT) $(BINARY_NAME)
	mkdir -p $(OUT)
	touch $(OUT)/.keep
