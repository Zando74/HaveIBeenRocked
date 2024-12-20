PROTOC = protoc
PROTO_DIR = internal/application/grpc_proto
OUT_DIR = .

DEV_CONFIG_PATH = ./internal/config/config.yml
TEST_CONFIG_PATH = ../config/config.yml

.PHONY: proto
protoc:
	$(PROTOC) --go_out=$(OUT_DIR) \
		--go-grpc_out=$(OUT_DIR) \
		--go-grpc_opt=paths=source_relative \
		--go_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

.PHONY: install
install:
	go mod download

.PHONY: build
build:
	go build -o main cmd/main.go

.PHONY: run-dev
run-dev:
	CONFIG_PATH=$(DEV_CONFIG_PATH) go run cmd/main.go

.PHONY: test-integration
test-integration:
	CONFIG_PATH=$(TEST_CONFIG_PATH) go test -v ./internal/integration/...

.PHONY: test-integration-ci
test-integration-ci:
	CONFIG_PATH=$(TEST_CONFIG_PATH) go test ./internal/integration/...