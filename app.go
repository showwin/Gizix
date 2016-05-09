package main

import (
	"fmt"
	"net/http"
	"os"

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

	r.GET("/manual", func(c *gin.Context) {
		c.HTML(http.StatusOK, "manual.tmpl", gin.H{})
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

	Candidate PeerCandidateChild `json:"candidate"`
}

// PeerOffer struct
type PeerOffer struct {
	Type string `json:"type"`
	Sdp  string `json:"sdp"`
}

// PeerCandidate struct
type PeerCandidate struct {
	Candidate PeerCandidateChild `json:"candidate"`
}

// PeerCandidateChild struct
type PeerCandidateChild struct {
	Candidate     string `json:"candidate"`
	SdpMLineIndex int    `json:"sdpMLineIndex"`
	SdpMid        string `json:"sdpMid"`
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	addClient(conn)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade.")
		return
	}

	for {
		var signal PeerSignal
		var offer PeerOffer
		var candidate PeerCandidate
		offerFlg := false
		conn.ReadJSON(&signal)
		if signal.Type == "" {
			candidate.Candidate = signal.Candidate
		} else {
			offer.Type = signal.Type
			offer.Sdp = signal.Sdp
			offerFlg = true
		}
		fmt.Println(signal)

		// broadcast する
		for _, c := range pool {
			if c == conn {
				continue
			}
			if offerFlg {
				c.WriteJSON(&offer)
			} else {
				c.WriteJSON(&candidate)
			}
		}
	}
}
