package main

import (
	"fmt"
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

func allUser() (users []User) {
	rows, err := db.Query("SELECT id, name, admin FROM users")
	if err != nil {
		return nil
	}

	defer rows.Close()
	for rows.Next() {
		u := User{}
		err = rows.Scan(&u.ID, &u.Name, &u.Admin)
		if err != nil {
			fmt.Println(err)
		}
		users = append(users, u)
	}
	return
}

func currentUser(session sessions.Session) (u User) {
	uID := session.Get("uid")
	r := db.QueryRow("SELECT id, name, password, admin FROM users WHERE id = ? LIMIT 1", uID)
	r.Scan(&u.ID, &u.Name, &u.Password, &u.Admin)
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

// JoinedRooms : return rooms which user joined
func (u *User) JoinedRooms() (rooms []Room) {
	rows, err := db.Query(
		"SELECT id, name "+
			"FROM rooms "+
			"WHERE id IN ( "+
			"SELECT room_id "+
			"FROM user_room "+
			"WHERE user_id = ?) "+
			"ORDER BY called_at DESC", u.ID)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer rows.Close()
	for rows.Next() {
		r := Room{}
		err = rows.Scan(&r.ID, &r.Name)
		if err != nil {
			fmt.Println(err)
		}
		rooms = append(rooms, r)
	}
	return
}

// NotJoinedRooms : return rooms which user not joined
func (u *User) NotJoinedRooms() (rooms []Room) {
	rows, err := db.Query(
		"SELECT id, name "+
			"FROM rooms "+
			"WHERE id NOT IN ( "+
			"SELECT room_id "+
			"FROM user_room "+
			"WHERE user_id = ?) "+
			"ORDER BY name", u.ID)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer rows.Close()
	for rows.Next() {
		r := Room{}
		err = rows.Scan(&r.ID, &r.Name)
		if err != nil {
			fmt.Println(err)
		}
		rooms = append(rooms, r)
	}
	return
}
