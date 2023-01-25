PROJECT_ROOT ?= $(PWD)

.PHONY: managerService
managerService:
	@go run ./manager

.PHONY: clientService
clientService:
	@go run ./client/

.PHONY: kafka-producer
kafka-producer:
	@go run ./producer

.PHONY: protobuf
protobuf:
	#brew install protoc-gen-go
	#brew install protoc-gen-go-grpc
	protoc --go_out=$(PROJECT_ROOT)/proto/manager --go-grpc_out=$(PROJECT_ROOT)/proto/manager ./proto/manager/server.proto

.PHONY: vendor
vendor:
	@go mod tidy; go mod vendor

.PHONY: kafka-up
kafka-up:
	@cd $(PROJECT_ROOT)/docker/; docker-compose up -d
	docker exec -it broker kafka-topics --bootstrap-server broker:9092 --create --topic test-consume-once --partitions 3

.PHONY: kafka-down
kafka-down:
	@cd $(PROJECT_ROOT)/docker/; docker-compose down --remove-orphans
