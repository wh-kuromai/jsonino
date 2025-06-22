package jsonino

import (
	"encoding/json"
)

func NewSchema(buf []byte) (*Schema, error) {
	scm := &Schema{}
	err := json.Unmarshal(buf, scm)
	if err != nil {
		return nil, err
	}

	return scm, nil
}

type PathResolver func(path string) Validator

var sSchemeFactory map[string]func(data []byte) (schemaNode, error)

func init() {
	sSchemeFactory = make(map[string]func(data []byte) (schemaNode, error))

	sSchemeFactory["string"] = stringSchemeFactory
	sSchemeFactory["number"] = numberSchemeFactory
	sSchemeFactory["array"] = arraySchemeFactory
	sSchemeFactory["object"] = objectSchemeFactory

}

func getFactory(typ string) func(data []byte) (schemaNode, error) {
	return sSchemeFactory[typ]
}

type Validator interface {
	Type() string
	Validate(buf []byte) bool
}

type schemaNode interface {
	Validator

	validateNode(pr PathResolver, node *Node) bool
	parse(pr PathResolver, buf []byte) (*Node, error)
}
