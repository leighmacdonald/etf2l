VERSION=v0.0.1
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_FLAGS = -trimpath -ldflags="-s -w -X github.com/leighmacdonald/etf2l/internal/app.BuildVersion=$(VERSION)"
DEBUG_FLAGS = -gcflags "all=-N -l"

vet:
	@go vet . ./...

fmt:
	gci write . --skip-generated -s standard -s default
	gofumpt -l -w .

bump_deps:
	go get -u ./...

test:
	@go test $(GO_FLAGS) -race -cover . ./...

check_deps:
	go install github.com/daixiang0/gci@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
	go install honnef.co/go/tools/cmd/staticcheck@latest

check:
	golangci-lint run --timeout 3m ./...
	staticcheck -go 1.21 ./...
