package jsonino

import "encoding/json"

type NumberScheme struct {
	TypeName     string   `json:"type"`
	Value        *float64 `json:"value,omitempty"`
	Max          *float64 `json:"maximum,omitempty"`
	ExclusiveMax *float64 `json:"exclusiveMaximum,omitempty"`
	Min          *float64 `json:"minimum,omitempty"`
	ExclusiveMin *float64 `json:"exclusiveMinimum,omitempty"`
	MultipleOf   *int     `json:"multipleOf,omitempty"`
	//Pattern      *string  `json:"pattern,omitempty"`
}

func numberSchemeFactory(data []byte) (schemaNode, error) {
	s := &NumberScheme{}
	err := json.Unmarshal(data, s)
	return s, err
}

func (s *NumberScheme) Type() string {
	return s.TypeName
}

func (s *NumberScheme) parse(pr PathResolver, buf []byte) (*Node, error) {
	var num float64
	err := json.Unmarshal(buf, &num)
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:        "number",
		NumberValue: &num,
	}, nil
}

func (s *NumberScheme) validateNode(pr PathResolver, n *Node) bool {
	if n.Type != "number" {
		return false
	}

	if s.Max != nil {
		if *n.NumberValue > *s.Max {
			return false
		}
	}

	if s.ExclusiveMax != nil {
		if *n.NumberValue >= *s.ExclusiveMax {
			return false
		}
	}

	if s.Min != nil {
		if *n.NumberValue < *s.Min {
			return false
		}
	}

	if s.ExclusiveMin != nil {
		if *n.NumberValue <= *s.ExclusiveMin {
			return false
		}
	}

	if s.MultipleOf != nil {
		valint := int64(*n.NumberValue)
		if *n.NumberValue != float64(valint) {
			return false
		}

		if valint%int64(*s.MultipleOf) != 0 {
			return false
		}
	}

	return true
}

func (s *NumberScheme) Validate(buf []byte) bool {

	n, err := s.parse(nil, buf)
	if err != nil {
		return false
	}

	return s.validateNode(nil, n)
}
