.PHONY: gen
gen:
	@protoc ngrpc/calculator.proto --go_out=plugins=grpc:.

.PHONY: server
server:
	@go run server/server.go

.PHONY: client
client:
	@go run client/client.go

.PHONY: ssl
ssl:
	@cd ./ssl; ./gen_ssl.sh; cd ..;