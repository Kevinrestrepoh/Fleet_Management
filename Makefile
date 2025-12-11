PROTO_DIR=./proto

VEHICLES_PB=./vehicle-simulator

generate:
	protoc \
	  --proto_path=$(PROTO_DIR) \
	  --go_out=$(VEHICLES_PB) \
	  --go-grpc_out=$(VEHICLES_PB) \
	  $(PROTO_DIR)/*.proto
