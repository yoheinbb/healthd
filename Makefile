.PHONY: build
build: 
	CGO_ENABLED=0 go build -v -o build/healthd

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
lint:
	golangci-lint run

.PHONY: run-local
run-local:
	go run main.go -global-config configs/conf/global.json -script-config configs/conf/script.json -v
