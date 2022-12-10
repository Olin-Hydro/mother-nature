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

type RACache struct {
	ActuationTimes map[string]time.Time
	SensorLogs     map[string]SensorLog
	RAs            map[string]RA
}

var (
	Cache RACache
)

func init() {
	Cache = RACache{}
	// TODO: fill in the actuation times by API call
	// For each id in the raconfig
	// Make a CreateRALogsReq request
	// Take the last time in there and store it here
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
		switch raConfig.ThresholdType {
		case 0:

		case 1:

		}
		fmt.Print(raConfig)
	}
	return cmds, nil
}

func isPassedInterval(raConfig RAConfig) bool {
	timeDiff := time.Now().UTC().Sub(Cache.ActuationTimes[raConfig.Id]).Seconds()
	if timeDiff < raConfig.Interval {
		return false
	}
	return true
}

func raCommandNeeded(raConfig RAConfig) (bool, error) {
	if !isPassedInterval(raConfig) {
		return false, nil
	}
	// switch raConfig.ThresholdType {
	// case 0:
	// 	if raConfig.Threshold <
	// case 1:

	// }
	return true, nil
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

func getRaLogs(raConfigId string, store Storage, client HTTPClient) (RALog, error) {
	raLog := RALog{}
	raLogs := RALogs{}
	req, err := store.CreateRALogsReq(raConfigId, "1")
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
		raConfigId := raConfigs[i].Id
		if _, ok := Cache.RAs[raConfigId]; !ok {
			ra, err := getRa(raConfigId, store, client)
			if err != nil {
				return Cache, fmt.Errorf("#UpdateRACache: %e", err)
			}
			Cache.RAs[raConfigId] = ra
		}
		sensorLog, err := getSensorLog(Cache.RAs[raConfigId].Id, store, client)
		if err != nil {
			return Cache, fmt.Errorf("#UpdateRACache: %e", err)
		}
		Cache.SensorLogs[raConfigId] = sensorLog
		raLog, err := getRaLogs(raConfigId, store, client)
		if err != nil {
			return Cache, fmt.Errorf("#UpdateRACache: %e", err)
		}
		actTime, err := StrToTime(raLog.CreatedAt)
		if err != nil {
			return Cache, fmt.Errorf("UpdateRACache: %e", err)
		}
		Cache.ActuationTimes[raConfigId] = actTime
	}
	return Cache, nil
}
