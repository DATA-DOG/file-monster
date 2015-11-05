package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/jmcvetta/napping.v1"
)

type oauthAccess struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	Ok          bool   `json:"ok"`
}

type authTest struct {
	UserID string `json:"user_id"`
}

func main() {
	loadConfig()
	openSQLConnection()

	router := httprouter.New()

	router.GET("/", index)
	router.POST("/slash", slashCommand)
	router.GET("/oauth", getOauth)

	log.Print("Listening on port ", config.Port, "...")
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", config.Port), router))
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.Header.Set("Content-Type", "text/html")

	fmt.Fprint(w, "Nothing to see here...")
}

func slashCommand(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := r.PostFormValue("user_id")
	text := r.PostFormValue("text")
	userName := r.PostFormValue("user_name")

	if userID == "" {
		fmt.Fprint(w, "No user provided")
		return
	}

	token := dbGetUserToken(userID)
	missingToken := len(token) == 0

	if text == "register" {
		if missingToken {
			fmt.Fprint(w, "Click on this link: <https://slack.com/oauth/authorize?client_id=", config.ClientID, "&scope=identify,read,post&state=1>")
			return
		}
		fmt.Fprint(w, "You are already registered")
		return
	}

	if missingToken {
		fmt.Fprint(w, "You are not registered. Type /flushme register")
		return
	}

	if text == "all" {
		fmt.Fprint(w, "Flushing public files...")
		go eatPublicFiles(token)
		return
	}

	fmt.Fprint(w, "Monster is starting to eat your files...")
	go eatUserFiles(userID, token, userName)
}

func getOauth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	error := r.URL.Query().Get("error")
	code := r.URL.Query().Get("code")

	if len(error) > 0 {
		fmt.Fprint(w, "Please accept authorisation request.")
		return
	}

	params := &napping.Params{"client_id": config.ClientID, "client_secret": config.ClientSecret, "code": code}
	access := oauthAccess{}
	var err interface{}

	napping.Get("https://slack.com/api/oauth.access", params, &access, err)

	if err != nil || !access.Ok {
		fmt.Fprint(w, "Invalid access token")
		return
	}

	params = &napping.Params{"token": access.AccessToken}
	test := authTest{}

	napping.Get("https://slack.com/api/auth.test", params, &test, err)

	if err != nil {
		fmt.Fprint(w, "Could not access user using access token")
		return
	}

	dbSaveUserToken(test.UserID, access.AccessToken)

	fmt.Fprint(w, "You may now close this window.")
}
