PROJECT_ROOT ?= $(PWD)

managerService:
	go run ./manager

clientService:
	go run ./client/

protobuf:
	#brew install protoc-gen-go
	#brew install protoc-gen-go-grpc
	protoc --go_out=$(PROJECT_ROOT)/proto/manager --go-grpc_out=$(PROJECT_ROOT)/proto/manager ./proto/manager/server.proto

vendor:
	go mod tidy; go mod vendor
