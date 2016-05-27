package model

import (
	"testing"
	"time"

	db "github.com/showwin/Gizix/database"
)

func TestUpdateDomain(t *testing.T) {
	db.Initialize(true)
	newDomain := time.Now().String() + ".com"
	UpdateDomain(newDomain)

	var d Domain
	db.Engine.QueryRow("SELECT name FROM domain WHERE id = 1 LIMIT 1").Scan(&d.Name)
	if d.Name != newDomain {
		t.Errorf("got: %v\nexpected: %v", d.Name, newDomain)
	}
}

func TestGetDomain(t *testing.T) {
	db.Initialize(true)
	newDomain := time.Now().String() + ".com"
	UpdateDomain(newDomain)

	d := GetDomain()
	if d.Name != newDomain {
		t.Errorf("got: %v\nexpected: %v", d.Name, newDomain)
	}
}
