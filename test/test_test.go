package greetings

import (
	"testing"
)

func TestHelloName(t *testing.T) {
	if "Hello" == "World!" {
		t.Fatalf("The world is broken!")
	}
}
