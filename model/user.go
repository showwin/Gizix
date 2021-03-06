package model

import (
	"fmt"
	"time"

	db "github.com/showwin/Gizix/database"

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

// BcryptCost : cost value of encrypting by bcrypt
const BcryptCost = 10

// Authenticate : user authentication
func Authenticate(name string, password string) (user User, result bool) {
	dbErr := db.Engine.QueryRow("SELECT id, name, password, admin FROM users WHERE name = ? LIMIT 1", name).
		Scan(&user.ID, &user.Name, &user.Password, &user.Admin)
	authErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return user, (dbErr == nil && authErr == nil)
}

// AllUser : get all users
func AllUser() (users []User) {
	rows, err := db.Engine.Query("SELECT id, name, admin FROM users")
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

// CurrentUser : get current user
func CurrentUser(uID int) (u User) {
	r := db.Engine.QueryRow("SELECT id, name, password, admin FROM users WHERE id = ? LIMIT 1", uID)
	r.Scan(&u.ID, &u.Name, &u.Password, &u.Admin)
	return u
}

// CreateUser : create user with the name
func CreateUser(name string) bool {
	password := []byte("password")
	hashedPass, _ := bcrypt.GenerateFromPassword(password, BcryptCost)

	_, err := db.Engine.Exec(
		"INSERT INTO users (name, password, admin, created_at) VALUES (?, ?, ?, ?)",
		name, hashedPass, false, time.Now())
	return err == nil
}

// UpdatePassword : return success to update password or not
func (u *User) UpdatePassword(oldPassword string, newPassword string) bool {
	authErr := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oldPassword))
	if authErr != nil {
		return false
	}

	// update password
	password := []byte(newPassword)
	hashedPass, _ := bcrypt.GenerateFromPassword(password, BcryptCost)
	_, err := db.Engine.Exec("UPDATE users SET password = ? WHERE id = ?", hashedPass, u.ID)
	return err == nil
}

// UpdateAdmin : return success to update admin or not
func (u *User) UpdateAdmin(admin bool) bool {
	// update admin
	_, err := db.Engine.Exec("UPDATE users SET admin = ? WHERE id = ?", admin, u.ID)
	return err == nil
}

// JoinedRooms : return rooms which user joined
func (u *User) JoinedRooms() (rooms []Room) {
	rows, err := db.Engine.Query(
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
	rows, err := db.Engine.Query(
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

// IsJoin : return user joined the room or not
func (u *User) IsJoin(roomID int) bool {
	var count int
	db.Engine.QueryRow("SELECT count(*) FROM user_room WHERE user_id = ? AND room_id = ?", u.ID, roomID).Scan(&count)
	return count != 0
}

// JoinRoom : user join the room
func (u *User) JoinRoom(roomID int) bool {
	_, err := db.Engine.Exec(
		"INSERT INTO user_room (user_id, room_id) VALUES (?, ?)", u.ID, roomID)
	return err == nil
}
