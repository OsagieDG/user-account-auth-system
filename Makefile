build-api:
	@go build -o bin/api ./cmd/api/

run: build-api
	@./bin/api

lint:
	@golangci-lint run ./...

cyclomatic:
	@gocyclo -over 7 . 

clean: 
	@rm -rf bin



