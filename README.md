# mother-nature

Install some dependencies:

- golangci-lint: https://golangci-lint.run/usage/install/
- mockgen: https://github.com/golang/mock

To generate mocks using mockgen (storage example):

```
mockgen -source pkg/storage.go -destination mocks/storage.go -package mocks
```

Before pushing run some checks:

```
go build cmd/main.go
go vet cmd/
golangci-lint run
```

To run the test suite:

```
go test pkg/tests/* -v
```
