package processor

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/dyptan-io/log-management/v2/api"
)

type (
	DecoderJSON struct{}
	raw         map[string]any
)

func (DecoderJSON) Decode(b []byte) (api.Log, error) {
	var r raw

	if err := json.Unmarshal(b, &r); err != nil {
		return api.Log{}, err
	}

	log := api.Log{
		Id:        r.string("id"),
		Severity:  r.string("@l"),
		Message:   r.string("@m"),
		Timestamp: r.time("@t"),
	}

	// Set remaining fields as extra attributes.
	log.Attributes = r

	return log, nil
}

func (r raw) string(name string) string {
	if str, ok := r[name]; ok {
		delete(r, name)
		return str.(string)
	}

	return ""
}

func (r raw) time(name string) time.Time {
	rawT := r.string(name)
	if rawT == "" {
		return time.Time{}
	}

	// TODO: Use custom layout parsing instead of trimming microseconds.
	t, err := time.Parse(time.DateTime, rawT[:strings.LastIndex(rawT, ":")])
	if err != nil {
		return time.Time{}
	}

	return t
}
