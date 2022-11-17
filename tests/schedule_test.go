package tests

import (
	"testing"

	"github.com/Olin-Hydro/mother-nature/pkg"
	"github.com/stretchr/testify/assert"
)

func TestNewSchedule(t *testing.T) {
	gardenConfig := mockGardenConfig()
	sched, err := pkg.NewSchedule(gardenConfig)
	assert.NoError(t, err)
	assert.Equal(t, mockSchedule(), sched)
}
