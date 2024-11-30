package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var oauth2Config = oauth2.Config{
	ClientID:     "<842701492130-od8gib8hvm4tomb8luqhkdv0rq6avh5a.apps.googleusercontent.com>",
	ClientSecret: "<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>",
	RedirectURL:  "http://localhost:8080/auth/login",
	Scopes: []string{
		"openid",
		"profile",
		"email",
	},
	Endpoint: google.Endpoint,
}

var oauthStateString = "randomstate"

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {

	url := oauth2Config.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code")

	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to exchange token: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	client := oauth2Config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user info: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode user info: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User Info: %+v", userInfo)
}

func main() {
	// Set up routing
	r := mux.NewRouter()
	r.HandleFunc("/auth/login", handleGoogleLogin)
	r.HandleFunc("/auth/callback", handleGoogleCallback)

	// Start the server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
