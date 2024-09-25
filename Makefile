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
