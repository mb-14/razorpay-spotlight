package main

import (
	"strings"
	"time"

	"github.com/buger/jsonparser"
)

type Json struct {
	data []byte
}

func (j Json) GetString(path string) (string, error) {
	return jsonparser.GetString(j.data, strings.Split(path, ".")...)
}

func (j Json) GetInt(path string) (int64, error) {
	return jsonparser.GetInt(j.data, strings.Split(path, ".")...)
}

func (j Json) GetTime(path string) (time.Time, error) {
	epoch, err := jsonparser.GetInt(j.data, strings.Split(path, ".")...)
	if err != nil {
		return time.Now(), err
	}
	return time.Unix(epoch, 0), nil
}
