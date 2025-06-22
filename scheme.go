package jsonino

import (
	"encoding/json"
)

type Schema interface {
	Validate(buf []byte) bool
}

type SchemaBase struct {
	Root     SchemaNode
	Resolver PathResolver
}

func (s *SchemaBase) Validate(buf []byte) bool {
	return s.Root.Validate(s.Resolver, buf)
}

func NewSchema(buf []byte) (Schema, error) {
	scm := &SchemaNodeBase{}
	err := json.Unmarshal(buf, scm)
	if err != nil {
		return nil, err
	}

	return &SchemaBase{
		Root: scm.This,
		Resolver: func(path string) SchemaNode {
			if path == "#" {
				return scm.This
			}
			return nil
		},
	}, nil
}

type PathResolver func(path string) SchemaNode

var sSchemeFactory map[string]func(data []byte) (SchemaNode, error)

func init() {
	sSchemeFactory = make(map[string]func(data []byte) (SchemaNode, error))

	sSchemeFactory["string"] = StringSchemeFactory
	sSchemeFactory["number"] = NumberSchemeFactory
	sSchemeFactory["array"] = ArraySchemeFactory
	sSchemeFactory["object"] = ObjectSchemeFactory

}

func getFactory(typ string) func(data []byte) (SchemaNode, error) {
	return sSchemeFactory[typ]
}

type SchemaNode interface {
	Type() string
	Validate(pr PathResolver, buf []byte) bool
	ValidateNode(pr PathResolver, node *Node) bool
	Parse(pr PathResolver, buf []byte) (*Node, error)
}
