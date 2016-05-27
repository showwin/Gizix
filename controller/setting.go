package controller

import (
	"net/http"

	m "github.com/showwin/Gizix/model"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// GetSetting response from GET /setting
func GetSetting(c *gin.Context) {
	session := sessions.Default(c)
	cUser := m.CurrentUser(session.Get("uid").(int))

	if cUser.Admin {
		allUser := m.AllUser()
		domain := m.GetDomain()
		skyway := m.GetSkyWayKey()

		// Flash Message
		var updatePasswordMessage interface{}
		if f := session.Flashes("UpdatePassword"); len(f) != 0 {
			updatePasswordMessage = f[0]
		}
		var createUserMessage interface{}
		if f := session.Flashes("CreateUser"); len(f) != 0 {
			createUserMessage = f[0]
		}
		var updateDomainMessage interface{}
		if f := session.Flashes("UpdateDomain"); len(f) != 0 {
			updateDomainMessage = f[0]
		}
		var updateSkyWayKeyMessage interface{}
		if f := session.Flashes("UpdateSkyWayKey"); len(f) != 0 {
			updateSkyWayKeyMessage = f[0]
		}
		session.Save()
		c.HTML(http.StatusOK, "setting.tmpl", gin.H{
			"CurrentUser":            cUser,
			"AllUser":                allUser,
			"Domain":                 domain,
			"SkyWay":                 skyway,
			"UpdatePasswordMessage":  updatePasswordMessage,
			"CreateUserMessage":      createUserMessage,
			"UpdateDomainMessage":    updateDomainMessage,
			"UpdateSkyWayKeyMessage": updateSkyWayKeyMessage,
		})
	} else {
		c.HTML(http.StatusOK, "setting.tmpl", gin.H{
			"CurrentUser": cUser,
		})
	}
}

// PostPassword response from POST /password
func PostPassword(c *gin.Context) {
	session := sessions.Default(c)
	cUser := m.CurrentUser(session.Get("uid").(int))
	oldPass := c.PostForm("old_password")
	newPass := c.PostForm("new_password")
	confPass := c.PostForm("confirm_password")
	if newPass != confPass {
		session.AddFlash("新しいパスワードが一致しません。", "UpdatePassword")
	} else if cUser.UpdatePassword(oldPass, newPass) {
		session.AddFlash("パスワードを更新しました。", "UpdatePassword")
	} else {
		session.AddFlash("パスワードの更新に失敗しました。", "UpdatePassword")
	}
	session.Save()
	c.Redirect(http.StatusSeeOther, "/setting")
}

// PostUser response from POST /user
func PostUser(c *gin.Context) {
	session := sessions.Default(c)

	userName := c.PostForm("name")
	if m.CreateUser(userName) {
		session.AddFlash("アカウント: "+userName+"を作成しました。パスワードは'password'です。", "CreateUser")
	} else {
		session.AddFlash("すでにそのアカウント名は作成されています。別の名前でお試しください。", "CreateUser")
	}
	session.Save()
	c.Redirect(http.StatusSeeOther, "/setting")
}

// PostDomain response from POST /domain
func PostDomain(c *gin.Context) {
	session := sessions.Default(c)

	domainName := c.PostForm("name")
	if m.UpdateDomain(domainName) {
		session.AddFlash("ドメイン名:"+domainName+" に設定しました。", "UpdateDomain")
	} else {
		session.AddFlash("ドメイン名の設定に失敗しました。", "UpdateDomain")
	}
	session.Save()
	c.Redirect(http.StatusSeeOther, "/setting")
}

// PostSkyWay response from POST /skyway
func PostSkyWay(c *gin.Context) {
	session := sessions.Default(c)

	skyWayKey := c.PostForm("key")
	if m.UpdateSkyWayKey(skyWayKey) {
		session.AddFlash("SkyWay API Key:"+skyWayKey+" に設定しました。", "UpdateSkyWayKey")
	} else {
		session.AddFlash("SkyWay API Key の設定に失敗しました。", "UpdateSkyWayKey")
	}
	session.Save()
	c.Redirect(http.StatusSeeOther, "/setting")
}
