.PHONY: build dev test fmt lint

build:
	go build ./...

dev:
	cd packages/web && pnpm dev

test:
	go test ./...

fmt:
	go fmt ./...
	cd packages/web && pnpm fmt

lint:
	go vet ./...
	cd packages/web && pnpm lint
