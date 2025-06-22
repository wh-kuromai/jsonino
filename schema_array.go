package jsonino

import "encoding/json"

type ArrayScheme struct {
	TypeName string  `json:"type"`
	Items    *Schema `json:"items,omitempty"`
	MaxItems *int    `json:"maxItems,omitempty"`
	MinItems *int    `json:"minItems,omitempty"`
	//UniqueItems *bool      `json:"uniqueItems,omitempty"`
	//MaxContains *int       `json:"exclusiveMinimum,omitempty"`
	//MultipleOf  *int       `json:"multipleOf,omitempty"`
	//Pattern     *string    `json:"pattern,omitempty"`
}

func (s *ArrayScheme) Type() string {
	return s.TypeName
}

func arraySchemeFactory(data []byte) (schemaNode, error) {
	s := &ArrayScheme{}
	err := json.Unmarshal(data, s)
	return s, err
}

func (s *ArrayScheme) parse(pr PathResolver, buf []byte) (*Node, error) {
	ary := make([]json.RawMessage, 0, 10)
	err := json.Unmarshal(buf, &ary)
	if err != nil {
		return nil, err
	}

	l := len(ary)
	nary := make([]*Node, l)

	for i := range ary {
		n, err2 := s.Items.Parse(pr, ary[i])
		if err2 != nil {
			return nil, err2
		}

		nary[i] = n
	}

	return &Node{
		Type:       "array",
		ArrayValue: nary,
	}, nil
}

func (s *ArrayScheme) validateNode(pr PathResolver, n *Node) bool {
	if n.Type != "array" || n.ArrayValue == nil {
		return false
	}

	l := len(n.ArrayValue)
	if s.MaxItems != nil {
		if l > *s.MaxItems {
			return false
		}
	}

	if s.MinItems != nil {
		if l < *s.MinItems {
			return false
		}
	}

	if s.Items != nil {

		for _, val := range n.ArrayValue {
			if !s.Items.ValidateNode(pr, val) {
				return false
			}
		}
	}

	return true
}

func (s *ArrayScheme) Validate(buf []byte) bool {

	n, err := s.parse(nil, buf)
	if err != nil {
		return false
	}

	return s.validateNode(nil, n)
}
