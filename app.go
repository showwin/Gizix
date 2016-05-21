package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/contrib/commonlog"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
)

var db *sql.DB

func main() {
	// database setting
	user := "root"
	dbname := "gizix"
	db, _ = sql.Open("mysql", user+"@/"+dbname)
	db.SetMaxIdleConns(5)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Use(static.Serve("/js", static.LocalFile("public/js", true)))
	r.Use(static.Serve("/css", static.LocalFile("public/css", true)))
	r.Use(static.Serve("/img", static.LocalFile("public/img", true)))
	r.Use(static.Serve("/fonts", static.LocalFile("public/fonts", true)))

	// session store
	store := sessions.NewCookieStore([]byte("gizix_happy"))
	store.Options(sessions.Options{HttpOnly: true})
	r.Use(sessions.Sessions("mysession", store))

	r.Use(commonlog.New())

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{})
	})

	r.POST("/login", func(c *gin.Context) {
		userName := c.PostForm("name")
		password := c.PostForm("password")

		session := sessions.Default(c)
		user, result := authenticate(userName, password)
		if result {
			//認証成功
			session.Set("uid", user.ID)
			session.Save()

			c.Redirect(http.StatusSeeOther, "/dashboard")
		} else {
			//認証失敗
			c.HTML(http.StatusOK, "login.tmpl", gin.H{
				"Message": "アカウント名かパスワードが間違っています。",
			})
		}
	})

	r.GET("/dashboard", func(c *gin.Context) {
		cUser := currentUser(sessions.Default(c))
		// adminかどうかで場合分け
		c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{
			"CurrentUser": cUser,
		})
	})

	r.GET("/room", func(c *gin.Context) {
		c.HTML(http.StatusOK, "room.tmpl", gin.H{})
	})

	r.GET("/setting", func(c *gin.Context) {
		cUser := currentUser(sessions.Default(c))

		if cUser.Admin {
			allUser := allUser()
			c.HTML(http.StatusOK, "setting.tmpl", gin.H{
				"CurrentUser": cUser,
				"AllUser":     allUser,
			})
		} else {
			c.HTML(http.StatusOK, "setting.tmpl", gin.H{
				"CurrentUser": cUser,
			})
		}
	})

	r.POST("/user", func(c *gin.Context) {
		userName := c.PostForm("name")
		if createUser(userName) {
			c.HTML(http.StatusOK, "setting.tmpl", gin.H{
				"CreateUserMessage": "アカウント: " + userName + "を作成しました。パスワードは'password'です。",
			})
		} else {
			c.HTML(http.StatusOK, "setting.tmpl", gin.H{
				"CreateUserMessage": "すでにそのアカウント名は作成されています。別の名前でお試しください。",
			})
		}
	})

	r.GET("/test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "test.tmpl", gin.H{})
	})

	// websocket interface
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

// PeerOffer struct: use when offering to start P2P
type PeerOffer struct {
	Type string `json:"type"`
	Sdp  string `json:"sdp"`
	From string `json:"from"`
}

// PeerCandidate struct: use after offering
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

// PeerPool struct: connected client ids
type PeerPool struct {
	Type string   `json:"type"`
	Ids  []string `json:"ids"`
}

// PeerConfig struct: client config
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
	// response your client id
	config := PeerConfig{Type: "config", ID: cID}
	conn.WriteJSON(&config)

	for {
		var signal PeerSignal
		var offer PeerOffer
		var candidate PeerCandidate
		var peerPool PeerPool
		offerFlg := false

		conn.ReadJSON(&signal)
		if signal.Type == "" {
			// Type: candidate
			candidate.Candidate = signal.Candidate
			candidate.From = cID
		} else if signal.Type == "initialize" {
			// Type: initialize
			peerPool.Type = "initialize"
			for i := range pool {
				id := strconv.Itoa(i)
				// exclude yourself
				if id == cID {
					continue
				}
				peerPool.Ids = append(peerPool.Ids, id)
			}
		} else {
			// Type: offer
			offer.Type = signal.Type
			offer.Sdp = signal.Sdp
			offer.From = cID
			offerFlg = true
		}
		fmt.Println("Connection Pool:")
		fmt.Println(pool)

		if signal.To == "myself" {
			conn.WriteJSON(peerPool)
			fmt.Println("Send Pool Info")
			fmt.Println(peerPool)
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
