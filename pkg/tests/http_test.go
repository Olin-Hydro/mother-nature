package pkg

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/Olin-Hydro/mother-nature/mocks"
	"github.com/Olin-Hydro/mother-nature/pkg"
	"github.com/stretchr/testify/assert"
)

func TestDecodeJson(t *testing.T) {
	b, err := json.Marshal(mocks.MockGarden())
	assert.NoError(t, err)
	garden := pkg.Garden{}
	r := io.NopCloser(bytes.NewReader(b))
	err = pkg.DecodeJson(&garden, r)
	assert.NoError(t, err)
	assert.Equal(t, mocks.MockGarden(), garden)
}
