package greetings

import (
	"testing"
)

func TestHelloWorld(t *testing.T) {
	if "Hello" == "World!" {
		t.Fatalf("The world is broken!")
	}
}
