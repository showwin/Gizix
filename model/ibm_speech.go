package model

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"

	db "github.com/showwin/Gizix/database"
)

// IBMAccount model
type IBMAccount struct {
	UserName string
	Password string
}

// UpdateIBMAccount : update account info
func UpdateIBMAccount(username string, password string) bool {
	_, err := db.Engine.Exec("UPDATE ibm_acount SET name = ?, password = ? WHERE id = 1", username, password)
	return err == nil
}

// GetIBMAccount : get account info
func GetIBMAccount() (i IBMAccount) {
	db.Engine.QueryRow("SELECT name, password FROM ibm_acount WHERE id = 1 LIMIT 1").Scan(&i.UserName, &i.Password)
	return i
}

// GetIBMToken : issue token for watson speech to text
func GetIBMToken() string {
	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: cookieJar,
	}

	account := GetIBMAccount()
	req, _ := http.NewRequest("GET", "https://stream.watsonplatform.net/authorization/api/v1/token?url=https://stream.watsonplatform.net/speech-to-text/api", nil)
	req.SetBasicAuth(account.UserName, account.Password)
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	return string(body)
}
