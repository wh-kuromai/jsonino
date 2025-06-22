package jsonino

import "encoding/json"

type ObjectScheme struct {
	TypeName   string                     `json:"type"`
	Properties map[string]*SchemaNodeBase `json:"properties,omitempty"`
	Required   []string                   `json:"required,omitempty"`
	Order      []string                   `json:"order,omitempty"`
}

func ObjectSchemeFactory(data []byte) (SchemaNode, error) {
	s := &ObjectScheme{}
	err := json.Unmarshal(data, s)
	return s, err
}

func (s *ObjectScheme) Type() string {
	return s.TypeName
}

func (s *ObjectScheme) Parse(pr PathResolver, buf []byte) (*Node, error) {
	m := make(map[string]json.RawMessage)
	err := json.Unmarshal(buf, &m)
	if err != nil {
		return nil, err
	}

	nm := make(map[string]*Node)
	for key, schema := range s.Properties {
		val := m[key]

		n, err2 := schema.Parse(pr, val)
		if err2 != nil {
			return nil, err2
		}

		nm[key] = n
	}

	return &Node{
		Type:        "object",
		ObjectValue: nm,
	}, nil
}

func (s *ObjectScheme) ValidateNode(pr PathResolver, n *Node) bool {
	if n.Type != "object" || n.ObjectValue == nil {
		return false
	}

	if s.Properties != nil {
		for key, scheme := range s.Properties {
			vnode := n.ObjectValue[key]
			if !scheme.ValidateNode(pr, vnode) {
				return false
			}
		}
	}

	return true
}

func (s *ObjectScheme) Validate(pr PathResolver, buf []byte) bool {
	n, err := s.Parse(pr, buf)
	if err != nil {
		return false
	}

	return s.ValidateNode(pr, n)
}
