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
		var updateAdminsMessage interface{}
		if f := session.Flashes("UpdateAdmins"); len(f) != 0 {
			updateAdminsMessage = f[0]
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
			"UpdateAdminsMessage":     updateAdminsMessage,
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
		session.AddFlash("Please make sure passwords are identical.", "UpdatePassword")
	} else if cUser.UpdatePassword(oldPass, newPass) {
		session.AddFlash("Update password successfully.", "UpdatePassword")
	} else {
		session.AddFlash("Failed to update password.", "UpdatePassword")
	}
	session.Save()
	c.Redirect(http.StatusSeeOther, "/setting")
}

// PostUser response from POST /user
func PostUser(c *gin.Context) {
	session := sessions.Default(c)

	userName := c.PostForm("name")
	if m.CreateUser(userName) {
		session.AddFlash("Create a new account: '"+userName+"'. Default password: 'password'.", "CreateUser")
	} else {
		session.AddFlash("The user already exists.", "CreateUser")
	}
	session.Save()
	c.Redirect(http.StatusSeeOther, "/setting")
}

// PostAdmins response from POST /admins
func PostAdmins(c *gin.Context) {
	session := sessions.Default(c)

	succFlg := true
	for _, user := range m.AllUser() {
		// Gizix always must be Admin
		admin := c.PostForm(user.Name) == "true" || user.Name == "Gizix"
		succFlg = succFlg && user.UpdateAdmin(admin)
	}

	if succFlg {
		session.AddFlash("Update authorities succssfully.", "UpdateAdmins")
	} else {
		session.AddFlash("Failed to update authorities.", "UpdateAdmins")
	}
	session.Save()
	c.Redirect(http.StatusSeeOther, "/setting")
}

// PostDomain response from POST /domain
func PostDomain(c *gin.Context) {
	session := sessions.Default(c)

	domainName := c.PostForm("name")
	if m.UpdateDomain(domainName) {
		session.AddFlash("Update domain name succssfully.", "UpdateDomain")
	} else {
		session.AddFlash("Failed to update domain name.", "UpdateDomain")
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
		session.AddFlash("Update speech language succssfully.", "UpdateLanguage")
	} else {
		session.AddFlash("Failed to update speech language", "UpdateLanguage")
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
		session.AddFlash("Update IBM Account succssfully.", "UpdateIBMAccount")
	} else {
		session.AddFlash("Failed to update IBM Account.", "UpdateIBMAccount")
	}
	session.Save()
	c.Redirect(http.StatusSeeOther, "/setting")
}
