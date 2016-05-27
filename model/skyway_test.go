package model

import (
	"testing"
	"time"

	db "github.com/showwin/Gizix/database"
)

func TestUpdateSkyWayKey(t *testing.T) {
	db.Initialize(true)
	newKey := "Gizix" + time.Now().String()
	UpdateSkyWayKey(newKey)

	var s SkyWay
	db.Engine.QueryRow("SELECT api_key FROM skyway WHERE id = 1 LIMIT 1").Scan(&s.Key)
	if s.Key != newKey {
		t.Errorf("got: %v\nexpected: %v", s.Key, newKey)
	}
}

func TestGetSkyWayKey(t *testing.T) {
	db.Initialize(true)
	newKey := "Gizix" + time.Now().String()
	UpdateSkyWayKey(newKey)

	s := GetSkyWayKey()
	if s.Key != newKey {
		t.Errorf("got: %v\nexpected: %v", s.Key, newKey)
	}
}
