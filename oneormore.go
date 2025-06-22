package jsonino

import (
	"encoding/json"

	"github.com/goccy/go-yaml"
)

type OneOrMore[T any] struct {
	One  *T
	More []*T
}

func (o *OneOrMore[T]) Iter() chan<- *T {
	if o.More != nil {
		ch := make(chan *T, len(o.More))
		for _, v := range o.More {
			ch <- v
		}
		return ch

	}

	if o.One != nil {
		ch := make(chan *T, 1)
		ch <- o.One
		return ch
	}

	return make(chan *T)
}

func (o *OneOrMore[T]) Do(f func(t *T) bool) bool {
	if o.More != nil {
		for _, c := range o.More {
			if f(c) {
				return true
			}
		}
		return false
	}

	if o.One != nil {
		return f(o.One)
	}

	return false
}

func (o *OneOrMore[T]) MarshalJSON() ([]byte, error) {
	if o.More != nil {
		return json.Marshal(o.More)
	}

	return json.Marshal(o.One)
}

func (o *OneOrMore[T]) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &o.More)
	if err == nil {
		return nil
	}

	one := new(T)
	o.One = one
	return json.Unmarshal(data, one)
}

func (o *OneOrMore[T]) MarshalYAML() ([]byte, error) {
	if o.More != nil {
		return yaml.Marshal(o.More)
	}

	return yaml.Marshal(o.One)
}

func (o *OneOrMore[T]) UnmarshalYAML(data []byte) error {
	err := yaml.Unmarshal(data, &o.More)
	if err == nil {
		return nil
	}

	one := new(T)
	o.One = one
	return yaml.Unmarshal(data, one)
}
