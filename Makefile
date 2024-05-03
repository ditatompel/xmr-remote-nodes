.PHONY: ui build linux64

BINARY_NAME = xmr-nodes

build: ui linux64

ui:
	go generate ./...

linux64:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-static-linux-amd64

clean:
	go clean
	rm -rfv ./bin
	rm -rf ./frontend/build
