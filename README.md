# mother-nature

## Overview

Mother nature is a controller that monitors information from
[hydrangea](https://hydrangea.kaz.farm/), and creates new commands as needed
that are store in hydrangea. Mother nature focuses solely on reactive actuators.
These are actuators such as ph down pumps or nutrients pumps. These actuators
turn on and off relative to a related sensor. For example, if the ph goes above
a certain threshold we want a command to be created that will turn on a ph down
pump.

Mother nature follows this logic:

- Get a garden config based on a garden_id
- Get the relevant reactive actuator objects for that garden
- For each reactive actuator:
  - Based on reactive acuator config, check that the reactive actuator has not
    been sent a command recently
  - Get the latest relevant sensor logs based on sensor_id and store it in
    RA_cache
  - If a sensor log is out of the defined range, create a reactive actuator
    command
- Send the commands to hydrangea

## Deployment

Mother nature is designed to run periodically (every ~10min). It is currently
deployed as a lambda function in AWS.

## Setup

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
