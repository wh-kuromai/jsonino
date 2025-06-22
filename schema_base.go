package jsonino

import (
	"encoding/json"

	"github.com/goccy/go-yaml"
)

type SchemaNodeBase struct {
	TypeName string     `json:"type"`
	This     SchemaNode `json:"-"`
}

func (s *SchemaNodeBase) Type() string {
	return s.This.Type()
}

func (s *SchemaNodeBase) Parse(pr PathResolver, buf []byte) (*Node, error) {
	return s.This.Parse(pr, buf)
}

func (s *SchemaNodeBase) ValidateNode(pr PathResolver, node *Node) bool {
	return s.This.ValidateNode(pr, node)
}

func (s *SchemaNodeBase) Validate(pr PathResolver, buf []byte) bool {
	return s.This.Validate(pr, buf)
}

func (s *SchemaNodeBase) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(s.This)
}

func (s *SchemaNodeBase) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.This)
}

func (s *SchemaNodeBase) UnmarshalJSON(data []byte) error {
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
