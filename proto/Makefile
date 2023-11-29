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
	protoc -I=. \
		--go_out=../controllerpb --go_opt=module="github.com/LogitsAI/kube-controller-api/controllerpb" \
		--go-grpc_out=../controllerpb --go-grpc_opt=module="github.com/LogitsAI/kube-controller-api/controllerpb" \
		kube_controller_api/*.proto
	python3 -m grpc_tools.protoc -I. \
	    --python_out=../python/src --pyi_out=../python/src \
		--grpc_python_out=../python/src \
		kube_controller_api/*.proto
