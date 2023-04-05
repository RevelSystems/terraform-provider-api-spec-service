GOARCH=$(shell go env GOARCH)
INSTALL_PATH=~/.terraform.d/plugins/localhost/providers/api-spec-service/0.0.1/darwin_$(GOARCH)

build:
	mkdir -p $(INSTALL_PATH)
	go build -o $(INSTALL_PATH)/terraform-provider-api-spec-service main.go

dev: build
	rm ./examples/local/.terraform.lock.hcl || true
	cd ./examples/local && terraform init
	cd ./examples/local && terraform destroy
	cd ./examples/local && TF_LOG=TRACE terraform apply -auto-approve
