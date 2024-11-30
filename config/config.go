package config

import "golang.org/x/oauth2"

var OAuthConfig = oauth2.Config{
	ClientID:     "your-client-id.apps.googleusercontent.com", // Replace with your Google OAuth credentials
	ClientSecret: "your-client-secret",                        // Replace with your Google OAuth credentials
	RedirectURL:  "http://localhost:8080/auth/login",          // Set this in Google Developer Console
	Scopes: []string{
		"openid",
		"profile",
		"email",
	},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://accounts.google.com/o/oauth2/auth",
		TokenURL: "https://accounts.google.com/o/oauth2/token",
	},
}

var OAuthStateString = "randomstate"
