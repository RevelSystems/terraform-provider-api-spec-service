HOSTNAME=registry.terraform.io
NAMESPACE=RevelSystems
NAME=api-spec-service
BINARY=terraform-provider-${NAME}
VERSION=0.1.0
OS_ARCH=darwin_$(shell uname -m)

default: install

build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

release:
	goreleaser release --rm-dist --snapshot --skip-publish  --skip-sign