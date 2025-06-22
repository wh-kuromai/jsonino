package jsonino

import (
	"encoding/json"
	"regexp"
)

type StringScheme struct {
	TypeName  string         `json:"type"`
	Value     *string        `json:"value,omitempty"`
	MaxLength *int           `json:"maxLength,omitempty"`
	MinLength *int           `json:"minLength,omitempty"`
	Pattern   *string        `json:"pattern,omitempty"`
	cpattern  *regexp.Regexp `json:"-"`
}

func StringSchemeFactory(data []byte) (SchemaNode, error) {
	s := &StringScheme{}
	err := json.Unmarshal(data, s)
	return s, err
}

func (s *StringScheme) Type() string {
	return s.TypeName
}

func (s *StringScheme) Parse(pr PathResolver, buf []byte) (*Node, error) {
	var str string
	err := json.Unmarshal(buf, &str)
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:        "string",
		StringValue: &str,
	}, nil
}

func (s *StringScheme) ValidateNode(pr PathResolver, n *Node) bool {
	if n.Type != "string" {
		return false
	}

	if n.StringValue == nil {
		return false
	}

	l := len(*n.StringValue)
	if s.MaxLength != nil {
		if l > *s.MaxLength {
			return false
		}
	}

	if s.MinLength != nil {
		if l < *s.MinLength {
			return false
		}
	}

	if s.Pattern != nil {
		if s.cpattern == nil {
			s.cpattern = regexp.MustCompile(*s.Pattern)
		}

		if !s.cpattern.MatchString(*n.StringValue) {
			return false
		}
	}

	return true
}

func (s *StringScheme) Validate(pr PathResolver, buf []byte) bool {

	n, err := s.Parse(pr, buf)
	if err != nil {
		return false
	}

	return s.ValidateNode(pr, n)
}
