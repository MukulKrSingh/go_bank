
build: postgres
	@go build -o bin/gobank

run: build
	@./bin/gobank


test: 
	@go test -v ./...


postgres:
	@docker container start go-bank-db

stopDb:
	@docker container stop go-bank-db	


# build_proto:
# 	protoc --go_out=. --go_opt=paths=source_relative \
#        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
#        protos/shop.proto

