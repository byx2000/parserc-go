package json

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	json := `
	{
		"a": 123,
		"b": 3.14,
		"c": "hello",
		"d": {
			"x": 100,
			"y": "world!"
		},
		"e": [
			12,
			34.56,
			{
				"name": "Xiao Ming",
				"age": 18,
				"score": [99.8, 87.5, 60.0]
			},
			"abc"
		],
		"f": [],
		"g": {},
		"h": [true, {"m": false}]
	}`
	m := map[string]interface{}{
		"a": 123,
		"b": 3.14,
		"c": "hello",
		"d": map[string]interface{}{
			"x": 100,
			"y": "world!",
		},
		"e": []interface{}{
			12,
			34.56,
			map[string]interface{}{
				"name":  "Xiao Ming",
				"age":   18,
				"score": []interface{}{99.8, 87.5, 60.0},
			},
			"abc",
		},
		"f": []interface{}{},
		"g": map[string]interface{}{},
		"h": []interface{}{true, map[string]interface{}{"m": false}},
	}

	r := Parse(json)
	assert.Equal(t, m, r)
}
