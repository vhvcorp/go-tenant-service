.PHONY: proto clean

PROTO_DIR := .
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)
OUT_DIR := ../internal/pb
GOOGLEAPIS_DIR := $(GOPATH)/src/github.com/googleapis/googleapis

proto:
	@echo "Generating protobuf code for tenant-service..."
	@mkdir -p $(OUT_DIR)
	protoc -I$(PROTO_DIR) \
		-I$(GOOGLEAPIS_DIR) \
		--go_out=$(OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(OUT_DIR) \
		--grpc-gateway_opt=paths=source_relative \
		--grpc-gateway_opt=logtostderr=true \
		$(PROTO_FILES)
	@echo "Proto generation complete!"

clean:
	@echo "Cleaning generated files..."
	@rm -rf $(OUT_DIR)/*.pb.go
	@echo "Clean complete!"

install-tools:
	@echo "Installing proto generation tools..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@echo "Tools installed!"

setup-googleapis:
	@echo "Setting up googleapis..."
	@mkdir -p $(GOPATH)/src/github.com/googleapis
	@if [ ! -d "$(GOOGLEAPIS_DIR)" ]; then \
		git clone https://github.com/googleapis/googleapis.git $(GOOGLEAPIS_DIR); \
	fi
	@echo "googleapis ready!"
