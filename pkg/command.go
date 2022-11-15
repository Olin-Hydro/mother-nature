package pkg

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

func newCommand(cmdType CommandType, id string, cmd int) Command {
	command := Command{
		CmdType: cmdType,
		Id:      id,
		Cmd:     cmd,
	}
	return command
}
