.PHONY: gen help

PROTO_DIR = pb

gen:
	protoc \
	--proto_path=$(PROTO_DIR) \
	--go_out=$(PROTO_DIR) \
	--go_opt=paths=source_relative \
	--go-grpc_out=$(PROTO_DIR) \
	--go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=$(PROTO_DIR) \
	--grpc-gateway_opt=paths=source_relative \
	hello.proto
	$(shell find $(PROTO_DIR) -iname "*.proto")

help:
	@echo "make gen - 生成pb、grpc以及gateway代码"