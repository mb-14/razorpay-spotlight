package json

import (
	"strings"
	"time"

	"github.com/buger/jsonparser"
)

type Json struct {
	Data []byte
}

func (j Json) GetString(path string) (string, error) {
	return jsonparser.GetString(j.Data, strings.Split(path, ".")...)
}

func (j Json) GetInt(path string) (int64, error) {
	return jsonparser.GetInt(j.Data, strings.Split(path, ".")...)
}

func (j Json) GetTime(path string) (time.Time, error) {
	epoch, err := jsonparser.GetInt(j.Data, strings.Split(path, ".")...)
	if err != nil {
		return time.Now(), err
	}
	return time.Unix(epoch, 0), nil
}

func (j *Json) Set(path string, value []byte) error {
	var err error
	j.Data, err = jsonparser.Set(j.Data, value, strings.Split(path, ".")...)
	return err
}
