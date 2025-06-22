package jsonino

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/goccy/go-yaml"
)

type UnixTime struct {
	t time.Time
}

func (o UnixTime) Time() time.Time {
	return (time.Time)(o.t)
}

func (o UnixTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Time().Unix())
}

func (o *UnixTime) UnmarshalJSON(data []byte) error {
	//fmt.Println("UnmarshalJSON")
	var num int64
	err := json.Unmarshal(data, &num)
	if err == nil {
		o.t = time.Unix(num, 0)
		return nil
	}

	var str string
	err = json.Unmarshal(data, &str)
	if err == nil {
		o.t = ParseAnyTime(str)
		return nil
	}

	return errors.New("UnmarshalJSON failed")
}

func (o UnixTime) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(o.Time().Unix())
}

func (o *UnixTime) UnmarshalYAML(data []byte) error {
	var num int64
	err := yaml.Unmarshal(data, &num)
	if err == nil {
		o.t = time.Unix(num, 0)
		return nil
	}

	var str string
	err = yaml.Unmarshal(data, &str)
	if err == nil {
		o.t = ParseAnyTime(str)
		return nil
	}

	return errors.New("UnmarshalYAML failed")
}
