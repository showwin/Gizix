package model

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"

	db "github.com/showwin/Gizix/database"
)

// WatsonAccount model
type WatsonAccount struct {
	UserName string
	Password string
	Model    string
}

// UpdateWatsonAccount : update account info
func UpdateWatsonAccount(username string, password string) bool {
	_, err := db.Engine.Exec("UPDATE watson_account SET name = ?, password = ? WHERE id = 1", username, password)
	return err == nil
}

// GetWatsonAccount : get account info
func GetWatsonAccount() (w WatsonAccount) {
	db.Engine.QueryRow("SELECT name, password FROM watson_account WHERE id = 1 LIMIT 1").Scan(&w.UserName, &w.Password)
	return w
}

// UpdateLanguage : update account info
func UpdateLanguage(model string) bool {
	_, err := db.Engine.Exec("UPDATE watson_account SET model = ? WHERE id = 1", model)
	fmt.Println(err)
	return err == nil
}

// GetLanguage : get account info
func GetLanguage() (model string) {
	db.Engine.QueryRow("SELECT model FROM watson_account WHERE id = 1 LIMIT 1").Scan(&model)
	return model
}

// GetWatsonToken : issue token for watson speech to text
func GetWatsonToken() string {
	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: cookieJar,
	}

	account := GetWatsonAccount()
	req, _ := http.NewRequest("GET", "https://stream.watsonplatform.net/authorization/api/v1/token?url=https://stream.watsonplatform.net/speech-to-text/api", nil)
	req.SetBasicAuth(account.UserName, account.Password)
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	return string(body)
}
