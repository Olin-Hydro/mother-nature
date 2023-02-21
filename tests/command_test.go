package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/Olin-Hydro/mother-nature/mocks"
	"github.com/Olin-Hydro/mother-nature/pkg"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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

func TestCreateRACommands(t *testing.T) {
	pkg.Cache = mockRaCache()
	cmds, err := pkg.CreateRACommands(mockGardenConfig())
	expectedCmd := pkg.Command{
		CmdType:  pkg.ReactiveActuator,
		Id:       raId,
		Cmd:      1,
		GardenId: gardenId,
	}
	assert.NoError(t, err)
	assert.Equal(t, len(cmds), 1)
	assert.Equal(t, cmds[0], expectedCmd)
}

func TestUpdateRACache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorage := mocks.NewMockStorage(ctrl)
	mockClient := mocks.NewMockHTTPClient(ctrl)
	raConfigs := mockGardenConfig().RAConfigs

	b, err := json.Marshal(mockRA())
	assert.NoError(t, err)
	r := io.NopCloser(bytes.NewReader(b))
	res := http.Response{
		StatusCode: 200,
		Body:       r,
	}
	mockStorage.EXPECT().CreateRAReq(raConfigs[0].RAId).Return(&http.Request{}, nil)
	mockClient.EXPECT().Do(&http.Request{}).Return(&res, nil)
	b2, err := json.Marshal(mockSensorLogs())
	assert.NoError(t, err)
	r2 := io.NopCloser(bytes.NewReader(b2))
	res2 := http.Response{
		StatusCode: 200,
		Body:       r2,
	}
	mockStorage.EXPECT().CreateSensorLogsReq(sensorId, "1").Return(&http.Request{}, nil)
	mockClient.EXPECT().Do(&http.Request{}).Return(&res2, nil)
	mockStorage.EXPECT().CreateRALogsReq(raId, "1").Return(&http.Request{}, nil)
	b3, err := json.Marshal(mockRaLogs())
	assert.NoError(t, err)
	r3 := io.NopCloser(bytes.NewReader(b3))
	res3 := http.Response{
		StatusCode: 200,
		Body:       r3,
	}
	mockClient.EXPECT().Do(&http.Request{}).Return(&res3, nil)
	cache, err := pkg.UpdateRACache(pkg.Cache, raConfigs, mockRaCache().GardenId, mockClient, mockStorage)
	assert.NoError(t, err)
	assert.Equal(t, mockRaCache(), cache)
}
