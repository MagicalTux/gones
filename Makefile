#!/bin/make
GOROOT:=$(shell PATH="/pkg/main/dev-lang.go.dev/bin:$$PATH" go env GOROOT)
GOPATH:=$(shell $(GOROOT)/bin/go env GOPATH)

.PHONY: test deps

all:
	$(GOPATH)/bin/goimports -w -l .
	$(GOROOT)/bin/go build -v

run: all
	#./gones -ppudebug - nes-test-roms/ppu_vbl_nmi/ppu_vbl_nmi.nes
	./gones -apudebug - nes-test-roms/apu_test/apu_test.nes
	#./gones cur.nes

deps:
	$(GOROOT)/bin/go get -v -t .

test:
	$(GOROOT)/bin/go test -v
