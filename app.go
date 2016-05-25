package main

import (
	"database/sql"
	"net/http"
	"os"

	c "github.com/showwin/Gizix/controller"
	database "github.com/showwin/Gizix/database"
	m "github.com/showwin/Gizix/model"

	"github.com/gin-gonic/contrib/commonlog"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	// initialize database in production mode
	database.Initialize(false)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Use(static.Serve("/js", static.LocalFile("public/js", true)))
	r.Use(static.Serve("/css", static.LocalFile("public/css", true)))
	r.Use(static.Serve("/img", static.LocalFile("public/img", true)))
	r.Use(static.Serve("/fonts", static.LocalFile("public/fonts", true)))

	// use session store
	store := sessions.NewCookieStore([]byte("gizix_happy"))
	store.Options(sessions.Options{HttpOnly: true})
	r.Use(sessions.Sessions("mysession", store))

	// use log middleware
	r.Use(commonlog.New())

	// top page
	r.GET("/", c.GetIndex)

	// login
	r.POST("/login", c.PostLogin)

	// logout
	r.GET("/logout", c.GetLogout)

	// websocket interface
	r.GET("/ws", c.SocketHandler)

	authorized := r.Group("/")
	authorized.Use(AuthRequired())
	{
		// dashboard page
		authorized.GET("/dashboard", c.GetDashboard)

		// room page
		authorized.GET("/room/:roomID", c.GetRoom)

		// setting page
		authorized.GET("/setting", c.GetSetting)

		// change password
		authorized.POST("/password", c.PostPassword)

		// create room
		authorized.POST("/room", c.PostRoom)

		// join the room
		authorized.POST("/join", c.PostJoin)

		// Admin Required
		admin := authorized.Group("/")
		admin.Use(AdminRequired())
		{
			// create user
			admin.POST("/user", c.PostUser)

			// update Gizix domain (or ip-address)
			admin.POST("/domain", c.PostDomain)

			// update SkyWay API Key
			admin.POST("/skyway", c.PostSkyWay)
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
		cUser := m.CurrentUser(sessions.Default(c))
		if cUser.Admin == false {
			c.Redirect(http.StatusFound, "/dashboard")
		}
	}
}
