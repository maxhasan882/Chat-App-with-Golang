package auth

import (
	adapterRepo "chat/auth/adapter/repository"
	"chat/auth/application/usecase"
	"chat/auth/domain/repository"
	"net/http"
)

var googleAuthRepository = adapterRepo.AuthGoogle{}
var githubAuthRepository = adapterRepo.AuthGitHub{}

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("gologin-temporary-cookie")
	if err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}
	h.next.ServeHTTP(w, r)
}

func New(mux *http.ServeMux) *http.ServeMux {
	mux = usecase.DoGoogleAuthentication(googleAuthRepository, mux)
	mux = usecase.DoGitHubAuthentication(githubAuthRepository, mux)
	mux.HandleFunc("/logout", logoutHandler)
	return mux
}

func logoutHandler(w http.ResponseWriter, req *http.Request) {
	repository.SessionStore.Destroy(w, repository.SessionName)
	http.Redirect(w, req, "/login", http.StatusFound)
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}
