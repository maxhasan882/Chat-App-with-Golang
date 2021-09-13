package repository

import (
	"chat/auth/domain/repository"
	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/github"
	"golang.org/x/oauth2"
	githubOAuth2 "golang.org/x/oauth2/github"
	"net/http"
)

type AuthGitHub struct{}

func (a AuthGitHub) NewGitHub(mux *http.ServeMux) *http.ServeMux {
	config := &repository.Config{
		ClientID:     "01d1f87de60727244b99",
		ClientSecret: "c2d7896ca7643259ecdea78b66b6b9edc7a7a254",
	}
	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  "http://118.179.96.252:8080/auth/callback/github",
		Endpoint:     githubOAuth2.Endpoint,
	}
	stateConfig := gologin.DebugOnlyCookieConfig
	mux.Handle("/github/login", github.StateHandler(stateConfig, github.LoginHandler(oauth2Config, nil)))
	mux.Handle("/auth/callback/github", github.StateHandler(stateConfig, github.CallbackHandler(oauth2Config, issueSessionGitHub(), nil)))
	return mux
}

func issueSessionGitHub() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		githubUser, err := github.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session := repository.SessionStore.New(repository.SessionName)
		session.Values[repository.SessionUserKey] = *githubUser.ID
		session.Values[repository.SessionUsername] = *githubUser.Login
		session.Save(w)
		http.Redirect(w, req, "/", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}
