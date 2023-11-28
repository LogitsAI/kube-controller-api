tools:
	@# Check for protoc and advise using Homebrew if not found.
	@which -s protoc || echo "Please install protobuf compiler: brew install protobuf"
	@# Install Go plugins for the protobuf compiler.
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	@which -s protoc-gen-go || echo "Please add $(shell go env GOPATH)/bin to your PATH"
	@# Install Python plugins for the protobuf compiler.
	@python3 -m pip install grpcio grpcio-tools

build: tools
	protoc -I=controller \
		--go_out=../controllerpb --go_opt=paths=source_relative \
		--go-grpc_out=../controllerpb --go-grpc_opt=paths=source_relative \
		controller/*.proto
	python3 -m grpc_tools.protoc -Icontroller \
	    --python_out=../client/python --pyi_out=../client/python \
		--grpc_python_out=../client/python \
		controller/*.proto
