PROTO_DIR=proto
GEN_DIR=api
PACKAGE_NAME=github.com/roshsyed/virtinput-server

.PHONY: generate clean test help

help:
	@echo "Usage:"
	@echo "  make generate    - Generate Go code from proto files"
	@echo "  make clean       - Remove generated files"

generate:
	@mkdir -p $(GEN_DIR)
	protoc --proto_path=$(PROTO_DIR) \
		--go_out=$(GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative \
		$$(find $(PROTO_DIR) -name '*.proto')

clean:
	rm -rf $(GEN_DIR)/*
