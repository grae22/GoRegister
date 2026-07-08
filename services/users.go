package services

import "goregister/domain"

type UsersService struct {
}

func NewUsersService() *UsersService {
	return &UsersService{}
}

func (s *UsersService) GetUsers() []domain.User {
	admin := domain.PermissionLogin |
		domain.PermissionViewAllEvents |
		domain.PermissionManageEvents |
		domain.PermissionDeleteRegisterEntry

	viewer := domain.PermissionLogin | domain.PermissionViewAllEvents

	u1, _ := domain.NewUser("graemeb", "Graeme Bruschi", "666", admin)
	u2, _ := domain.NewUser("geraldc", "Gerald Camp", "666", admin)
	u3, _ := domain.NewUser("dylanw", "Dylan Williams", "666", admin)
	u4, _ := domain.NewUser("waldob", "Waldo Bekker", "777", viewer)

	return []domain.User{*u1, *u2, *u3, *u4}
}

func (s *UsersService) GetUserById(id string) (domain.User, bool) {
	for _, u := range s.GetUsers() {
		if u.Id == id {
			return u, true
		}
	}

	return domain.User{}, false
}

func (s *UsersService) ValidatePassword(userId string, password string) bool {
	u, ok := s.GetUserById(userId)
	if !ok {
		return false
	}

	return u.VerifyPassword(password)
}

func (s *UsersService) HasPermission(userId string, p domain.UserPermissions) bool {
	u, ok := s.GetUserById(userId)
	if !ok {
		return false
	}

	return u.HasPermission(p)
}
