# user-account-auth-system
User Account CRUD including login and logout functions as well as session validation.

- To build target: 
-      go build -o bin/api ./cmd/api/

-  To run target:
-      go run ./cmd/api 

- To lint code
-      golangci-lint run ./... 

- To check code complexity
-      gocyclo -over 7 . 

- To delete build:
-      rm -rf bin 

## Makefile commands for this are:
```
make build-api
make run
make lint
make cyclomatic 
make clean
```


