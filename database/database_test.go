package database

import "testing"

func TestInitialize(t *testing.T) {
	// initialize database in test mode
	Initialize(true)

	if Engine.Ping() != nil {
		t.Errorf("got: %v\nexpected: %v", Engine.Ping(), nil)
	}
}
