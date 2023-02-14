package tests

import (
	"bytes"
	"encoding/json"

	"io"
	"net/http"
	"testing"

	mn "github.com/Olin-Hydro/mother-nature/cmd"
	"github.com/Olin-Hydro/mother-nature/mocks"
	"github.com/Olin-Hydro/mother-nature/pkg"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	sensorId       = "abc432"
	raConfId       = "abc234"
	raId           = "abc112"
	saId           = "abc654"
	gardenConfigId = "abc321"
	gardenId       = "abc123"
	sensorValue    = 1.1
)

func mockRaCache() pkg.RACache {
	return pkg.RACache{
		SensorLogs: map[string]pkg.SensorLog{
			raConfId: {
				Id:        "abc",
				Name:      "sensor_name",
				SensorId:  sensorId,
				Value:     sensorValue,
				CreatedAt: "1970-01-01T00:00:00.000Z",
			},
		},
		RAs: map[string]pkg.RA{
			raConfId: {
				Id:        raId,
				Name:      "ra_name",
				SensorId:  sensorId,
				CreatedAt: "1970-01-01T00:00:00.000Z",
			},
		},
	}
}

func mockRA() pkg.RA {
	return pkg.RA{
		Id:        raId,
		Name:      "ra_name",
		SensorId:  sensorId,
		CreatedAt: "1970-01-01T00:00:00.000Z",
	}
}

func mockRALogs() pkg.RALogs {
	return pkg.RALogs{
		Logs: []pkg.RALog{{
			Id:         "abc",
			Name:       "ra_log_name",
			ActuatorId: raId,
			Data:       "1",
			CreatedAt:  "1970-01-01T00:00:00.000Z",
		}},
	}
}

func mockSensorLogs() pkg.SensorLogs {
	return pkg.SensorLogs{
		Logs: []pkg.SensorLog{
			{
				Id:        "abc",
				Name:      "sensor_name",
				SensorId:  sensorId,
				Value:     sensorValue,
				CreatedAt: "1970-01-01T00:00:00.000Z",
			},
		},
	}
}

func mockGardenConfig() pkg.GardenConfig {
	sensor := pkg.SensorConfig{
		Id:       sensorId,
		Interval: 300.0,
	}
	sensors := []pkg.SensorConfig{sensor}
	raConfig := pkg.RAConfig{
		RAId:          raConfId,
		Interval:      1200.0,
		Threshold:     8.0,
		Duration:      5.0,
		ThresholdType: 1,
	}
	ras := []pkg.RAConfig{raConfig}
	onTimes := []string{"1970-01-01T00:00:00.000Z"}
	offTimes := []string{"1970-01-01T00:01:00.000Z"}
	saConfig := pkg.SAConfig{
		SAId: saId,
		On:   onTimes,
		Off:  offTimes,
	}
	sas := []pkg.SAConfig{saConfig}
	gardenConfig := pkg.GardenConfig{
		Id:        gardenConfigId,
		Name:      "Config_1",
		Sensors:   sensors,
		SAConfigs: sas,
		RAConfigs: ras,
		CreatedAt: "1970-01-01T00:00:00.000Z",
	}
	return gardenConfig
}

func mockGarden() pkg.Garden {
	garden := pkg.Garden{
		Id:        gardenId,
		Name:      "Garden_1",
		Location:  "Mac_3_EndCap",
		Config:    mockGardenConfig(),
		CreatedAt: "1970-01-01T00:00:00.000Z",
	}
	return garden
}

func mockScheduledCmds() []pkg.ScheduledCommand {
	var schedCmds []pkg.ScheduledCommand
	cmdOn := pkg.Command{
		CmdType: pkg.ScheduledActuator,
		Id:      saId,
		Cmd:     1,
	}
	cmdOff := pkg.Command{
		CmdType: pkg.ScheduledActuator,
		Id:      saId,
		Cmd:     0,
	}
	saCmdOn := pkg.ScheduledCommand{
		Cmd:  cmdOn,
		Time: "00:00:00",
	}
	saCmdOff := pkg.ScheduledCommand{
		Cmd:  cmdOff,
		Time: "00:01:00",
	}
	schedCmds = append(schedCmds, saCmdOff, saCmdOn)
	return schedCmds
}

func mockSchedule() pkg.Schedule {
	return pkg.Schedule{
		Cmds: mockScheduledCmds(),
	}
}

func TestGetConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorage := mocks.NewMockStorage(ctrl)
	mockStorage.EXPECT().CreateGardenReq("abc").Return(&http.Request{}, nil)
	mockClient := mocks.NewMockHTTPClient(ctrl)
	b, err := json.Marshal(mockGarden())
	assert.NoError(t, err)
	r := io.NopCloser(bytes.NewReader(b))
	res := http.Response{
		StatusCode: 200,
		Body:       r,
	}
	mockClient.EXPECT().Do(&http.Request{}).Return(&res, nil)
	gardenConfig, err := mn.GetGardenConfig(mockStorage, mockClient, "abc")
	assert.NoError(t, err)
	assert.Equal(t, mockGarden().Config, gardenConfig)
}
