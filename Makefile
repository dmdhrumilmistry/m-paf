bump:
	@go mod tidy
	@go get -u ./...

test:
	@go test -v ./...