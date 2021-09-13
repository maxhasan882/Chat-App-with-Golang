package repository

import (
	"github.com/dghubble/sessions"
	"net/http"
)

const (
	SessionSecret   = "cookie signing secret"
	SessionName     = "Chat App Go Lang"
	SessionUserKey  = "userID"
	SessionUsername = "userName"
)

var SessionStore = sessions.NewCookieStore([]byte(SessionSecret), nil)

type Config struct {
	ClientID     string
	ClientSecret string
}

type IAuthGoogle interface {
	NewGoogle(mux *http.ServeMux) *http.ServeMux
}

type IAuthGitHub interface {
	NewGitHub(mux *http.ServeMux) *http.ServeMux
}
