package utils

import (
	"goregister/domain"
	"goregister/services"
	"net/http"
)

type RequestCtx struct {
	User *domain.User
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
	c, err := r.Cookie("user")
	if err != nil {
		return
	}

	userId := c.Value
	u, ok := users.GetUserById(userId)
	if !ok {
		return
	}

	ctx.User = &u
}
