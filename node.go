package jsonino

func matchRate(s1, s2 schemaNode) (numMatch int, rate float64) {
	if s1.Type() != s2.Type() {
		return 1, 0
	}

	rate = 1.0

	return numMatch, rate
}

type Node struct {
	Type        string
	BoolValue   *bool
	StringValue *string
	NumberValue *float64
	ArrayValue  []*Node
	ObjectValue map[string]*Node
}
