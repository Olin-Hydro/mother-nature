# mother-nature

Install some dependencies:

- golangci-lint: https://golangci-lint.run/usage/install/
- mockgen: https://github.com/golang/mock

Create a .env file with the following variables:

```
HYDRANGEA_GARDEN_URL=https://hydrangea.kaz.farm/garden
HYDRANGEA_RALOG_URL=https://hydrangea.kaz.farm/ra/logging/actions
HYDRANGEA_SENSORLOG_URL=https://hydrangea.kaz.farm/sensors/logging
HYDRANGEA_RA_URL=https://hydrangea.kaz.farm/ra
HYDRANGEA_COMMAND_URL=https://hydrangea.kaz.farm/cmd
HYDRANGEA_CONFIG_URL=https://hydrangea.kaz.farm/config
GARDEN_ID=
API_KEY=
```

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
