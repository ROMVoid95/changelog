DIST := dist
GO ?= go

ifneq ($(DRONE_TAG),)
	SHORT_VERSION ?= $(subst v,,$(DRONE_TAG))
	LONG_VERSION ?= $(SHORT_VERSION)
else
	SHORT_VERSION ?= $(shell git describe --tags --always --abbrev=0 | sed 's/-/+/' | sed 's/^v//')
	LONG_VERSION ?= $(shell git describe --tags --always | sed 's/-/+/' | sed 's/^v//')
endif

LDFLAGS := $(LDFLAGS) -X "main.Version=$(LONG_VERSION)"

.PHONY: build
build: generate
	$(GO) build -ldflags '-s -w $(LDFLAGS)'
	
.PHONY: generate
generate:
	$(GO) generate ./...

.PHONY: lint
lint:
	@hash golangci-lint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		export BINARY="golangci-lint"; \
		curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.23.1; \
	fi
	golangci-lint run --timeout 5m

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: release
release: release-dirs check-xgo release-windows release-linux release-darwin release-compress release-check

.PHONY: release-dirs
release-dirs:
	mkdir -p $(DIST)/

.PHONY: check-xgo
check-xgo:
	@hash xgo > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u src.techknowlogick.com/xgo; \
	fi

.PHONY: release-linux
release-linux:
	xgo -dest $(DIST)/ -targets 'linux/amd64,linux/386,linux/arm-5,linux/arm-6,linux/arm64,linux/mips64le,linux/mips,linux/mipsle' -out changelog-$(SHORT_VERSION) .

.PHONY: release-windows
release-windows:
	xgo -dest $(DIST)/ -targets 'windows/*' -out changelog-$(SHORT_VERSION) .

.PHONY: release-darwin
release-darwin:
	xgo -dest $(DIST)/ -targets 'darwin/*' -out changelog-$(SHORT_VERSION) .

.PHONY: release-check
release-check:
	cd $(DIST)/; for file in `find . -type f -name "*"`; do echo "checksumming $${file}" && $(SHASUM) `echo $${file} | sed 's/^..//'` > $${file}.sha256; done;

.PHONY: release-compress
release-compress:
	@hash gxz > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/ulikunitz/xz/cmd/gxz; \
	fi
	cd $(DIST)/; for file in `find . -type f -name "*"`; do echo "compressing $${file}" && gxz -k -9 $${file}; done;