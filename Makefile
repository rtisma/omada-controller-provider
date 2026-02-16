.PHONY: build install test testacc clean fmt docs

HOSTNAME=registry.terraform.io
NAMESPACE=your-org
NAME=omada
BINARY=terraform-provider-${NAME}
VERSION=0.1.0
OS_ARCH=darwin_arm64

default: build

build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	cp ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:
	go test -v ./...

testacc:
	TF_ACC=1 go test -v ./... -timeout 120m

clean:
	rm -f ${BINARY}
	rm -rf ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}

fmt:
	go fmt ./...
	terraform fmt -recursive ./examples/

lint:
	golangci-lint run

docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate

deps:
	go mod tidy
	go mod verify

validate:
	terraform fmt -check -recursive ./examples/
	go vet ./...

release-check:
	@echo "Checking release readiness..."
	@test -n "$(VERSION)" || (echo "VERSION is not set. Use: make release-check VERSION=v0.1.0" && exit 1)
	@echo "Building for release $(VERSION)..."
	@go build -o terraform-provider-omada
	@echo "✓ Build successful"
	@echo "✓ Ready to release $(VERSION)"
	@echo ""
	@echo "To create the release, run:"
	@echo "  git tag -a $(VERSION) -m 'Release $(VERSION)'"
	@echo "  git push origin $(VERSION)"

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the provider binary"
	@echo "  install       - Build and install the provider locally"
	@echo "  test          - Run unit tests"
	@echo "  testacc       - Run acceptance tests (requires TF_ACC=1)"
	@echo "  clean         - Remove built binaries and local plugins"
	@echo "  fmt           - Format Go and Terraform files"
	@echo "  lint          - Run golangci-lint"
	@echo "  docs          - Generate documentation with tfplugindocs"
	@echo "  deps          - Tidy and verify dependencies"
	@echo "  validate      - Validate Terraform files and Go code"
	@echo "  release-check - Check if ready for release (use VERSION=v0.1.0)"
