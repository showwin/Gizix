package model

import db "github.com/showwin/Gizix/database"

// SkyWay model
type SkyWay struct {
	Key string
}

// UpdateSkyWayKey : update skyway API key
func UpdateSkyWayKey(key string) bool {
	_, err := db.Engine.Exec("UPDATE skyway SET api_key = ? WHERE id = 1", key)
	return err == nil
}

// GetSkyWayKey : get skyway API key
func GetSkyWayKey() (s SkyWay) {
	db.Engine.QueryRow("SELECT api_key FROM skyway WHERE id = 1 LIMIT 1").Scan(&s.Key)
	return s
}
