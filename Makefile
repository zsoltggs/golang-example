protos:
	mkdir -p ./port-domain-service/pkg/generated
	protoc --go_out=./port-domain-service/pkg/generated \
	--go-grpc_out=require_unimplemented_servers=true:./port-domain-service/pkg/generated \
	./port-domain-service/resources/port-domain-service.proto

lint:
	golangci-lint run


.PHONY: build
build:
	GOOS=linux go build -o ./client-api-service/build/cas ./client-api-service/cmd/client-api-service/main.go
	GOOS=linux go build -o ./port-domain-service/build/pds ./port-domain-service/cmd/port-domain-service/main.go

.PHONY: clean
clean:
	rm -rf build
	rm -rf ./client-api-service/build
	rm -rf ./port-domain-service/build

.PHONY: build-docker
build-docker: build
	docker build --no-cache -t cas:local ./client-api-service --build-arg SERVICE=client-api-service
	docker build --no-cache -t pds:local ./port-domain-service --build-arg SERVICE=port-domain-service