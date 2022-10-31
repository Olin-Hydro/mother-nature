package main

import (
	"encoding/json"
	"fmt"
	mn "olinhydro/mothernature/internal"
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

type Schedule struct {
	Cmds []ScheduledCommand `json:"commands"`
}

type ScheduledCommand struct {
	Cmd      Command `json:"cmd"`
	Datetime int64   `json:"datetime"`
}

func encodeCmd(cmd Command) (b []byte, e error) {
	b, err := json.Marshal(cmd)
	if err != nil {
		return nil, fmt.Errorf("encodeCmd: %e", err)
	}
	return b, nil
}

func decodeCmd(b []byte) (cmd Command, e error) {
	err := json.Unmarshal(b, &cmd)
	if err != nil {
		return cmd, fmt.Errorf("decodeCmd: %e", err)
	}
	return cmd, nil
}

func main() {
	_, err := mn.NewHydrangea(mn.HYDRANGEA_GARDEN_URL, mn.HYDRANGEA_RALOG_URL, mn.HYDRANGEA_SENSORLOG_URL)
	if err != nil {
		fmt.Println(err)
	}
	senseCmd := Command{
		CmdType: Sensor,
		Id:      "sdhb3k3j3kn",
	}
	fmt.Println(senseCmd)
	b, err := encodeCmd(senseCmd)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
	decodedCmd, err := decodeCmd(b)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(decodedCmd)
}
