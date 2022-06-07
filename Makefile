.phony: fmt all lib tst cov vet lnt

build_lib = build

all: lib

fmt:
	go fmt ./pkg/...

lib: fmt
	go build -mod=mod

tst: fmt
	go test -coverprofile=coverage.out ./pkg/...

cov: test
	go tool cover -html=coverage.out


vet: all
	go vet ./pkg/...

lnt: fmt
	golangci-lint run -v --timeout 5m
