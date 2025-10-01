.PHONY: all compile test clean

all: compile test

# ====================================================================================
# Go
# ====================================================================================

test:
	@echo "--> Running go tests"
	go test -v -race ./...

# ====================================================================================
# Protobuf
# ====================================================================================

compile:
	@echo "--> Compiling protobuf files"
	protoc api/v1/*.proto \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative

clean:
	@echo "--> Cleaning generated protobuf files"
	rm -f api/v1/*.pb.go
