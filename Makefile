.PHONY: ui build linux-amd64 linux-arm64

BINARY_NAME = xmr-nodes

build: ui linux-amd64 linux-arm64

ui:
	go generate ./...

linux-amd64:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-static-linux-amd64

linux-arm64:
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -o bin/${BINARY_NAME}-static-linux-arm64

clean:
	go clean
	rm -rfv ./bin
	rm -rf ./frontend/build
