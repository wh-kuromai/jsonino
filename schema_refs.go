package jsonino

import (
	"encoding/json"
)

type RefsSchema struct {
	Refs *string `json:"refs,omitempty"`
}

func RefsSchemaFactory(data []byte) (SchemaNode, error) {
	s := &RefsSchema{}
	err := json.Unmarshal(data, s)
	return s, err
}

func (s *RefsSchema) Type() string {
	return "refs"
}

func (s *RefsSchema) Parse(pr PathResolver, buf []byte) (*Node, error) {
	return pr(*s.Refs).Parse(pr, buf)
}

func (s *RefsSchema) ValidateNode(pr PathResolver, n *Node) bool {
	return pr(*s.Refs).ValidateNode(pr, n)
}

func (s *RefsSchema) Validate(pr PathResolver, buf []byte) bool {
	return pr(*s.Refs).Validate(pr, buf)
}
