package main

import (
	"encoding/json"
	"fmt"
)

type CommandType int

const (
	Sensor CommandType = iota
	Switch
)

type Command struct {
	CmdType CommandType `json:"cmdType"`
	Id      int         `json:"id"`
	Cmd     int         `json:"cmd"`
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
	senseCmd := Command{
		CmdType: Sensor,
		Id:      0,
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
