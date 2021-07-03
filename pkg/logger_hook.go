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
	if v, ok := v.([]byte); ok {
		return string(v)
	}
	if val, err := json.Marshal(v); err == nil {
		return string(val)
	}
	return v
}
