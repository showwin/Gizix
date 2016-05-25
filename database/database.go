package database

import (
	"database/sql"
	"log"
	"os"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// Engine of database
var Engine *sql.DB

// Initialize database setting
func Initialize(testEnv bool) {
	// override settings depending on test env
	onCircleCI := os.Getenv("CIRCLECI")
	if onCircleCI == "true" {
		os.Setenv("DATABASE_USERNAME", "ubuntu")
		os.Setenv("DATABASE_NAME", "circle_test")
	} else if testEnv {
		os.Setenv("DATABASE_NAME", "gizix_test")
	}

	// read settings
	username := os.Getenv("DATABASE_USERNAME")
	dbName := os.Getenv("DATABASE_NAME")

	var err error
	Engine, err = sql.Open("mysql", username+"@/"+dbName)
	if err != nil {
		log.Fatal("Initialize Database: ", err)
	}
	Engine.SetMaxIdleConns(5)
}
