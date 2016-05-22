package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

// SocketBody struct
type SocketBody struct {
	Type   string `json:"type"`
	Sdp    string `json:"sdp"`
	To     string `json:"to"`
	UID    string `json:"uid"`
	RoomID string `json:"roomID"`

	Candidate PeerCandidateChild `json:"candidate"`
}

// ClientInfo struct
type ClientInfo struct {
	Type string `json:"type"`
	UID  string `json:"uid"`
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 接続したクライアントの一覧
var pool = map[string]*websocket.Conn{}

func registerClient(uid string, c *websocket.Conn) {
	pool[uid] = c
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade.")
		fmt.Println(err)
		return
	}

	for {
		var body SocketBody
		var offer PeerOffer
		var candidate PeerCandidate
		var peerPool PeerPool

		conn.ReadJSON(&body)
		switch body.Type {
		case "register":
			fmt.Println("Type register")
			registerClient(body.UID, conn)

			// send to yourself
			info := ClientInfo{Type: "info", UID: body.UID}
			conn.WriteJSON(info)
			fmt.Println("Register client of: " + body.UID)
		case "initialize":
			// Type: initialize
			fmt.Println("Type initialize")
			peerPool.Type = "initialize"
			// return all user id in the room except you
			roomID, err := strconv.Atoi(body.RoomID)
			if err != nil {
				fmt.Println(err)
			}
			room := getRoom(roomID)
			ru := room.WithUsers()
			for _, u := range ru.Users {
				id := strconv.Itoa(u.ID)
				if id == body.UID {
					fmt.Println("=== EXCLUDE" + body.UID)
					continue
				}
				peerPool.Ids = append(peerPool.Ids, id)
			}
			fmt.Println("Users in Room: ")
			fmt.Println(peerPool.Ids)

			//send to yourself
			conn.WriteJSON(peerPool)
		case "offer", "answer":
			// Type: offer
			fmt.Println("Type offer")
			offer.Type = body.Type
			offer.Sdp = body.Sdp
			offer.From = body.UID

			// send to other
			c := pool[body.To]
			c.WriteJSON(&offer)
			fmt.Printf("Send to " + body.To + "\n")
		default:
			// Type: candidate
			fmt.Println("Type Candidate")
			candidate.Candidate = body.Candidate
			candidate.From = body.UID

			// send to other
			c := pool[body.To]
			c.WriteJSON(&candidate)
			fmt.Printf("Send to " + body.To + "\n")
		}
	}
}
