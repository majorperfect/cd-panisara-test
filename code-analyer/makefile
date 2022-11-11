generate:
	protoc -I ./proto \
    --go_out ./proto --go_opt paths=source_relative \
    --go-grpc_out ./proto --go-grpc_opt paths=source_relative \
    --grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative \
    ./proto/code-analyzer/code-analyzer.proto; \
    sed -i "" -e "s/,omitempty//g" ./proto/code-analyzer/*.pb.go; \
    statik -m -f -src ./proto/OpenAPI/
