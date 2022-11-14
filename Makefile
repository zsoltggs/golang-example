protos:
	# PDS
	mkdir -p ./pkg
	protoc --go_out=${GOPATH}/src \
	--go-grpc_out=require_unimplemented_servers=true:${GOPATH}/src \
	./services/port-domain-service/resources/port-domain-service.proto
	# Users
	protoc --go_out=${GOPATH}/src \
	--go-grpc_out=require_unimplemented_servers=true:${GOPATH}/src \
	./services/users/resources/users.proto

lint:
	golangci-lint run


.PHONY: build
build:
	GOOS=linux go build -o ./services/client-api-service/build/cas ./services/client-api-service/cmd/client-api-service/main.go
	GOOS=linux go build -o ./services/port-domain-service/build/pds ./services/port-domain-service/cmd/port-domain-service/main.go
	GOOS=linux go build -o ./services/users/build/users ./services/users/cmd/users/main.go
	GOOS=linux go build -o ./services/validation-service/build/validationsvc ./services/validation-service/cmd/validation-service/main.go

.PHONY: clean
clean:
	rm -rf build
	rm -rf ./services/client-api-service/build
	rm -rf ./services/port-domain-service/build
	rm -rf ./services/users/build
	rm -rf ./services/validation-service/build

.PHONY: build-docker
build-docker: build
	docker build --no-cache -t cas:local ./services/client-api-service --build-arg SERVICE=client-api-service
	docker build --no-cache -t pds:local ./services/port-domain-service --build-arg SERVICE=port-domain-service
	docker build --no-cache -t users:local ./services/users --build-arg SERVICE=users
	docker build --no-cache -t validationsvc:local ./services/validation-service --build-arg SERVICE=validation-service

.PHONY: generate-mocks
generate-mocks:
	for service in validation-service; do \
		echo ./srv/$$service ; \
		$(MAKE) -C ./services/$$service generate-mocks ; \
	done