package controller

import (
	"net/http"

	m "github.com/showwin/Gizix/model"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// GetIndex response from GET /
func GetIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", gin.H{})
}

// PostLogin response from POST /login
func PostLogin(c *gin.Context) {
	userName := c.PostForm("name")
	password := c.PostForm("password")

	session := sessions.Default(c)
	user, result := m.Authenticate(userName, password)
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
}

// GetLogout response from GET /Logout
func GetLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusFound, "/")
}

// GetDashboard response from Get /dashboard
func GetDashboard(c *gin.Context) {
	session := sessions.Default(c)
	cUser := m.CurrentUser(session)
	var joinedRooms = []m.RoomUsers{}
	var otherRooms = []m.RoomUsers{}
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
	domain := m.GetDomain()

	// Flash Message
	var createRoomMessage interface{}
	if f := session.Flashes("CreateRoom"); len(f) != 0 {
		createRoomMessage = f[0]
	}
	session.Save()
	c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{
		"CurrentUser":       cUser,
		"JoinedRooms":       joinedRooms,
		"OtherRooms":        otherRooms,
		"Domain":            domain,
		"CreateRoomMessage": createRoomMessage,
	})
}
