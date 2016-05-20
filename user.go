package main

import (
	"time"

	"github.com/gin-gonic/contrib/sessions"

	"golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	ID        int
	Name      string
	Password  string
	IconPath  string
	Admin     bool
	CreatedAt string
	LoginedAt string
}

func authenticate(name string, password string) (user User, result bool) {
	dbErr := db.QueryRow("SELECT id, name, password, admin FROM users WHERE name = ? LIMIT 1", name).
		Scan(&user.ID, &user.Name, &user.Password, &user.Admin)
	authErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return user, (dbErr == nil && authErr == nil)
}

func currentUser(session sessions.Session) (u User) {
	uID := session.Get("uid")
	r := db.QueryRow("SELECT * FROM users WHERE id = ? LIMIT 1", uID)
	r.Scan(&u.ID, &u.Name, &u.Password, &u.IconPath, &u.Admin, &u.CreatedAt, &u.LoginedAt)
	return u
}

func createUser(name string) bool {
	password := []byte("password")
	cost := 10
	hashedPass, _ := bcrypt.GenerateFromPassword(password, cost)

	_, err := db.Exec(
		"INSERT INTO users (name, password, admin, created_at) VALUES (?, ?, ?, ?)",
		name, hashedPass, false, time.Now())
	return err == nil
}
