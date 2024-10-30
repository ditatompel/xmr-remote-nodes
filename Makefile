BINARY_NAME = xmr-nodes

# These build are modified version of rclone's Makefile
# https://github.com/rclone/rclone/blob/master/Makefile
VERSION := $(shell cat VERSION)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
# Last tag on this branch (eg. v1.0.0)
LAST_TAG := $(shell git describe  --tags --abbrev=0)
# Tag of the current commit, if any. If this is not "" then we are building a release
RELEASE_TAG := $(shell git tag -l --points-at HEAD)
# If we are working on a release, override branch to main
ifdef RELEASE_TAG
	BRANCH := main
	LAST_TAG := $(shell git describe --abbrev=0 --tags $(VERSION)^)
endif
# Make version suffix -beta.NNNN.CCCCCCCC (N=Commit number, C=Commit)
VERSION_SUFFIX := -beta.$(shell git rev-list --count HEAD).$(shell git show --no-patch --no-notes --pretty='%h' HEAD)
TAG_BRANCH := .$(BRANCH)
# If building HEAD or master then unset TAG_BRANCH
ifeq ($(subst HEAD,,$(subst main,,$(BRANCH))),)
	TAG_BRANCH :=
endif
# TAG is current version + commit number + commit + branch
TAG := $(VERSION)$(VERSION_SUFFIX)$(TAG_BRANCH)
ifdef RELEASE_TAG
	TAG := $(RELEASE_TAG)
endif
# end modified rclone's Makefile

BUILD_LDFLAGS := -s -w -X github.com/ditatompel/xmr-remote-nodes/internal/config.Version=$(TAG)

# This called from air cmd (see .air.toml)
.PHONY: dev
dev: templ tailwind
	go build -ldflags="$(BUILD_LDFLAGS)" -tags server -o ./tmp/main .

.PHONY: build
build: client server

.PHONY: client
client:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build \
		-ldflags="$(BUILD_LDFLAGS)"                \
		-o bin/${BINARY_NAME}-client-linux-amd64
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build \
		-ldflags="$(BUILD_LDFLAGS)"                \
		-o bin/${BINARY_NAME}-client-linux-arm64

.PHONY: server
server: prepare templ tailwind
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build \
		-ldflags="$(BUILD_LDFLAGS)" -tags server   \
		-o bin/${BINARY_NAME}-server-linux-amd64
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build \
		-ldflags="$(BUILD_LDFLAGS)" -tags server   \
		-o bin/${BINARY_NAME}-server-linux-arm64

.PHONY: prepare
prepare:
	bun install --frozen-lockfile
	@mkdir -p ./internal/handler/views/assets/js
	cp ./node_modules/htmx.org/dist/htmx.min.js ./internal/handler/views/assets/js

# Compile template
.PHONY: templ
templ:
	@echo "Compiling Templ template..."
	templ generate

.PHONY: tailwind
tailwind:
	mkdir -p ./internal/handler/views/assets/css
	@echo "Compiling TailwindCSS..."
	bun tailwindcss -i ./internal/handler/views/src/css/main.css \
		-o ./internal/handler/views/assets/css/main.min.css \
		-c ./tailwind.config.js \
		--minify
	bun build ./internal/handler/views/src/js/main.js --minify \
		--outfile ./internal/handler/views/assets/js/main.min.js

.PHONY: clean
clean:
	go clean
	rm -rfv ./bin
	rm -rfv ./tmp/main
	rm -rf ./internal/handler/views/*_templ.go
	rm -rf ./internal/handler/views/assets/css/

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: bench
bench:
	go test ./... -bench=. -benchmem -run=^#

# Deploying new binary file to server and probers host
# The deploy-* command doesn't build the binary file, so you need to run `make build` first.
# And make sure the inventory and deploy-*.yml file is properly configured.
.PHONY: deploy-server
deploy-server:
	ansible-playbook -i ./deployment/ansible/inventory.ini \
		-l server ./deployment/ansible/deploy-server.yml -K

.PHONY: deploy-prober
deploy-prober:
	ansible-playbook -i ./deployment/ansible/inventory.ini \
		-l prober ./deployment/ansible/deploy-prober.yml -K
