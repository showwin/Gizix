package model

import (
	"testing"
	"time"

	db "github.com/showwin/Gizix/database"
)

func TestCreateRoom(t *testing.T) {
	db.Initialize(true)

	roomName := "testRoom" + time.Now().String()
	CreateRoom(roomName)

	var r Room
	db.Engine.QueryRow("SELECT name FROM rooms ORDER BY id DESC LIMIT 1").Scan(&r.Name)
	if r.Name != roomName {
		t.Errorf("got: %v\nexpected: %v", r.Name, roomName)
	}
}

func TestGetRoom(t *testing.T) {
	db.Initialize(true)

	roomName := "testRoom" + time.Now().String()
	CreateRoom(roomName)

	var rID int
	db.Engine.QueryRow("SELECT id FROM rooms ORDER BY id DESC LIMIT 1").Scan(&rID)
	r := GetRoom(rID)

	if r.Name != roomName {
		t.Errorf("got: %v\nexpected: %v", r.Name, roomName)
	}
}

func TestWithUser(t *testing.T) {
	db.Initialize(true)
	roomName := "testRoom3"
	CreateRoom(roomName)
	var rID int
	db.Engine.QueryRow("SELECT id FROM rooms ORDER BY id DESC LIMIT 1").Scan(&rID)
	r := GetRoom(rID)

	db.Engine.Exec("INSERT INTO user_room (user_id, room_id) VALUES (?, ?), (?, ?)", 1, r.ID, 2, r.ID)
	ru := r.WithUsers()

	if ru.ID != r.ID {
		t.Errorf("got: %v\nexpected: %v", ru.ID, r.ID)
	}
	if ru.Name != r.Name {
		t.Errorf("got: %v\nexpected: %v", ru.Name, r.Name)
	}
	if ru.Users[0].Name != "Gizix" {
		t.Errorf("got: %v\nexpected: %v", ru.Users[0].Name, "Gizix")
	}
	if ru.Users[1].Name != "Gizix2" {
		t.Errorf("got: %v\nexpected: %v", ru.Users[1].Name, "Gizix2")
	}
}
