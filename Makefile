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

run: stop run-ingestion wait-ingestion run-rest

run-ingestion:
	@cd ingestion-server && cargo run > /dev/null 2>&1 &

wait-ingestion:
	@until nc -z localhost 50051; do sleep 1; done

run-rest:
	@cd control-api && go run . > /dev/null 2>&1 &
	@cd vehicle-simulator && go run . > /dev/null 2>&1 &

	@cd fleet-tui && go run .

stop:
	@-kill $$(lsof -t -i :50051) 2>/dev/null || true
