package model

import db "github.com/showwin/Gizix/database"

// Domain model
type Domain struct {
	Name string
}

// UpdateDomain : update domain name
func UpdateDomain(name string) bool {
	_, err := db.Engine.Exec("UPDATE domain SET name = ? WHERE id = 1", name)
	return err == nil
}

// GetDomain : get domain name
func GetDomain() (d Domain) {
	db.Engine.QueryRow("SELECT name FROM domain WHERE id = 1 LIMIT 1").Scan(&d.Name)
	return d
}
