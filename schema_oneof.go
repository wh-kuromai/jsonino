package jsonino

import (
	"encoding/json"
	"errors"
)

type OneOfSchema struct {
	OneOf []*Schema `json:"oneOf,omitempty"`
}

func oneOfSchemaFactory(data []byte) (schemaNode, error) {
	s := &OneOfSchema{}
	err := json.Unmarshal(data, s)
	return s, err
}

func (s *OneOfSchema) Type() string {
	return "oneOf"
}

func (s *OneOfSchema) parse(pr PathResolver, buf []byte) (*Node, error) {
	for _, scm := range s.OneOf {
		if scm.Validate(buf) {
			return scm.Parse(pr, buf)
		}
	}

	return nil, errors.New("oneof parse error")
}

func (s *OneOfSchema) validateNode(pr PathResolver, n *Node) bool {
	for _, scm := range s.OneOf {
		if scm.ValidateNode(pr, n) {
			return true
		}
	}

	return false
}

func (s *OneOfSchema) Validate(buf []byte) bool {
	for _, scm := range s.OneOf {
		if scm.Validate(buf) {
			n, err := scm.Parse(nil, buf)
			if err != nil {
				return false
			}
			return s.validateNode(nil, n)
		}
	}

	return false
}
