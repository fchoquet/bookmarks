package oembed

import (
	"encoding/json"
	"testing"
)

func TestStringInt(t *testing.T) {
	t.Run("it unmarshals json from both int and strings", func(t *testing.T) {

		fixtures := []struct {
			input    string
			expected int
			hasError bool
		}{
			// Value is a int
			{`{"Value":1234}`, 1234, false},
			{`{"Value":-1234}`, -1234, false},

			// Value is a string
			{`{"Value":"1234"}`, 1234, false},
			{`{"Value":"-1234"}`, -1234, false},
			{`{"Value":"blah"}`, 0, true},

			// Let's accept null
			{`{"Value":null}`, 0, false},

			// Value is another type
			{`{"Value":true}`, 0, true},
			{`{"Value":{}}`, 0, true},
			{`{"Value":[]}`, 0, true},
		}

		for _, f := range fixtures {
			var v struct {
				Value StringInt
			}
			err := json.Unmarshal([]byte(f.input), &v)

			if f.hasError && err == nil {
				t.Errorf("an error was expected with entry: %+v - returned %d instead", f, int(v.Value))
			}

			if err != nil && !f.hasError {
				t.Errorf("unexepected error: %s", err)
			}

			if int(v.Value) != f.expected {
				t.Errorf("expected %d, got %d", f.expected, int(v.Value))
			}
		}
	})
}
