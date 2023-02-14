package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/Olin-Hydro/mother-nature/mocks"
	"github.com/Olin-Hydro/mother-nature/pkg"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateRACommands(t *testing.T) {
	pkg.Cache = mockRaCache()
	cmds, err := pkg.CreateRACommands(mockGardenConfig())
	expectedCmd := pkg.Command{
		CmdType: pkg.ReactiveActuator,
		Id:      raId,
		Cmd:     1,
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
	cache, err := pkg.UpdateRACache(pkg.Cache, raConfigs, mockClient, mockStorage)
	assert.NoError(t, err)
	assert.Equal(t, cache, mockRaCache())
}
