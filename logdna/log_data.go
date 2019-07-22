package logdna

import (
	"encoding/json"
	"time"
)

type LogData struct {
	Time    time.Time
	Message string

	Level string
	App   string
	Env   string
	Meta  map[string]interface{}
}

func (d LogData) GetTime() int64 {
	if d.Time.IsZero() {
		return time.Now().UnixNano() / int64(time.Millisecond)
	}
	return d.Time.UnixNano() / int64(time.Millisecond)
}

func (d LogData) toPayload() *logPayload {
	return &logPayload{
		Timestamp: d.GetTime(),
		Line:      d.Message,
		Level:     d.Level,
		App:       d.App,
		Env:       d.Env,
		Meta:      d.Meta,
	}
}

type logPayload struct {
	Timestamp int64  `json:"timestamp"`
	Line      string `json:"line"`

	Hostname string                 `json:"hostname,omitempty"`
	Level    string                 `json:"level,omitempty"`
	App      string                 `json:"app,omitempty"`
	Env      string                 `json:"env,omitempty"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
}

func logsToJSON(logs []*logPayload) (string, error) {
	byt, err := json.Marshal(map[string]interface{}{
		"lines": logs,
	})
	return string(byt), err
}
