package webapi

import (
	"errors"
	"goregister/services"
	"net/http"
	"strings"
)

type UserSessionApi struct {
	users *services.UsersService
}

func NewUserSessionApi(users *services.UsersService) (*UserSessionApi, error) {
	if users == nil {
		return nil, errors.New("Nil users service")
	}

	return &UserSessionApi{
			users: users,
		},
		nil
}

func (a *UserSessionApi) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		a.handleLogin(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *UserSessionApi) handleLogin(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.FormValue("username"))
	if len(username) == 0 {
		http.Error(w, "No username supplied", http.StatusBadRequest)
		return
	}

	password := r.FormValue("password")
	if len(password) == 0 {
		http.Error(w, "No password supplied", http.StatusBadRequest)
		return
	}

	if a.users.ValidatePassword(username, password) {
		c := http.Cookie{
			Name:     "user",
			Value:    username,
			Path:     "/",
			MaxAge:   2678400,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}

		http.SetCookie(w, &c)

		http.Redirect(
			w,
			r,
			"/events",
			http.StatusSeeOther)
	} else {
		http.Redirect(
			w,
			r,
			"/login?failed",
			http.StatusSeeOther)
	}
}
