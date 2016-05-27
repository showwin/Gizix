package model

import (
	"testing"
	"time"

	db "github.com/showwin/Gizix/database"
)

func TestAuthenticate(t *testing.T) {
	u, result := Authenticate("Gizix", "password")
	if result != true {
		t.Error("User Authentication FAILED")
	}

	if u.Name != "Gizix" {
		t.Errorf("got: %v\nexpected: %v", u.Name, "Gizix")
	}
}

func TestAllUser(t *testing.T) {
	db.Initialize(true)

	var allBefore int
	db.Engine.QueryRow("SELECT count(*) FROM users").Scan(&allBefore)

	name := "user" + time.Now().String()
	CreateUser(name)

	allUser := AllUser()
	if len(allUser) != allBefore+1 {
		t.Errorf("got: %v\nexpected: %v", len(allUser), allBefore+1)
	}

	if allUser[len(allUser)-1].Name != name {
		t.Errorf("got: %v\nexpected: %v", allUser[len(allUser)-1].Name, name)
	}
}

func TestCurrentUser(t *testing.T) {
	uID := 1
	u := CurrentUser(uID)

	if u.ID != uID {
		t.Errorf("got: %v\nexpected: %v", u.ID, uID)
	}
	if u.Name != "Gizix" {
		t.Errorf("got: %v\nexpected: %v", u.Name, "Gizix")
	}
}

func TestCreateUser(t *testing.T) {
	db.Initialize(true)
	name := "user" + time.Now().String()
	CreateUser(name)

	_, result := Authenticate(name, "password")
	if result != true {
		t.Errorf("user '%v' is not created with password: 'password'", name)
	}

	var admin bool
	db.Engine.QueryRow("SELECT admin FROM users WHERE name = ?", name).Scan(&admin)
	if admin != false {
		t.Errorf("user '%v' is not created with admin: %v", name, false)
	}
}

func TestUpdatePassword(t *testing.T) {
	name := "user" + time.Now().String()
	CreateUser(name)

	// before change password
	_, result := Authenticate(name, "password")
	if result != true {
		t.Errorf("user '%v' is not created with password: 'password'", name)
	}

	var u User
	db.Engine.QueryRow("SELECT id, password FROM users WHERE name = ? LIMIT 1", name).
		Scan(&u.ID, &u.Password)
	newPassword := "newPassword"
	err := u.UpdatePassword("password", newPassword)
	if err != true {
		t.Errorf("cannot update user password")
	}

	// after change password
	_, result = Authenticate(name, newPassword)
	if result != true {
		t.Errorf("user password is not changed to: %v", newPassword)
	}
}

func TestJoinedRooms(t *testing.T) {
	db.Initialize(true)
	uName := "user" + time.Now().String()
	CreateUser(uName)
	rName := "room" + time.Now().String()
	CreateRoom(rName)

	var u User
	var rID int
	db.Engine.QueryRow("SELECT id FROM users WHERE name = ?", uName).Scan(&u.ID)
	db.Engine.QueryRow("SELECT id FROM rooms WHERE name = ?", rName).Scan(&rID)
	u.JoinRoom(rID)

	rooms := u.JoinedRooms()
	if rooms[0].ID != rID {
		t.Errorf("got: %v\nexpected: %v", rooms[0].ID, rID)
	}
	if rooms[0].Name != rName {
		t.Errorf("got: %v\nexpected: %v", rooms[0].Name, rName)
	}
}

func TestNotJoinedRooms(t *testing.T) {
	db.Initialize(true)
	uName := "user" + time.Now().String()
	CreateUser(uName)
	var notJoinedRoomCount int
	db.Engine.QueryRow("SELECT count(*) FROM rooms").Scan(&notJoinedRoomCount)

	rName := "room" + time.Now().String()
	CreateRoom(rName)

	var u User
	var rID int
	db.Engine.QueryRow("SELECT id FROM users WHERE name = ?", uName).Scan(&u.ID)
	db.Engine.QueryRow("SELECT id FROM rooms WHERE name = ?", rName).Scan(&rID)
	u.JoinRoom(rID)

	rooms := u.NotJoinedRooms()
	if len(rooms) != notJoinedRoomCount {
		t.Errorf("got: %v\nexpected: %v", len(rooms), notJoinedRoomCount)
	}
}

// testing also func JoinRoom
func TestIsJoin(t *testing.T) {
	db.Initialize(true)
	uName := "user" + time.Now().String()
	CreateUser(uName)
	rName := "room" + time.Now().String()
	CreateRoom(rName)

	var u User
	var rID int
	db.Engine.QueryRow("SELECT id FROM users WHERE name = ?", uName).Scan(&u.ID)
	db.Engine.QueryRow("SELECT id FROM rooms WHERE name = ?", rName).Scan(&rID)
	u.JoinRoom(rID)

	if u.IsJoin(rID) != true {
		t.Errorf("got: %v\nexpected: %v", u.IsJoin(rID), true)
	}
}
