package pkg

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

func TestCreateRACommands(t *testing.T) {
	pkg.Cache = mocks.MockRaCache()
	cmds, err := pkg.CreateRACommands(mocks.MockGardenConfig())
	expectedCmd := pkg.Command{
		CmdType:  pkg.ReactiveActuator,
		Id:       mocks.RAId,
		Cmd:      1,
		GardenId: mocks.GardenId,
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
	raConfigs := mocks.MockGardenConfig().RAConfigs

	b, err := json.Marshal(mocks.MockRA())
	assert.NoError(t, err)
	r := io.NopCloser(bytes.NewReader(b))
	res := http.Response{
		StatusCode: 200,
		Body:       r,
	}
	mockStorage.EXPECT().CreateRAReq(raConfigs[0].RAId).Return(&http.Request{}, nil)
	mockClient.EXPECT().Do(&http.Request{}).Return(&res, nil)
	b2, err := json.Marshal(mocks.MockSensorLogs())
	assert.NoError(t, err)
	r2 := io.NopCloser(bytes.NewReader(b2))
	res2 := http.Response{
		StatusCode: 200,
		Body:       r2,
	}
	mockStorage.EXPECT().CreateSensorLogsReq(mocks.SensorId, "1").Return(&http.Request{}, nil)
	mockClient.EXPECT().Do(&http.Request{}).Return(&res2, nil)
	mockStorage.EXPECT().CreateRALogsReq(mocks.RAId, "1").Return(&http.Request{}, nil)
	b3, err := json.Marshal(mocks.MockRaLogs())
	assert.NoError(t, err)
	r3 := io.NopCloser(bytes.NewReader(b3))
	res3 := http.Response{
		StatusCode: 200,
		Body:       r3,
	}
	mockClient.EXPECT().Do(&http.Request{}).Return(&res3, nil)
	cache, err := pkg.UpdateRACache(pkg.Cache, raConfigs, mocks.MockRaCache().GardenId, mockClient, mockStorage)
	assert.NoError(t, err)
	assert.Equal(t, mocks.MockRaCache(), cache)
}
