.DEFAULT_GOAL := build

.PHONY: swag
swag:
	swag init --parseDependency --dir internal/app -g apiserver.go -o api/swagger

.PHONY: build
build:
	go mod download && CGO_ENABLED=0 go build -o ./.bin/apiserver ./cmd/apiserver

.PHONY: gen
gen:
	go generate ./...

test:
	go test -v ./...

.PHONY: cover
cover:
	go test -count=1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out