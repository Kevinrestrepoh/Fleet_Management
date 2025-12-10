PROTO_DIR=./proto

VEHICLES_PB=./go-simulator/pb
INGEST_PB=./go-ingestion/pb

all: proto

vehicles:
	protoc \
	  --proto_path=$(PROTO_DIR) \
	  --go_out=$(VEHICLES_PB) \
	  --go-grpc_out=$(VEHICLES_PB) \
	  $(PROTO_DIR)/*.proto

ingest:
	protoc \
	  --proto_path=$(PROTO_DIR) \
	  --go_out=$(INGEST_PB) \
	  --go-grpc_out=$(INGEST_PB) \
	  $(PROTO_DIR)/*.proto

clean:
	rm -rf $(VEHICLES_PB)/*
	rm -rf $(INGEST_PB)/*
