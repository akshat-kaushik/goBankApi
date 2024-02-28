build:
	@go build -o bin/pricefetcher

run: build
	@./bin/pricefetcher

proto:
	@protoc --go-grpc_out=. --go_out=. *.proto
	