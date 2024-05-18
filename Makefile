.PHONY: deploy-prober deploy-server ui client server build

BINARY_NAME = xmr-nodes

build: client server

ui:
	go generate ./...

client:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/${BINARY_NAME}-client-linux-amd64
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/${BINARY_NAME}-client-linux-arm64

server: ui
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -tags server -o bin/${BINARY_NAME}-server-linux-amd64
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -tags server -o bin/${BINARY_NAME}-server-linux-arm64

clean:
	go clean
	rm -rfv ./bin
	rm -rf ./frontend/build

# Deploying new binary file to server and probers host
# The deploy-* command doesn't build the binary file, so you need to run `make build` first.
# And make sure the inventory and deploy-*.yml file is properly configured.

deploy-server:
	ansible-playbook -i ./tools/ansible/inventory.ini -l server ./tools/ansible/deploy-server.yml -K

deploy-prober:
	ansible-playbook -i ./tools/ansible/inventory.ini -l prober ./tools/ansible/deploy-prober.yml -K
