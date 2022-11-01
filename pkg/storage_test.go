package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHydrangea(t *testing.T) {
	h, err := NewHydrangea("https://test.com/garden", "https://test.com/ra", "https://test.com/sensor")
	assert.NoError(t, err)
	assert.Equal(t, "/garden", h.GardenURL.Path)
	assert.Equal(t, "/ra", h.RALogURL.Path)
	assert.Equal(t, "/sensor", h.SensorLogURL.Path)
}
