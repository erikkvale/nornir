PROTOC = protoc
PROTO_DIR = proto
PROTO_FILES = $(wildcard $(PROTO_DIR)/*/*.proto)

proto:
	@for file in $(PROTO_FILES); do \
		$(PROTOC) \
			--proto_path=$(PROTO_DIR) \
			--go_out=$(PROTO_DIR) \
			--go_opt=paths=source_relative \
			--go-grpc_out=$(PROTO_DIR) \
			--go-grpc_opt=paths=source_relative \
			$$file; \
	done

clean:
	rm -f $(PROTO_DIR)/*/*.pb.go $(PROTO_DIR)/*/*_grpc.pb.go

.PHONY: proto clean
