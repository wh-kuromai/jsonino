package jsonino

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbs(t *testing.T) {

	dat := `
		{"a":1}
	`
	scm := `
		{
			"type":"object",
			"properties": {
				"a" : {
					"type": "number"
				}
			}
		}
	`

	scheme, err := NewSchema([]byte(scm))
	if err != nil {
		t.Fatal(err)
	}

	result := scheme.Validate([]byte(dat))
	assert.EqualValues(t, true, result)

}

func TestOneOf(t *testing.T) {

	dat := `
		{"a":1}
	`
	scm := `
		{
			"oneOf" : [
				{
					"type":"object",
					"properties": {
						"a" : {
							"type": "number"
						}
					}
				},
				{
					"type":"object",
					"properties": {
						"b" : {
							"type": "number"
						}
					}
				}
			]
		}
	`

	scheme, err := NewSchema([]byte(scm))
	if err != nil {
		t.Fatal(err)
	}

	result := scheme.Validate([]byte(dat))
	assert.EqualValues(t, true, result)

}
