.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	bash scripts/build

.PHONY: docs
docs:
	bash scripts/docs