package tests

import (
	"testing"

	"github.com/Olin-Hydro/mother-nature/pkg"
	"github.com/stretchr/testify/assert"
)

const (
	BaseURL     = "https://test.com"
	GardenRoute = "/garden"
	RaLogRoute  = "/ra/logging/actions"
	RaRoute     = "/ra"
	SensorRoute = "/sensor"
)

func MockHydrangea() (pkg.Hydrangea, error) {
	return pkg.NewHydrangea(BaseURL+GardenRoute, BaseURL+RaLogRoute, BaseURL+RaRoute, BaseURL+SensorRoute)
}

func TestNewHydrangea(t *testing.T) {
	h, err := MockHydrangea()
	assert.NoError(t, err)
	assert.Equal(t, GardenRoute, h.GardenURL.Path)
	assert.Equal(t, RaLogRoute, h.RALogURL.Path)
	assert.Equal(t, RaRoute, h.RAURL.Path)
	assert.Equal(t, SensorRoute, h.SensorLogURL.Path)
}

func TestCreateGardenReq(t *testing.T) {
	h, err := MockHydrangea()
	assert.NoError(t, err)
	id := "abc"
	req, err := h.CreateGardenReq(id)
	assert.NoError(t, err)
	assert.Equal(t, BaseURL+GardenRoute+"/"+id, req.URL.String())
}

func TestCreateRALogsReq(t *testing.T) {
	h, err := MockHydrangea()
	assert.NoError(t, err)
	id := "abc"
	limit := "10"
	req, err := h.CreateRALogsReq(id, limit)
	assert.NoError(t, err)
	assert.Equal(t, BaseURL+RaLogRoute+"/"+id+"?limit="+limit, req.URL.String())
}

func TestCreateRAReq(t *testing.T) {
	h, err := MockHydrangea()
	assert.NoError(t, err)
	id := "abc"
	req, err := h.CreateRAReq(id)
	assert.NoError(t, err)
	assert.Equal(t, BaseURL+RaRoute+"/"+id, req.URL.String())
}

func TestCreateSensorLogsReq(t *testing.T) {
	h, err := MockHydrangea()
	assert.NoError(t, err)
	id := "abc"
	limit := "10"
	req, err := h.CreateSensorLogsReq(id, limit)
	assert.NoError(t, err)
	assert.Equal(t, BaseURL+SensorRoute+"/"+id+"?limit="+limit, req.URL.String())
}
