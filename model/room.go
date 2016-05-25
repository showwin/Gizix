package model

import (
	"fmt"
	"time"

	db "github.com/showwin/Gizix/database"
)

// Room model
type Room struct {
	ID        int
	Name      string
	CreatedAt string
	CalledAt  string
}

// RoomUsers model
type RoomUsers struct {
	ID       int
	Name     string
	CalledAt string
	Users    []User
}

// CreateRoom : create room with name
func CreateRoom(name string) bool {
	_, err := db.Engine.Exec(
		"INSERT INTO rooms (name, created_at) VALUES (?, ?)", name, time.Now())
	return err == nil
}

// GetRoom : get room by id
func GetRoom(id int) (r Room) {
	db.Engine.QueryRow("SELECT id, name FROM rooms WHERE id = ?", id).Scan(&r.ID, &r.Name)
	return
}

// WithUsers : return room with users info who belongs.
func (r *Room) WithUsers() (ru RoomUsers) {
	ru.ID = r.ID
	ru.Name = r.Name
	ru.CalledAt = r.CalledAt
	rows, err := db.Engine.Query(
		"SELECT id, name FROM users WHERE id IN (SELECT user_id FROM user_room WHERE room_id = ?)", r.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		u := User{}
		err = rows.Scan(&u.ID, &u.Name)
		if err != nil {
			fmt.Println(err)
		}
		ru.Users = append(ru.Users, u)
	}
	return
}
