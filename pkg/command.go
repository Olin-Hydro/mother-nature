package pkg

import (
	"fmt"
	"time"
)

type CommandType int

const (
	Sensor CommandType = iota
	ScheduledActuator
	ReactiveActuator
)

type Command struct {
	CmdType CommandType `json:"cmdType"`
	Id      string      `json:"id"`
	Cmd     int         `json:"cmd"`
}

var (
	ActuationTimes RActuationTimes
)

func init() {
	ActuationTimes = RActuationTimes{}
	// TODO: fill in the actuation times by API call
	// For each id in the raconfig
	// Make a CreateRALogsReq request
	// Take the last time in there and store it here
}

type RActuationTimes struct {
	Times map[string]time.Time
}

func newCommand(cmdType CommandType, id string, cmd int) Command {
	command := Command{
		CmdType: cmdType,
		Id:      id,
		Cmd:     cmd,
	}
	return command
}

func checkRACommands(conf GardenConfig) ([]Command, error) {
	raConfigs := conf.ReactiveActuators
	var cmds []Command
	for i := 0; i < len(raConfigs); i++ {
		raConfig := raConfigs[i]
		fmt.Print(raConfig)
	}
	return cmds, nil
}
