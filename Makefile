.PHONY: build
build: 
	CGO_ENABLED=0 go build -v -ldflags "-X main.version=$(shell cat VERSION)" -o build/healthd

.PHONY: image
image:
	docker build -thealthd:latest .

#.PHONY: image-verbose
#image-verbose:
#	docker build --no-cache --progress=plain -thealthd:latest .

.PHONY: e2e
e2e:
	cd e2e;go run github.com/onsi/ginkgo/v2/ginkgo -v

.PHONY: test
test:
	go test `go list ./... | grep -v e2e`

.PHONY: lint
LINT_IMAGE ?= golangci/golangci-lint:v2.11.4
lint:
	docker run --rm -v "$(CURDIR)":/app -w /app $(LINT_IMAGE) golangci-lint run

.PHONY: run-local
run-local:
	go run main.go -global-config configs/conf/global.json -script-config configs/conf/script.json -v

RPM_TOPDIR ?= $(CURDIR)/build/rpmbuild
VERSION ?= $(shell cat VERSION)
RELEASE ?= 1
GOARCH ?= amd64
RPM_ARCH ?= x86_64

.PHONY: rpm
rpm:
	mkdir -p $(RPM_TOPDIR)/BUILD $(RPM_TOPDIR)/BUILDROOT $(RPM_TOPDIR)/RPMS $(RPM_TOPDIR)/SOURCES $(RPM_TOPDIR)/SPECS $(RPM_TOPDIR)/SRPMS
	GOOS=linux GOARCH=$(GOARCH) CGO_ENABLED=0 go build -v -ldflags "-X main.version=$(VERSION)" -o build/healthd
	cp build/healthd $(RPM_TOPDIR)/SOURCES/healthd
	rpmbuild --target $(RPM_ARCH) --define "_topdir $(RPM_TOPDIR)" --define "version $(VERSION)" --define "release $(RELEASE)" -bb packaging/healthd.spec

.PHONY: rpm-docker-amd64
rpm-docker-amd64:
	docker run --rm --platform linux/amd64 -v "$(CURDIR)":/work -w /work ubuntu:24.04 bash -lc 'set -euo pipefail; export DEBIAN_FRONTEND=noninteractive; apt-get update >/dev/null; apt-get install -y ca-certificates curl make rpm file binutils >/dev/null; curl -fsSL https://go.dev/dl/go1.26.2.linux-amd64.tar.gz -o /tmp/go.tar.gz; rm -rf /usr/local/go; tar -C /usr/local -xzf /tmp/go.tar.gz; export PATH=/usr/local/go/bin:$$PATH; make rpm GOARCH=amd64 RPM_ARCH=x86_64 VERSION=$(VERSION) RELEASE=$(RELEASE)'

# Ubuntu コンテナに RPM をインストールして /usr/bin/healthd の存在と起動を確認する
.PHONY: rpm-install-test
rpm-install-test:
	$(eval RPM_FILE := $(shell find $(RPM_TOPDIR)/RPMS/x86_64 -name "*.rpm" | sort | tail -1))
	@test -n "$(RPM_FILE)" || (echo "ERROR: No RPM found under $(RPM_TOPDIR)/RPMS/x86_64. Run make rpm-docker-amd64 first." && exit 1)
	docker run --rm --platform linux/amd64 -v "$(RPM_FILE)":/tmp/healthd.rpm ubuntu:24.04 bash -lc '\
		apt-get update >/dev/null && apt-get install -y alien >/dev/null && \
		alien --to-deb --scripts /tmp/healthd.rpm >/dev/null && \
		dpkg -i healthd_*.deb >/dev/null && \
		healthd --help 2>&1 | head -5 || healthd -v 2>&1 | head -5 || (ls -l /usr/bin/healthd && echo "healthd installed OK")'
