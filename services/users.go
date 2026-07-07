package services

import "goregister/domain"

type UsersService struct {
}

func NewUsersService() *UsersService {
	return &UsersService{}
}

func (s *UsersService) GetUsers() []domain.User {
	u1, _ := domain.NewUser("graemeb", "Graeme Bruschi", "666")
	u2, _ := domain.NewUser("geraldc", "Gerald Camp", "666")
	u3, _ := domain.NewUser("dylanw", "Dylan Williams", "666")

	return []domain.User{*u1, *u2, *u3}
}

func (s *UsersService) GetUserById(id string) (domain.User, bool) {
	for _, u := range s.GetUsers() {
		if u.Id == id {
			return u, true
		}
	}

	return domain.User{}, false
}

func (s *UsersService) ValidatePassword(username string, password string) bool {
	u, ok := s.GetUserById(username)
	if !ok {
		return false
	}

	return u.VerifyPassword(password)
}
