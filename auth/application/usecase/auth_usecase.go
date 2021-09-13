package usecase

import (
	"chat/auth/domain/repository"
	"net/http"
)

func DoGoogleAuthentication(a repository.IAuthGoogle, mux *http.ServeMux) *http.ServeMux {
	a.NewGoogle(mux)
	return mux
}

func DoGitHubAuthentication(a repository.IAuthGitHub, mux *http.ServeMux) *http.ServeMux {
	a.NewGitHub(mux)
	return mux
}
