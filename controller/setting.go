package controller

import (
	"fmt"
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
		lang := m.GetLanguage()
		ibm := m.GetIBMAccount()

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
		var updateLanguageMessage interface{}
		if f := session.Flashes("UpdateLanguage"); len(f) != 0 {
			updateLanguageMessage = f[0]
		}
		var updateIBMAccountMessage interface{}
		if f := session.Flashes("UpdateIBMAccount"); len(f) != 0 {
			updateIBMAccountMessage = f[0]
		}
		session.Save()
		c.HTML(http.StatusOK, "setting.tmpl", gin.H{
			"CurrentUser":             cUser,
			"AllUser":                 allUser,
			"Domain":                  domain,
			"Language":                lang,
			"IBMAccount":              ibm,
			"UpdatePasswordMessage":   updatePasswordMessage,
			"CreateUserMessage":       createUserMessage,
			"UpdateDomainMessage":     updateDomainMessage,
			"UpdateLanguageMessage":   updateLanguageMessage,
			"UpdateIBMAccountMessage": updateIBMAccountMessage,
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

// PostLangugage response from POST /language
func PostLanguage(c *gin.Context) {
	session := sessions.Default(c)

	language := c.PostForm("language")
	fmt.Println(language)
	if m.UpdateLanguage(language) {
		session.AddFlash("言語の設定を更新しました。", "UpdateLanguage")
	} else {
		session.AddFlash("言語の設定に失敗しました。", "UpdateLanguage")
	}
	session.Save()
	c.Redirect(http.StatusSeeOther, "/setting")
}

// PostIBMAccount response from POST /ibm_account
func PostIBMAccount(c *gin.Context) {
	session := sessions.Default(c)

	userName := c.PostForm("username")
	password := c.PostForm("password")
	if m.UpdateIBMAccount(userName, password) {
		session.AddFlash("IBM Account の設定を更新しました。", "UpdateIBMAccount")
	} else {
		session.AddFlash("IBM Account の設定に失敗しました。", "UpdateIBMAccount")
	}
	session.Save()
	c.Redirect(http.StatusSeeOther, "/setting")
}
