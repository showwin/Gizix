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

	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()

		c.Redirect(http.StatusFound, "/")
	})

	// websocket interface
	r.GET("/ws", func(c *gin.Context) {
		socketHandler(c.Writer, c.Request)
	})

	authorized := r.Group("/")
	authorized.Use(AuthRequired())
	{
		authorized.GET("/dashboard", func(c *gin.Context) {
			session := sessions.Default(c)
			cUser := currentUser(session)
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

			createRoomMessage := session.Flashes("CreateRoom")[0]
			session.Save()
			c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{
				"CurrentUser":       cUser,
				"JoinedRooms":       joinedRooms,
				"OtherRooms":        otherRooms,
				"Domain":            domain,
				"CreateRoomMessage": createRoomMessage,
			})
		})

		authorized.GET("/room/:roomID", func(c *gin.Context) {
			session := sessions.Default(c)
			cUser := currentUser(session)
			domain := getDomain()
			roomID, _ := strconv.Atoi(c.Param("roomID"))
			room := getRoom(roomID)
			skyway := getSkyWayKey()
			joinedFlg := cUser.IsJoin(roomID)
			joinRoomMessage := session.Flashes("JoinRoom")[0]
			session.Save()
			c.HTML(http.StatusOK, "room.tmpl", gin.H{
				"CurrentUser":     cUser,
				"Domain":          domain,
				"Room":            room,
				"SkyWay":          skyway,
				"JoinedFlg":       joinedFlg,
				"JoinRoomMessage": joinRoomMessage,
			})
		})

		authorized.GET("/setting", func(c *gin.Context) {
			session := sessions.Default(c)
			cUser := currentUser(session)

			if cUser.Admin {
				allUser := allUser()
				domain := getDomain()
				skyway := getSkyWayKey()
				createUserMessage := session.Flashes("CreateUser")[0]
				updateDomainMessage := session.Flashes("UpdateDomain")[0]
				updateSkyWayKeyMessage := session.Flashes("UpdateSkyWayKey")[0]
				session.Save()
				c.HTML(http.StatusOK, "setting.tmpl", gin.H{
					"CurrentUser":            cUser,
					"AllUser":                allUser,
					"Domain":                 domain,
					"SkyWay":                 skyway,
					"CreateUserMessage":      createUserMessage,
					"UpdateDomainMessage":    updateDomainMessage,
					"UpdateSkyWayKeyMessage": updateSkyWayKeyMessage,
				})
			} else {
				c.HTML(http.StatusOK, "setting.tmpl", gin.H{
					"CurrentUser": cUser,
				})
			}
		})

		// create room
		authorized.POST("/room", func(c *gin.Context) {
			session := sessions.Default(c)
			roomName := c.PostForm("name")
			if !createRoom(roomName) {
				session.AddFlash("すでにその Room は作成されています。別の名前でお試しください。", "CreateRoom")
			}
			session.Save()
			c.Redirect(http.StatusSeeOther, "/dashboard")
		})

		// join the room
		authorized.POST("/join", func(c *gin.Context) {
			session := sessions.Default(c)
			cUser := currentUser(session)
			roomID, _ := strconv.Atoi(c.PostForm("roomID"))
			if cUser.JoinRoom(roomID) {
				c.Redirect(http.StatusSeeOther, "/dashboard")
			} else {
				session.AddFlash("Room の参加に失敗しました。", "JoinRoom")
				session.Save()
				c.Redirect(http.StatusSeeOther, "/room/"+c.Param("roomID"))
			}
		})

		// Admin Required
		admin := authorized.Group("/")
		admin.Use(AdminRequired())
		{
			// create user
			admin.POST("/user", func(c *gin.Context) {
				session := sessions.Default(c)

				userName := c.PostForm("name")
				if createUser(userName) {
					session.AddFlash("アカウント: "+userName+"を作成しました。パスワードは'password'です。", "CreateUser")
				} else {
					session.AddFlash("すでにそのアカウント名は作成されています。別の名前でお試しください。", "CreateUser")
				}
				session.Save()
				c.Redirect(http.StatusSeeOther, "/setting")
			})

			// update Gizix domain (or ip-address)
			admin.POST("/domain", func(c *gin.Context) {
				session := sessions.Default(c)

				domainName := c.PostForm("name")
				if updateDomain(domainName) {
					session.AddFlash("ドメイン名:"+domainName+" に設定しました。", "UpdateDomain")
				} else {
					session.AddFlash("ドメイン名の設定に失敗しました。", "UpdateDomain")
				}
				session.Save()
				c.Redirect(http.StatusSeeOther, "/setting")
			})

			// update SkyWay API Key
			admin.POST("/skyway", func(c *gin.Context) {
				session := sessions.Default(c)

				skywayKey := c.PostForm("key")
				if updateSkyWayKey(skywayKey) {
					session.AddFlash("SkyWay API Key:"+skywayKey+" に設定しました。", "UpdateSkyWayKey")
				} else {
					session.AddFlash("SkyWay API Key の設定に失敗しました。", "UpdateSkyWayKey")
				}
				session.Save()
				c.Redirect(http.StatusSeeOther, "/setting")
			})
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		r.Run(":5000")
	} else {
		r.Run(":" + port)
	}
}

// AuthRequired : redirect to "/" if not authorized
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := sessions.Default(c).Get("uid")
		if uid == nil {
			c.Redirect(http.StatusFound, "/")
		}
	}
}

// AdminRequired : redirect to "/dashboard" if not admin
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		cUser := currentUser(sessions.Default(c))
		if cUser.Admin == false {
			c.Redirect(http.StatusFound, "/dashboard")
		}
	}
}
