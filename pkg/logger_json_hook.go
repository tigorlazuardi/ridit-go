package pkg

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type JSONHook struct{}

func (JSONHook) Fire(entry *logrus.Entry) error {
	if _, ok := entry.Logger.Formatter.(*logrus.JSONFormatter); ok {
		return nil
	}
	for k, v := range entry.Data {
		entry.Data[k] = jsonize(v)
	}
	return nil
}

func (JSONHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func jsonize(v interface{}) interface{} {
	switch v := v.(type) {
	case []byte:
		return string(v)
	case string:
		return v
	}
	if val, err := json.Marshal(v); err == nil {
		if string(val) == "{}" {
			if e, ok := v.(error); ok {
				return e.Error()
			}
		}
		return string(val)
	}
	return v
}
