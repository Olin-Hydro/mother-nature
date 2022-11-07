package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	mn "github.com/Olin-Hydro/mother-nature/cmd"
	"github.com/Olin-Hydro/mother-nature/mocks"
	"github.com/Olin-Hydro/mother-nature/pkg"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func mockGarden() pkg.Garden {
	sensor := pkg.SensorConfig{
		Id:       "abc432",
		Interval: 300.0,
	}
	sensors := []pkg.SensorConfig{sensor}
	raConfig := pkg.RAConfig{
		Id:            "abc234",
		Interval:      1200.0,
		Threshold:     8.0,
		Duration:      5.0,
		ThresholdType: 1,
	}
	ras := []pkg.RAConfig{raConfig}
	onTimes := []string{"1970-01-01T00:00:00.000Z"}
	offTimes := []string{"1970-01-01T00:01:00.000Z"}
	saConfig := pkg.SAConfig{
		Id:  "abc654",
		On:  onTimes,
		Off: offTimes,
	}
	sas := []pkg.SAConfig{saConfig}
	gardenConfig := pkg.GardenConfig{
		Id:                 "abc321",
		Name:               "Config_1",
		Sensors:            sensors,
		ScheduledActuators: sas,
		ReactiveActuators:  ras,
		CreatedAt:          "1970-01-01T00:00:00.000Z",
	}
	garden := pkg.Garden{
		Id:        "abc123",
		Name:      "Garden_1",
		Location:  "Mac_3_EndCap",
		Config:    gardenConfig,
		CreatedAt: "1970-01-01T00:00:00.000Z",
	}
	return garden
}

func TestGetConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorage := mocks.NewMockStorage(ctrl)
	mockStorage.EXPECT().CreateGardenReq("abc").Return(&http.Request{}, nil)
	mockClient := mocks.NewMockHTTPClient(ctrl)
	b, err := json.Marshal(mockGarden())
	assert.NoError(t, err)
	//nolint:staticcheck
	r := ioutil.NopCloser(bytes.NewReader(b))
	res := http.Response{
		StatusCode: 200,
		Body:       r,
	}
	mockClient.EXPECT().Do(&http.Request{}).Return(&res, nil)
	gardenConfig, err := mn.GetGardenConfig(mockStorage, mockClient, "abc")
	assert.NoError(t, err)
	assert.Equal(t, mockGarden().Config, gardenConfig)
}
