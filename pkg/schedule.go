package pkg

import (
	"encoding/json"
	"fmt"
	"time"
)

type Schedule struct {
	Cmds []ScheduledCommand `json:"commands"`
}

type ScheduledCommand struct {
	Cmd  Command `json:"cmd"`
	Time string  `json:"time"`
}

func StrToTime(dtStr string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.000Z"
	dt, err := time.Parse(layout, dtStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("#StrToTime: %e", err)
	}
	return dt, nil
}

func newScheduledCommand(cmd Command, dt time.Time) ScheduledCommand {
	schedCmd := ScheduledCommand{
		Cmd:  cmd,
		Time: dt.Format("15:04:05"),
	}
	return schedCmd
}

func createScheduledCmds(SAconf SAConfig) ([]ScheduledCommand, error) {
	var schedCmds []ScheduledCommand
	for t := 0; t < len(SAconf.Off); t++ {
		dt, err := StrToTime(SAconf.Off[t])
		if err != nil {
			return schedCmds, fmt.Errorf("#NewSchedule: %e", err)
		}
		cmd := newCommand(ScheduledActuator, SAconf.SAId, 0)
		schedCmd := newScheduledCommand(cmd, dt)
		schedCmds = append(schedCmds, schedCmd)
	}
	for t := 0; t < len(SAconf.On); t++ {
		dt, err := StrToTime(SAconf.On[t])
		if err != nil {
			return schedCmds, fmt.Errorf("#NewSchedule: %e", err)
		}
		cmd := newCommand(ScheduledActuator, SAconf.SAId, 1)
		schedCmd := newScheduledCommand(cmd, dt)
		schedCmds = append(schedCmds, schedCmd)
	}
	return schedCmds, nil
}

func NewSchedule(gardenConfig GardenConfig) (Schedule, error) {
	sched := Schedule{}
	SAS := gardenConfig.SAConfigs
	for i := 0; i < len(SAS); i++ {
		SAconf := SAS[i]
		schedCmds, err := createScheduledCmds(SAconf)
		if err != nil {
			return sched, fmt.Errorf("#NewSchedule: %e", err)
		}
		sched.Cmds = append(sched.Cmds, schedCmds...)
	}
	return sched, nil
}

func EncodeSchedule(s Schedule) (b []byte, e error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("#EncodeSchedule: %e", err)
	}
	return b, nil
}
