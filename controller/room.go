package controller

import (
	"net/http"
	"strconv"

	m "github.com/showwin/Gizix/model"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// GetRoom response from GET /room/:id
func GetRoom(c *gin.Context) {
	session := sessions.Default(c)
	cUser := m.CurrentUser(session.Get("uid").(int))
	domain := m.GetDomain()
	roomID, _ := strconv.Atoi(c.Param("roomID"))
	room := m.GetRoom(roomID)
	wToken := m.GetWatsonToken()
	joinedFlg := cUser.IsJoin(roomID)

	// Flash Message
	var joinRoomMessage interface{}
	if f := session.Flashes("JoinRoom"); len(f) != 0 {
		joinRoomMessage = f[0]
	}
	session.Save()
	c.HTML(http.StatusOK, "room.tmpl", gin.H{
		"CurrentUser":     cUser,
		"Domain":          domain,
		"Room":            room,
		"WatsonToken":     wToken,
		"JoinedFlg":       joinedFlg,
		"JoinRoomMessage": joinRoomMessage,
	})
}

// PostRoom response from POST /room
func PostRoom(c *gin.Context) {
	session := sessions.Default(c)
	roomName := c.PostForm("name")
	if !m.CreateRoom(roomName) {
		session.AddFlash("The room already exists. Please try with a different name.", "CreateRoom")
	}
	session.Save()
	c.Redirect(http.StatusSeeOther, "/dashboard")
}

// PostJoin response from POST /join
func PostJoin(c *gin.Context) {
	session := sessions.Default(c)
	cUser := m.CurrentUser(session.Get("uid").(int))
	roomID, _ := strconv.Atoi(c.PostForm("roomID"))
	if cUser.JoinRoom(roomID) {
		c.Redirect(http.StatusSeeOther, "/dashboard")
	} else {
		session.AddFlash("Sorry, failed to join this room.", "JoinRoom")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/room/"+c.Param("roomID"))
	}
}
