GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

build:
	mkdir -p dist/$(GOOS)-$(GOARCH)
	go build \
		-o dist/$(GOOS)-$(GOARCH)/nats-subject-profiler \
		.

zip-build:
	cd dist/$(GOOS)-$(GOARCH) && zip ../nats-subject-profiler-$(GOOS)-$(GOARCH).zip ./*

dist:
	GOOS=linux GOARCH=amd64 make build zip-build
	GOOS=darwin GOARCH=amd64 make build zip-build
	GOOS=windows GOARCH=amd64 make build zip-build
