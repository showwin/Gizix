package main

// SkyWay model
type SkyWay struct {
	Key string
}

func updateSkyWayKey(key string) bool {
	_, err := db.Exec("UPDATE skyway SET api_key = ? WHERE id = 1", key)
	return err == nil
}

func getSkyWayKey() (s SkyWay) {
	db.QueryRow("SELECT api_key FROM skyway WHERE id = 1 LIMIT 1").Scan(&s.Key)
	return s
}
