package jsonino

import (
	"encoding/json"

	"github.com/goccy/go-yaml"
)

type Schema struct {
	TypeName string     `json:"type"`
	This     schemaNode `json:"-"`
}

func (s *Schema) Type() string {
	return s.This.Type()
}

func (s *Schema) Parse(pr PathResolver, buf []byte) (*Node, error) {
	return s.This.parse(pr, buf)
}

func (s *Schema) ValidateNode(pr PathResolver, node *Node) bool {
	return s.This.validateNode(pr, node)
}

func (s *Schema) Validate(buf []byte) bool {
	return s.This.Validate(buf)
}

func (s *Schema) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(s.This)
}

func (s *Schema) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.This)
}

func (s *Schema) UnmarshalJSON(data []byte) error {
	a := &struct {
		TypeName *string         `json:"type"`
		OneOf    json.RawMessage `json:"oneOf,omitempty"`
		Refs     json.RawMessage `json:"refs,omitempty"`
	}{}

	err := json.Unmarshal(data, a)
	if err != nil {
		return err
	}

	if a.OneOf != nil {
		sc := &OneOfSchema{}
		err = json.Unmarshal(data, sc)
		if err != nil {
			return err
		}

		s.TypeName = sc.Type()
		s.This = sc
		return nil
	}

	factory := getFactory(*a.TypeName)

	sc, err := factory(data)
	if err != nil {
		return err
	}

	s.TypeName = *a.TypeName
	s.This = sc
	return nil
}
