package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type Garden struct {
	Id        string `json:"_id"`
	Name      string `json:"name"`
	Location  string `json:"location"`
	ConfigID  string `json:"config_id"`
	CreatedAt string `json:"created_at"`
}

type GardenConfig struct {
	Id        string         `json:"_id"`
	Name      string         `json:"name"`
	Sensors   []SensorConfig `json:"sensors"`
	SAConfigs []SAConfig     `json:"sa_schedule"`
	RAConfigs []RAConfig     `json:"ra_schedule"`
	CreatedAt string         `json:"created_at"`
}

type SensorConfig struct {
	Id       string  `json:"_id"`
	Interval float64 `json:"interval"`
}

type SAConfig struct {
	SAId string   `json:"sa_id"`
	On   []string `json:"on"`
	Off  []string `json:"off"`
}

type RAConfig struct {
	RAId          string  `json:"ra_id"`
	Interval      float64 `json:"interval"`
	Threshold     float64 `json:"threshold"`
	Duration      float64 `json:"duration"`
	ThresholdType int     `json:"threshold_type"`
}

type RALogs struct {
	Logs []RALog `json:"logs"`
}

type RA struct {
	Id        string `json:"_id"`
	Name      string `json:"name"`
	SensorId  string `json:"sensor_id"`
	CreatedAt string `json:"created_at"`
}

type RALog struct {
	Id         string `json:"_id"`
	Name       string `json:"name"`
	ActuatorId string `json:"actuator_id"`
	Data       string `json:"data"`
	CreatedAt  string `json:"created_at"`
}

type SensorLogs struct {
	Logs []SensorLog `json:"logs"`
}

type SensorLog struct {
	Id        string  `json:"_id"`
	SensorId  string  `json:"sensor_id"`
	Value     float64 `json:"value"`
	CreatedAt string  `json:"created_at"`
}

type Storage interface {
	CreateGardenReq(gardenId string) (*http.Request, error)
	CreateConfigReq(configId string) (*http.Request, error)
	CreateRALogsReq(RAId string, limit string) (*http.Request, error)
	CreateRAReq(RAConfigId string) (*http.Request, error)
	CreateSensorLogsReq(SensorId string, limit string) (*http.Request, error)
	CreateCommandReq(commands []Command) (*http.Request, error)
}

type Hydrangea struct {
	GardenURL    url.URL
	RALogURL     url.URL
	SensorLogURL url.URL
	RAURL        url.URL
	CommandURL   url.URL
	ConfigURL    url.URL
	ApiKey       string
}

func NewHydrangea(gardenURL string, raLogURL string, raURL string, sensorLogURL string, commmandURL string, configURL string, apiKey string) (Hydrangea, error) {
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
	u, err = url.Parse(raURL)
	if err != nil {
		return h, fmt.Errorf("#NewHydrangea: %e", err)
	}
	h.RAURL = *u
	u, err = url.Parse(sensorLogURL)
	if err != nil {
		return h, fmt.Errorf("#NewHydrangea: %e", err)
	}
	h.SensorLogURL = *u
	u, err = url.Parse(commmandURL)
	if err != nil {
		return h, fmt.Errorf("#NewHydrangea: %e", err)
	}
	h.CommandURL = *u
	u, err = url.Parse(configURL)
	if err != nil {
		return h, fmt.Errorf("#NewHydrangea: %e", err)
	}
	h.ConfigURL = *u
	h.ApiKey = apiKey
	return h, nil
}

func (h Hydrangea) CreateRAReq(RAConfigId string) (*http.Request, error) {
	h.RAURL.Path = path.Join(h.RAURL.Path, RAConfigId)
	req, err := http.NewRequest("GET", h.RAURL.String(), strings.NewReader(""))
	req.Header.Add("x-api-key", h.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("#Hydrangea.CreateGardenReq: %e", err)
	}
	return req, nil
}

func (h Hydrangea) CreateGardenReq(gardenId string) (*http.Request, error) {
	h.GardenURL.Path = path.Join(h.GardenURL.Path, gardenId)
	req, err := http.NewRequest("GET", h.GardenURL.String(), strings.NewReader(""))
	req.Header.Add("x-api-key", h.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("#Hydrangea.CreateGardenReq: %e", err)
	}
	return req, nil
}

func (h Hydrangea) CreateConfigReq(configId string) (*http.Request, error) {
	h.ConfigURL.Path = path.Join(h.ConfigURL.Path, configId)
	req, err := http.NewRequest("GET", h.ConfigURL.String(), strings.NewReader(""))
	req.Header.Add("x-api-key", h.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("#Hydrangea.CreateConfigReq: %e", err)
	}
	return req, nil
}

func (h Hydrangea) CreateRALogsReq(RAId string, limit string) (*http.Request, error) {
	h.RALogURL.Path = path.Join(h.RALogURL.Path, RAId)
	values := h.RALogURL.Query()
	values.Set("limit", limit)
	h.RALogURL.RawQuery = values.Encode()
	req, err := http.NewRequest("GET", h.RALogURL.String(), strings.NewReader(""))
	req.Header.Add("x-api-key", h.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("#Hydrangea.CreateRALogsReq: %e", err)
	}
	return req, nil
}

func (h Hydrangea) CreateSensorLogsReq(SensorId string, limit string) (*http.Request, error) {
	h.SensorLogURL.Path = path.Join(h.SensorLogURL.Path, SensorId)
	values := h.SensorLogURL.Query()
	values.Set("limit", limit)
	h.SensorLogURL.RawQuery = values.Encode()
	req, err := http.NewRequest("GET", h.SensorLogURL.String(), strings.NewReader(""))
	req.Header.Add("x-api-key", h.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("#Hydrangea.CreateSensorLogsReq: %e", err)
	}
	return req, nil
}

func (h Hydrangea) CreateCommandReq(commands []Command) (*http.Request, error) {
	values := h.CommandURL.Query()
	cmdStr, err := json.Marshal(commands)
	if err != nil {
		return nil, fmt.Errorf("#Hydrangea.CreateCommandReq: %e", err)
	}
	h.CommandURL.RawQuery = values.Encode()
	req, err := http.NewRequest("POST", h.CommandURL.String(), bytes.NewBuffer(cmdStr))
	req.Header.Add("x-api-key", h.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("#Hydrangea.CreateCommandReq: %e", err)
	}
	return req, nil
}
