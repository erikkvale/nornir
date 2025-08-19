PROTOC = protoc
PROTO_DIR = proto
PROTO_FILES = $(wildcard $(PROTO_DIR)/*/*.proto)

REGISTRY = local
GATEWAY_IMAGE = nornir-gateway
WORKER_IMAGE = nornir-worker
TAG = latest

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

proto-clean:
	rm -f $(PROTO_DIR)/*/*.pb.go $(PROTO_DIR)/*/*_grpc.pb.go

test:
	go test ./...

build:
	docker build -f gateway-service/Dockerfile -t $(GATEWAY_IMAGE):$(TAG) .
	docker build -f worker-service/Dockerfile -t $(WORKER_IMAGE):$(TAG) .

helm-install:
	helm upgrade --install nornir ./charts/nornir

helm-install-local: build
	helm upgrade --install nornir ./charts/nornir \
		--set gateway.image.repository=$(GATEWAY_IMAGE) \
		--set gateway.image.tag=$(TAG) \
		--set worker.image.repository=$(WORKER_IMAGE) \
		--set worker.image.tag=$(TAG)


deploy: helm-install-local


.PHONY: proto proto-clean test build helm-install helm-install-local deploy 
