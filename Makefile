PROTO_DIR=./proto

VEHICLES_PB=./vehicle-simulator
CONTROL_PB=./control-api

generate: control vehicle

control:
	protoc \
	  --proto_path=$(PROTO_DIR) \
	  --go_out=$(CONTROL_PB) \
	  --go-grpc_out=$(CONTROL_PB) \
	  $(PROTO_DIR)/*.proto

vehicle:
	protoc \
	  --proto_path=$(PROTO_DIR) \
	  --go_out=$(VEHICLES_PB) \
	  --go-grpc_out=$(VEHICLES_PB) \
	  $(PROTO_DIR)/*.proto
