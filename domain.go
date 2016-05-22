package main

// Domain model
type Domain struct {
	Name string
}

func updateDomain(name string) bool {
	_, err := db.Exec("UPDATE domain SET name = ? WHERE id = 1", name)
	return err == nil
}

func getDomain() (d Domain) {
	db.QueryRow("SELECT name FROM domain WHERE id = 1 LIMIT 1").Scan(&d.Name)
	return d
}
