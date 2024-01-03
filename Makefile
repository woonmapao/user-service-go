build:
	@go build -o bin/user-service

run: build
	@./bin/user-service

test:
	 @go test -v ./...




