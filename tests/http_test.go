package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/Olin-Hydro/mother-nature/pkg"
	"github.com/stretchr/testify/assert"
)

func TestDecodeJson(t *testing.T) {
	b, err := json.Marshal(mockGarden())
	assert.NoError(t, err)
	garden := pkg.Garden{}
	r := ioutil.NopCloser(bytes.NewReader(b))
	err = pkg.DecodeJson(&garden, r)
	assert.NoError(t, err)
	assert.Equal(t, mockGarden(), garden)
}
