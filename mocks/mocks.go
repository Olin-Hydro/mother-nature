package mocks

import (
	"time"

	pkg "github.com/Olin-Hydro/mother-nature/pkg"
)

const (
	SensorId       = "abc432"
	raConfId       = "abc234"
	RAId           = "abc112"
	saId           = "abc654"
	GardenConfigId = "abc321"
	GardenId       = "abc123"
	sensorValue    = 1.1
)

func MockCommands() []pkg.Command {
	cmds := []pkg.Command{
		{
			CmdType:  pkg.ReactiveActuator,
			Id:       "abc",
			Cmd:      1,
			GardenId: "bcd",
		},
		{
			CmdType:  pkg.ScheduledActuator,
			Id:       "aab",
			Cmd:      0,
			GardenId: "bcd",
		},
	}
	return cmds
}

func MockRaCache() pkg.RACache {
	t, _ := pkg.StrToTime("1970-01-01T00:00:00.000Z")
	return pkg.RACache{
		SensorLogs: map[string]pkg.SensorLog{
			RAId: {
				Id:        "abc",
				Name:      "sensor_name",
				SensorId:  SensorId,
				Value:     sensorValue,
				CreatedAt: "1970-01-01T00:00:00.000Z",
			},
		},
		RAs: map[string]pkg.RA{
			RAId: {
				Id:        RAId,
				Name:      "ra_name",
				SensorId:  SensorId,
				CreatedAt: "1970-01-01T00:00:00.000Z",
			},
		},
		GardenId: GardenId,
		ActuationTimes: map[string]time.Time{
			RAId: t,
		},
	}
}

func MockRA() pkg.RA {
	return pkg.RA{
		Id:        RAId,
		Name:      "ra_name",
		SensorId:  SensorId,
		CreatedAt: "1970-01-01T00:00:00.000Z",
	}
}

func MockSensorLogs() pkg.SensorLogs {
	return pkg.SensorLogs{
		Logs: []pkg.SensorLog{
			{
				Id:        "abc",
				Name:      "sensor_name",
				SensorId:  SensorId,
				Value:     sensorValue,
				CreatedAt: "1970-01-01T00:00:00.000Z",
			},
		},
	}
}

func MockRaLogs() pkg.RALogs {
	return pkg.RALogs{
		Logs: []pkg.RALog{{
			Id:         "abc",
			Name:       "ra_name",
			ActuatorId: RAId,
			Data:       "on",
			CreatedAt:  "1970-01-01T00:00:00.000Z",
		}},
	}
}

func MockGardenConfig() pkg.GardenConfig {
	sensor := pkg.SensorConfig{
		Id:       SensorId,
		Interval: 300.0,
	}
	sensors := []pkg.SensorConfig{sensor}
	raConfig := pkg.RAConfig{
		RAId:          RAId,
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
		Id:        GardenConfigId,
		Name:      "Config_1",
		Sensors:   sensors,
		SAConfigs: sas,
		RAConfigs: ras,
		CreatedAt: "1970-01-01T00:00:00.000Z",
	}
	return gardenConfig
}

func MockGarden() pkg.Garden {
	garden := pkg.Garden{
		Id:        GardenId,
		Name:      "Garden_1",
		Location:  "Mac_3_EndCap",
		ConfigID:  GardenConfigId,
		CreatedAt: "1970-01-01T00:00:00.000Z",
	}
	return garden
}
