package domain

import (
	"errors"
	"strings"
)

type User struct {
	Id      string
	Name    string
	pwdHash string
}

func NewUser(
	id string,
	name string,
	pwdHash string,
) (*User, error) {

	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)

	if id == "" {
		return nil, errors.New("Invalid id")
	}

	if name == "" {
		return nil, errors.New("Invalid name")
	}

	o := User{
		Id:      id,
		Name:    name,
		pwdHash: pwdHash,
	}

	return &o, nil
}
