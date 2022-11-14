generate-mocks:
	rm -rf ./internal/mocks
	mockery --all --dir=./internal --keeptree --recursive --output ./internal/mocks
	rm -rf ./pkg/mocks
	mockery --all --dir=./pkg --recursive --output ./pkg/mocks