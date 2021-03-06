# Author:       nghiatc
# Email:        congnghia0609@gmail.com

.PHONY: deps
deps:
	@./deps.sh

.PHONY: gen
gen:
	@protoc example/ngrpc/calculator.proto --go_out=plugins=grpc:.

.PHONY: server
server:
	@go run main.go

.PHONY: client
client:
	@go run cli.go

.PHONY: ssl
ssl:
	@cd ./ssl; ./gen_ssl.sh; cd ..;
