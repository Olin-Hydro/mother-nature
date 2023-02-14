package pkg

import (
	"fmt"
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

type RACache struct {
	SensorLogs map[string]SensorLog
	RAs        map[string]RA
}

var (
	Cache RACache
)

func init() {
	Cache = RACache{}
}

func newCommand(cmdType CommandType, id string, cmd int) Command {
	command := Command{
		CmdType: cmdType,
		Id:      id,
		Cmd:     cmd,
	}
	return command
}

func CreateRACommands(conf GardenConfig) ([]Command, error) {
	raConfigs := conf.RAConfigs
	var cmds []Command
	var errors []error
	for i := 0; i < len(raConfigs); i++ {
		raConfig := raConfigs[i]
		needed, err := raCommandNeeded(raConfig)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		if !needed {
			continue
		}
		cmd := newCommand(ReactiveActuator, Cache.RAs[raConfig.RAId].Id, 1)
		cmds = append(cmds, cmd)
	}
	if len(errors) != 0 {
		return cmds, fmt.Errorf("#CreateRACommands: %e", errors)
	}
	return cmds, nil
}

func raCommandNeeded(raConfig RAConfig) (bool, error) {
	switch raConfig.ThresholdType {
	case 0:
		if raConfig.Threshold < Cache.SensorLogs[raConfig.RAId].Value {
			return true, nil
		}
	case 1:
		if raConfig.Threshold > Cache.SensorLogs[raConfig.RAId].Value {
			return true, nil
		}
	default:
		return false, fmt.Errorf("#raCommandNeeded: Invalid thresholdtype: %d", raConfig.ThresholdType)
	}
	return false, nil
}

func getRa(raConfigId string, store Storage, client HTTPClient) (RA, error) {
	ra := RA{}
	req, err := store.CreateRAReq(raConfigId)
	if err != nil {
		return ra, fmt.Errorf("#getRa: %e", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return ra, fmt.Errorf("#getRa: %e", err)
	}
	err = DecodeJson(&ra, res.Body)
	if err != nil {
		return ra, fmt.Errorf("#getRa: %e", err)
	}
	return ra, nil
}

func getSensorLog(sensorId string, store Storage, client HTTPClient) (SensorLog, error) {
	sensorLog := SensorLog{}
	sensorLogs := SensorLogs{}
	req, err := store.CreateSensorLogsReq(sensorId, "1")
	if err != nil {
		return sensorLog, fmt.Errorf("#getSensorLog: %e", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return sensorLog, fmt.Errorf("#getSensorLog: %e", err)
	}
	err = DecodeJson(&sensorLogs, res.Body)
	if err != nil {
		return sensorLog, fmt.Errorf("#getSensorLog: %e", err)
	}
	if len(sensorLogs.Logs) != 1 {
		return sensorLog, fmt.Errorf("#getSensorLog: storage returned wrong number of logs: expected 1 got %d", len(sensorLogs.Logs))
	}
	return sensorLogs.Logs[0], nil
}

func getRaLogs(raId string, store Storage, client HTTPClient) (RALog, error) {
	raLog := RALog{}
	raLogs := RALogs{}
	req, err := store.CreateRALogsReq(raId, "1")
	if err != nil {
		return raLog, fmt.Errorf("#getRaLogs: %e", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return raLog, fmt.Errorf("#getRaLogs: %e", err)
	}
	err = DecodeJson(&raLogs, res.Body)
	if err != nil {
		return raLog, fmt.Errorf("getRaLogs: %e", err)
	}
	if len(raLogs.Logs) != 1 {
		return raLog, fmt.Errorf("#getSensorLog: storage returned wrong number of logs: expected 1 got %d", len(raLogs.Logs))
	}
	return raLogs.Logs[0], nil
}

func UpdateRACache(Cache RACache, raConfigs []RAConfig, client HTTPClient, store Storage) (RACache, error) {
	for i := 0; i < len(raConfigs); i++ {
		raConfigId := raConfigs[i].RAId
		ra, err := getRa(raConfigId, store, client)
		if err != nil {
			return Cache, fmt.Errorf("#UpdateRACache: %e", err)
		}
		Cache.RAs[raConfigId] = ra
		sensorLog, err := getSensorLog(ra.SensorId, store, client)
		if err != nil {
			return Cache, fmt.Errorf("#UpdateRACache: %e", err)
		}
		Cache.SensorLogs[raConfigId] = sensorLog
	}
	return Cache, nil
}
