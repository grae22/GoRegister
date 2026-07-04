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
