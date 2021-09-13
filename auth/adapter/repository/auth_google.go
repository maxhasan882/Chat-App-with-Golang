package repository

import (
	"chat/auth/domain/repository"
	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/google"
	"golang.org/x/oauth2"
	googleOAuth2 "golang.org/x/oauth2/google"
	"net/http"
)

type AuthGoogle struct{}

func (a AuthGoogle) NewGoogle(mux *http.ServeMux) *http.ServeMux {
	config := &repository.Config{
		ClientID:     "62477641542-u8lm0cv40os6fslpeph7m25p7nh70q73.apps.googleusercontent.com",
		ClientSecret: "QtWSiu6VOXK7_6KuHLEtZqLy",
	}
	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  "http://localhost:8080/auth_github/callback/google",
		Endpoint:     googleOAuth2.Endpoint,
		Scopes:       []string{"profile", "email"},
	}
	stateConfig := gologin.DebugOnlyCookieConfig
	mux.Handle("/google/login", google.StateHandler(stateConfig, google.LoginHandler(oauth2Config, nil)))
	mux.Handle("/auth_github/callback/google", google.StateHandler(stateConfig, google.CallbackHandler(oauth2Config, issueSessionGoogle(), nil)))
	return mux
}

func issueSessionGoogle() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		googleUser, err := google.UserFromContext(ctx)
		if err != nil {
			http.Error(w, "This is Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		session := repository.SessionStore.New(repository.SessionName)
		session.Values[repository.SessionUserKey] = googleUser.Id
		session.Values[repository.SessionUsername] = googleUser.Name
		session.Save(w)
		http.Redirect(w, req, "/", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}
