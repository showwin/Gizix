package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/contrib/commonlog"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Use(static.Serve("/js", static.LocalFile("public/js", true)))

	r.Use(commonlog.New())

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	r.GET("/test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "test.tmpl", gin.H{})
	})

	r.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})

	port := os.Getenv("PORT")
	if port == "" {
		r.Run(":5000")
	} else {
		r.Run(":" + port)
	}
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 接続したクライアントの一覧
var pool = []*websocket.Conn{}

func addClient(c *websocket.Conn) {
	pool = append(pool, c)
}

// PeerSignal struct
type PeerSignal struct {
	Type string `json:"type"`
	Sdp  string `json:"sdp"`
	To   string `json:"to"`

	Candidate PeerCandidateChild `json:"candidate"`
}

// PeerOffer struct
type PeerOffer struct {
	Type string `json:"type"`
	Sdp  string `json:"sdp"`
	From string `json:"from"`
}

// PeerCandidate struct
type PeerCandidate struct {
	Candidate PeerCandidateChild `json:"candidate"`
	From      string             `json:"from"`
}

// PeerCandidateChild struct
type PeerCandidateChild struct {
	Candidate     string `json:"candidate"`
	SdpMLineIndex int    `json:"sdpMLineIndex"`
	SdpMid        string `json:"sdpMid"`
}

// PeerConfig struct
type PeerConfig struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	cID := strconv.Itoa(len(pool))
	addClient(conn)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade.")
		fmt.Println(err)
		return
	}
	config := PeerConfig{Type: "config", ID: cID}
	conn.WriteJSON(&config)

	for {
		var signal PeerSignal
		var offer PeerOffer
		var candidate PeerCandidate
		offerFlg := false
		conn.ReadJSON(&signal)
		if signal.Type == "" {
			candidate.Candidate = signal.Candidate
			candidate.From = cID
		} else {
			offer.Type = signal.Type
			offer.Sdp = signal.Sdp
			offer.From = cID
			offerFlg = true
		}

		if signal.To == "broadcast" {
			for _, c := range pool {
				// Don't send it yourself
				if c == conn {
					continue
				}

				if offerFlg {
					c.WriteJSON(&offer)
				} else {
					c.WriteJSON(&candidate)
				}
			}
			fmt.Println("Send to Broadcast")
		} else {
			to, _ := strconv.Atoi(signal.To)
			c := pool[to]
			if offerFlg {
				c.WriteJSON(&offer)
			} else {
				c.WriteJSON(&candidate)
			}
			fmt.Printf("Send to %d\n", to)
		}
	}
}
