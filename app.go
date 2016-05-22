package main

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/contrib/commonlog"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
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
		var joinedRooms = []RoomUsers{}
		var otherRooms = []RoomUsers{}
		jr := cUser.JoinedRooms()
		or := cUser.NotJoinedRooms()
		for _, r := range jr {
			ru := r.WithUsers()
			joinedRooms = append(joinedRooms, ru)
		}
		for _, r := range or {
			ru := r.WithUsers()
			otherRooms = append(otherRooms, ru)
		}
		domain := getDomain()
		c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{
			"CurrentUser": cUser,
			"JoinedRooms": joinedRooms,
			"OtherRooms":  otherRooms,
			"Domain":      domain,
		})
	})

	r.GET("/room/:roomID", func(c *gin.Context) {
		cUser := currentUser(sessions.Default(c))
		domain := getDomain()
		roomID, _ := strconv.Atoi(c.Param("roomID"))
		room := getRoom(roomID)
		c.HTML(http.StatusOK, "room.tmpl", gin.H{
			"CurrentUser": cUser,
			"Domain":      domain,
			"Room":        room,
		})
	})

	r.GET("/setting", func(c *gin.Context) {
		cUser := currentUser(sessions.Default(c))

		if cUser.Admin {
			allUser := allUser()
			domain := getDomain()
			c.HTML(http.StatusOK, "setting.tmpl", gin.H{
				"CurrentUser": cUser,
				"AllUser":     allUser,
				"Domain":      domain,
			})
		} else {
			c.HTML(http.StatusOK, "setting.tmpl", gin.H{
				"CurrentUser": cUser,
			})
		}
	})

	// create user
	r.POST("/user", func(c *gin.Context) {
		cUser := currentUser(sessions.Default(c))
		allUser := allUser()
		domain := getDomain()

		userName := c.PostForm("name")
		if createUser(userName) {
			c.HTML(http.StatusOK, "setting.tmpl", gin.H{
				"CurrentUser":       cUser,
				"AllUser":           allUser,
				"Domain":            domain,
				"CreateUserMessage": "アカウント: " + userName + "を作成しました。パスワードは'password'です。",
			})
		} else {
			c.HTML(http.StatusOK, "setting.tmpl", gin.H{
				"CurrentUser":       cUser,
				"AllUser":           allUser,
				"CreateUserMessage": "すでにそのアカウント名は作成されています。別の名前でお試しください。",
			})
		}
	})

	// create room
	r.POST("/room", func(c *gin.Context) {
		roomName := c.PostForm("name")
		if createRoom(roomName) {
			c.Redirect(http.StatusSeeOther, "/dashboard")
		} else {
			c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{
				"CreateRoomMessage": "すでにその Room は作成されています。別の名前でお試しください。",
			})
		}
	})

	// update Gizix domain (or ip-address)
	r.POST("/domain", func(c *gin.Context) {
		domainName := c.PostForm("name")

		cUser := currentUser(sessions.Default(c))
		allUser := allUser()
		if updateDomain(domainName) {
			domain := getDomain()
			c.HTML(http.StatusOK, "setting.tmpl", gin.H{
				"CurrentUser":         cUser,
				"AllUser":             allUser,
				"Domain":              domain,
				"UpdateDomainMessage": "ドメイン名:" + domainName + " に設定しました。",
			})
		} else {
			domain := getDomain()
			c.HTML(http.StatusOK, "setting.tmpl", gin.H{
				"CurrentUser":         cUser,
				"AllUser":             allUser,
				"Domain":              domain,
				"UpdateDomainMessage": "設定に失敗しました。",
			})
		}
	})

	r.GET("/test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	// websocket interface
	r.GET("/ws", func(c *gin.Context) {
		socketHandler(c.Writer, c.Request)
	})

	port := os.Getenv("PORT")
	if port == "" {
		r.Run(":5000")
	} else {
		r.Run(":" + port)
	}
}
