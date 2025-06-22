package jsonino

import (
	"encoding/json"
	"errors"
)

type OneOfSchema struct {
	OneOf []*SchemaNodeBase `json:"oneOf,omitempty"`
}

func OneOfSchemaFactory(data []byte) (SchemaNode, error) {
	s := &OneOfSchema{}
	err := json.Unmarshal(data, s)
	return s, err
}

func (s *OneOfSchema) Type() string {
	return "oneOf"
}

func (s *OneOfSchema) Parse(pr PathResolver, buf []byte) (*Node, error) {
	for _, scm := range s.OneOf {
		if scm.Validate(pr, buf) {
			return scm.Parse(pr, buf)
		}
	}

	return nil, errors.New("oneof parse error")
}

func (s *OneOfSchema) ValidateNode(pr PathResolver, n *Node) bool {
	for _, scm := range s.OneOf {
		if scm.ValidateNode(pr, n) {
			return true
		}
	}

	return false
}

func (s *OneOfSchema) Validate(pr PathResolver, buf []byte) bool {
	for _, scm := range s.OneOf {
		if scm.Validate(pr, buf) {
			n, err := scm.Parse(pr, buf)
			if err != nil {
				return false
			}
			return s.ValidateNode(pr, n)
		}
	}

	return false
}
