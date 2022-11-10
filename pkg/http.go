package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

func DecodeJson(destination interface{}, body io.ReadCloser) error {
	err := json.NewDecoder(body).Decode(&destination)
	if err != nil {
		return fmt.Errorf("#DecodeJson: %e", err)
	}
	return nil
}
