package database_test

import (
	"os"
	"testing"
)

func setup() {

}

func shutdown() {

}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	shutdown()
	os.Exit(code)
}

func TestDemo(t *testing.T) {
	t.Error("Beep boop! Something went terribly wrong!")
}
