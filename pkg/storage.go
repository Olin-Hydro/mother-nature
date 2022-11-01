package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Garden struct {
	Id        string       `json:"id"`
	Name      string       `json:"name"`
	Location  string       `json:"location"`
	Config    GardenConfig `json:"config"`
	CreatedAt string       `json:"created_at"`
}

type GardenConfig struct {
	Id                 string         `json:"id"`
	Name               string         `json:"name"`
	Sensors            []SensorConfig `json:"sensors"`
	ScheduledActuators []SAConfig     `json:"scheduled_actuators"`
	ReactiveActuators  []RAConfig     `json:"reactive_actuators"`
	CreatedAt          string         `json:"created_at"`
}

type SensorConfig struct {
	Id       string  `json:"id"`
	Interval float64 `json:"interval"`
}

type SAConfig struct {
	Id  string   `json:"id"`
	On  []string `json:"on"`
	Off []string `json:"off"`
}

type RAConfig struct {
	Id            string  `json:"id"`
	Interval      float64 `json:"interval"`
	Threshold     float64 `json:"threshold"`
	Duration      float64 `json:"duration"`
	ThresholdType int     `json:"threshold_type"`
}

type RALogs struct {
	Logs []RALogs `json:"logs"`
}

type RALog struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	ActuatorId string `json:"actuator_id"`
	Data       string `json:"data"`
	CreatedAt  string `json:"created_at"`
}

type SensorLogs struct {
	Logs []SensorLog `json:"logs"`
}

type SensorLog struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	SensorId  string  `json:"sensor_id"`
	Value     float64 `json:"value"`
	CreatedAt string  `json:"created_at"`
}

type Storage interface {
	GetGarden(gardenId string) (Garden, error)
	GetRALogs(RAId string, limit string) (RALogs, error)
	GetSensorLogs(SensorId string, limit string) (SensorLogs, error)
}

type Hydrangea struct {
	GardenURL    url.URL
	RALogURL     url.URL
	SensorLogURL url.URL
}

func NewHydrangea(gardenURL string, raLogURL string, sensorLogURL string) (Hydrangea, error) {
	h := Hydrangea{}
	u, err := url.Parse(gardenURL)
	if err != nil {
		return h, fmt.Errorf("#NewHydrangea: %e", err)
	}
	h.GardenURL = *u
	u, err = url.Parse(raLogURL)
	if err != nil {
		return h, fmt.Errorf("#NewHydrangea: %e", err)
	}
	h.RALogURL = *u
	u, err = url.Parse(sensorLogURL)
	if err != nil {
		return h, fmt.Errorf("#NewHydrangea: %e", err)
	}
	h.SensorLogURL = *u
	return h, nil
}

func (h Hydrangea) GetGarden(gardenId string) (Garden, error) {
	garden := Garden{}
	values := h.GardenURL.Query()
	values.Add("id", gardenId)
	h.GardenURL.RawQuery = values.Encode()
	res, err := http.Get(h.GardenURL.String())
	if err != nil {
		return garden, fmt.Errorf("#Hydrangea.GetConfig: %e", err)
	}
	err = json.NewDecoder(res.Body).Decode(&garden)
	if err != nil {
		return garden, fmt.Errorf("#Hydrangea.GetConfig: %e", err)
	}
	return garden, nil
}

func (h Hydrangea) GetRALogs(RAId string, limit string) (RALogs, error) {
	logs := RALogs{}
	values := h.RALogURL.Query()
	values.Add("id", RAId)
	values.Add("limit", limit)
	h.RALogURL.RawQuery = values.Encode()
	res, err := http.Get(h.RALogURL.String())
	if err != nil {
		return logs, fmt.Errorf("#Hydrangea.GetRALogs: %e", err)
	}
	err = json.NewDecoder(res.Body).Decode(&logs)
	if err != nil {
		return logs, fmt.Errorf("#Hydrangea.GetRALogs: %e", err)
	}
	return logs, nil
}

func (h Hydrangea) GetSensorLogs(SensorId string, limit string) (SensorLogs, error) {
	logs := SensorLogs{}
	values := h.SensorLogURL.Query()
	values.Add("id", SensorId)
	values.Add("limit", limit)
	h.SensorLogURL.RawQuery = values.Encode()
	res, err := http.Get(h.SensorLogURL.String())
	if err != nil {
		return logs, fmt.Errorf("#Hydrangea.GetSensorLogs: %e", err)
	}
	err = json.NewDecoder(res.Body).Decode(&logs)
	if err != nil {
		return logs, fmt.Errorf("#Hydrangea.GetSensorLogs: %e", err)
	}
	return logs, nil
}
