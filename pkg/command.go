package pkg

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	Sensor            string = "Sensor"
	ScheduledActuator string = "Scheduled Actuator"
	ReactiveActuator  string = "Reactive Actuator"
)

type Command struct {
	RefId    string `json:"ref_id"`
	CmdType  string `json:"type"`
	Cmd      int    `json:"cmd"`
	GardenId string `json:"garden_id"`
}

type RACache struct {
	ActuationTimes map[string]time.Time
	SensorLogs     map[string]SensorLog
	RAs            map[string]RA
	GardenId       string
}

var (
	Cache RACache
)

func init() {
	Cache = RACache{
		ActuationTimes: make(map[string]time.Time),
		SensorLogs:     make(map[string]SensorLog),
		RAs:            make(map[string]RA),
	}
}

func newCommand(cmdType string, id string, cmd int, gardenId string) Command {
	command := Command{
		RefId:    id,
		CmdType:  cmdType,
		Cmd:      cmd,
		GardenId: gardenId,
	}
	return command
}

func StrToTime(dtStr string) (time.Time, error) {
	dt, err := time.Parse(time.RFC3339Nano, dtStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("#StrToTime: %e", err)
	}
	return dt, nil
}

func CreateRACommands(conf GardenConfig) ([]Command, error) {
	raConfigs := conf.RAConfigs
	var cmds []Command
	var errors []error
	for i := 0; i < len(raConfigs); i++ {
		log.Info(fmt.Sprintf("Checking RA %s with ID %s", Cache.RAs[raConfigs[i].RAId].Name, raConfigs[i].RAId))
		raConfig := raConfigs[i]
		needed, err := raCommandNeeded(raConfig)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		if !needed {
			continue
		}
		log.Info(fmt.Sprintf("Creating command for %s with type %s and value %d", Cache.RAs[raConfig.RAId].Name, ReactiveActuator, 1))
		cmd := newCommand(ReactiveActuator, Cache.RAs[raConfig.RAId].Id, 1, Cache.GardenId)
		cmds = append(cmds, cmd)
	}
	if len(errors) != 0 {
		return cmds, fmt.Errorf("#CreateRACommands: %e", errors)
	}
	return cmds, nil
}

func isPassedInterval(raConfig RAConfig) (bool, error) {
	actTime, ok := Cache.ActuationTimes[raConfig.RAId]
	if !ok {
		return false, fmt.Errorf("#isPassedInterval: RAId %s not found in RACache", raConfig.RAId)
	}
	timeDiff := time.Now().UTC().Sub(actTime).Seconds()
	log.Info(fmt.Sprintf("Comparing time since actuation of %f seconds and %f seconds interval", timeDiff, raConfig.Interval))
	return timeDiff <= raConfig.Interval, nil
}

func raCommandNeeded(raConfig RAConfig) (bool, error) {
	passed, err := isPassedInterval(raConfig)
	if err != nil {
		return false, fmt.Errorf("#raCommandNeeded: %d", raConfig.ThresholdType)
	} else if passed {
		return false, nil
	}
	switch raConfig.ThresholdType {
	case 0:
		if raConfig.Threshold < Cache.SensorLogs[raConfig.RAId].Value {
			log.Info(fmt.Sprintf("Threshold: %f is smaller than sensor value %f", raConfig.Threshold, Cache.SensorLogs[raConfig.RAId].Value))
			return true, nil
		}
	case 1:
		if raConfig.Threshold > Cache.SensorLogs[raConfig.RAId].Value {
			log.Info(fmt.Sprintf("Threshold: %f is greater than sensor value %f", raConfig.Threshold, Cache.SensorLogs[raConfig.RAId].Value))
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
	if res.StatusCode != http.StatusOK {
		return sensorLog, fmt.Errorf("#getSensorLog: %s", res.Status)
	}
	err = DecodeJson(&sensorLogs.Logs, res.Body)
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
	// No logs found, default time to be a long time ago
	if res.StatusCode == http.StatusNotFound {
		raLog.Id = raId
		raLog.CreatedAt = "2000-03-02T16:13:07.437575+00:00"
		log.Warn(fmt.Sprintf("No logs found for ra with id %s, allowing command", raId))
		return raLog, nil
	}
	err = DecodeJson(&raLogs.Logs, res.Body)
	if err != nil {
		return raLog, fmt.Errorf("getRaLogs: %e", err)
	}
	if len(raLogs.Logs) != 1 {
		return raLog, fmt.Errorf("#getRaLog: storage returned wrong number of logs: expected 1 got %d", len(raLogs.Logs))
	}
	return raLogs.Logs[0], nil
}

func UpdateRACache(Cache RACache, raConfigs []RAConfig, gardenId string, client HTTPClient, store Storage) (RACache, error) {
	Cache.GardenId = gardenId
	for i := 0; i < len(raConfigs); i++ {
		raId := raConfigs[i].RAId
		ra, err := getRa(raId, store, client)
		if err != nil {
			return Cache, fmt.Errorf("#UpdateRACache: %e", err)
		}
		Cache.RAs[raId] = ra
		sensorLog, err := getSensorLog(ra.SensorId, store, client)
		if err != nil {
			return Cache, fmt.Errorf("#UpdateRACache: %e", err)
		}
		Cache.SensorLogs[raId] = sensorLog
		raLog, err := getRaLogs(Cache.RAs[raId].Id, store, client)
		if err != nil {
			return Cache, fmt.Errorf("#UpdateRACache: %e", err)
		}
		actTime, err := StrToTime(raLog.CreatedAt)
		if err != nil {
			return Cache, fmt.Errorf("UpdateRACache: %e", err)
		}
		Cache.ActuationTimes[raId] = actTime
	}
	return Cache, nil
}
