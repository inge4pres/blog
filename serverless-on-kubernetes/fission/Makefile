build:
	docker run --rm -v "$$(pwd)":/usr/src/fissionapp \
		-v "$$(pwd)/vendor/github.com":/go/src/github.com \
		-w /usr/src/fissionapp golang:1.9 go build -buildmode=plugin -o function.so 

.PHONY: build