.PHONY: build image

build:
	go build -o bin/runtime-check cmd/main.go

image:
	docker build -t runtime-check .
