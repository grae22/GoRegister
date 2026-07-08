package webapi

import (
	"errors"
	"goregister/services"
	"net/http"
	"strings"
	"time"
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
	switch r.Method {
	case http.MethodPost:
		a.handleLogin(w, r)

	case http.MethodDelete:
		a.handleLogout(w, r)

	default:
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
		c := createUserCookie(username)

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

func (a *UserSessionApi) handleLogout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("user")
	if err == nil {
		newC := createUserCookie(c.Value)
		newC.Expires = time.Unix(0, 0)
		newC.MaxAge = -1
		http.SetCookie(w, &newC)
	}

	w.WriteHeader(http.StatusOK)
}

func createUserCookie(userId string) http.Cookie {
	return http.Cookie{
		Name:     "user",
		Value:    userId,
		Path:     "/",
		MaxAge:   2678400,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}
