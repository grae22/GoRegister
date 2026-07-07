package web

import (
	"errors"
	"goregister/domain"
	"goregister/services"
	"goregister/utils"
	"html/template"
	"net/http"
)

type UsersController struct {
	usersService *services.UsersService
}

type loginPageData struct {
	CurrentUser *domain.User
	HasFailed   bool
}

func NewUsersController(
	usersService *services.UsersService,
) (*UsersController, error) {

	if usersService == nil {
		return nil, errors.New("Nil users service")
	}

	c := UsersController{
		usersService: usersService,
	}

	return &c, nil
}

func (c *UsersController) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	requestCtx := utils.NewRequestContext(c.usersService, r)

	data := loginPageData{
		CurrentUser: requestCtx.User,
		HasFailed:   r.URL.Query().Has("failed"),
	}

	tmpl := template.Must(template.ParseFiles("html/layout.html", "html/login.html"))
	tmpl.ExecuteTemplate(w, "layout", data)
}
