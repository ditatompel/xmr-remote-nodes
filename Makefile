BINARY_NAME = xmr-nodes

.PHONY: build
build: client server

.PHONY: ui
ui:
	go generate ./...

.PHONY: client
client:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/${BINARY_NAME}-client-linux-amd64
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/${BINARY_NAME}-client-linux-arm64

.PHONY: server
server: ui
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -tags server -o bin/${BINARY_NAME}-server-linux-amd64
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -tags server -o bin/${BINARY_NAME}-server-linux-arm64

.PHONY: clean
clean:
	go clean
	rm -rfv ./bin
	rm -rf ./frontend/build

# Deploying new binary file to server and probers host
# The deploy-* command doesn't build the binary file, so you need to run `make build` first.
# And make sure the inventory and deploy-*.yml file is properly configured.

.PHONY: deploy-server
deploy-server:
	ansible-playbook -i ./deployment/ansible/inventory.ini -l server ./deployment/ansible/deploy-server.yml -K

.PHONY: deploy-prober
deploy-prober:
	ansible-playbook -i ./deployment/ansible/inventory.ini -l prober ./deployment/ansible/deploy-prober.yml -K
