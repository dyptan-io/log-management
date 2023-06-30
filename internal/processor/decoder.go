package processor

import (
	"encoding/json"

	"github.com/diptanw/log-management/api"
)

type DecoderJSON struct{}

func (DecoderJSON) Decode(b []byte) (api.Log, error) {
	var raw struct {
		Id string `json:"id"`
		L  string `json:"@l"`
		M  string `json:"@m"`
		T  string `json:"@t"`
	}

	if err := json.Unmarshal(b, &raw); err != nil {
		return api.Log{}, err
	}

	return api.Log{
		Id:        raw.Id,
		Severity:  raw.L,
		Message:   raw.M,
		Timestamp: raw.T,
	}, nil
}
