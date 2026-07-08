package utils

import (
	"goregister/domain"
	"goregister/services"
	"net/http"
)

type RequestCtx struct {
	User domain.User
}

func NewRequestContext(
	users *services.UsersService,
	r *http.Request,
) *RequestCtx {
	ctx := RequestCtx{}

	ctx.RetrieveUser(users, r)

	return &ctx
}

func (ctx *RequestCtx) RetrieveUser(users *services.UsersService, r *http.Request) {
	guest := domain.User{
		Id:          "guest",
		Name:        "Guest",
		Permissions: domain.PermissionNone,
	}

	c, err := r.Cookie("user")
	if err != nil {
		ctx.User = guest
		return
	}

	userId := c.Value
	u, ok := users.GetUserById(userId)
	if !ok {
		ctx.User = guest
		return
	}

	ctx.User = u
}
