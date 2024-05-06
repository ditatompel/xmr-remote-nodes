.PHONY: deploy-server ui build linux-amd64 linux-arm64

BINARY_NAME = xmr-nodes

# Deploy server
# To use this, make sure the inventory and deploy-server.yml file is properly configured
deploy-server: build
	ansible-playbook -i ./tools/ansible/inventory.ini -l server ./tools/ansible/deploy-server.yml -K

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
