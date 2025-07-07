
build: 
	@go build -o bin/gobank

run: build
	@./bin/gobank


test: 
	@go test -v ./...	 


# build_proto:
# 	protoc --go_out=. --go_opt=paths=source_relative \
#        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
#        protos/shop.proto

