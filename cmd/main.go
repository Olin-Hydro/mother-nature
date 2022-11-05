package main

import (
	"encoding/json"
	"fmt"

	"github.com/Olin-Hydro/mother-nature/pkg"
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

func getGardenConfig(store pkg.Storage, gardenId string) (pkg.GardenConfig, error) {
	garden, err := store.GetGarden(gardenId)
	if err != nil {
		return garden.Config, fmt.Errorf("#getConfig: %e", err)
	}
	return garden.Config, nil
}

func main() {
	conf := pkg.LoadConfigFromEnv()
	h, err := pkg.NewHydrangea(conf.HydrangeaGardenURL, conf.HydrangeaRALogURL, conf.HydrangeaSensorLogURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	gardenConfig, err := getGardenConfig(h, conf.GardenId)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(gardenConfig)
	// TODO:
	// Create schedule from config
	// Send config to gardener
	// TODO:
	// Check conditions
	// Send commands to gardener if needed

	// Placeholder command code
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
