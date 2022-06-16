.PHONY: check build image

getdeps:
	mkdir -p $(GOPATH)/bin
	which golangci-lint 1>/dev/null || (echo "Installing golangci-lint" && go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2)

lint: getdeps
	echo "Running $@ check"
	${GOPATH}/bin/golangci-lint cache clean
	${GOPATH}/bin/golangci-lint run --timeout=5m --config ./.golangci.yaml

check: lint

build:
	go build -o bin/runtime-check cmd/main.go

image:
	docker build -t runtime-check .
