package web

import (
	"goregister/domain"
	"goregister/utils"
)

type Layout struct {
	HasTopBar         bool
	ShowAllEventsLink bool
	CurrentUser       domain.User
}

func NewLayout(
	hasTopBar bool,
	ctx utils.RequestCtx,
) Layout {
	return Layout{
		HasTopBar:         hasTopBar,
		ShowAllEventsLink: ctx.User.HasPermission(domain.PermissionViewAllEvents),
		CurrentUser:       ctx.User,
	}
}
