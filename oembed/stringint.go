package oembed

import (
	"encoding/json"
	"strconv"
)

// A StringInt holds an int value that can be represented as a string
// We need this workaround for width and height
// Flickr uses strings, and Vimeo integers
type StringInt int

// UnmarshalJSON returns the parsed JSON value of StringInt
func (s *StringInt) UnmarshalJSON(b []byte) error {
	// try an int
	var intVal int
	// watch the == (we enter this condition if no error)
	if err := json.Unmarshal(b, &intVal); err == nil {
		*s = StringInt(intVal)
		return nil
	}

	// now, try a string
	var stringVal string
	if err := json.Unmarshal(b, &stringVal); err != nil {
		return err
	}

	intVal, err := strconv.Atoi(stringVal)
	if err != nil {
		return err
	}
	*s = StringInt(intVal)
	return nil
}
