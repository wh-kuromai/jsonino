package jsonino

import (
	"encoding/json"
	"regexp"

	"github.com/goccy/go-yaml"
)

type Regexp struct {
	Pattern string
	re      *regexp.Regexp `json:"-"`
}

func (o *Regexp) Regexp() *regexp.Regexp {
	if o.re == nil {
		o.re = regexp.MustCompile(o.Pattern)
	}
	return o.re
}

func (o *Regexp) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Pattern)
}

func (o *Regexp) UnmarshalJSON(data []byte) (err error) {
	return json.Unmarshal(data, &o.Pattern)
}

func (o *Regexp) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(o.Pattern)
}

func (o *Regexp) UnmarshalYAML(data []byte) (err error) {
	return yaml.Unmarshal(data, &o.Pattern)
}
