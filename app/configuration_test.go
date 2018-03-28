package app

import (
	"testing"
)

func TestParseUsers(t *testing.T) {
	fixtures := map[string]UserList{
		"":                      {},
		"foo:bar":               {"foo": "bar"},
		"user1:pwd1;user2:pwd2": {"user1": "pwd1", "user2": "pwd2"},
	}

	for input, expected := range fixtures {
		output, _ := ParseUsers(input)

		if len(output) != len(expected) ||
			output["user1"] != expected["user1"] ||
			output["user2"] != expected["user2"] {
			t.Errorf("expected %s - got %s", expected, output)
		}
	}
}

func TestParseUsersWhenInvalid(t *testing.T) {
	_, err := ParseUsers("foo:bar;blah")
	if err == nil {
		t.Error("Invalid input but no error returned")
	}
}
