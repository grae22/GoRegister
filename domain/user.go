package domain

import (
	"errors"
	"strings"
)

type User struct {
	Id       string
	Name     string
	password string
}

func NewUser(
	id string,
	name string,
	password string,
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
		Id:       id,
		Name:     name,
		password: encrypt(password),
	}

	return &o, nil
}

func (u *User) VerifyPassword(pwd string) bool {
	if len(pwd) != len(u.password) {
		return false
	}

	ePwd := encrypt(pwd)

	for i := range len(ePwd) {
		if ePwd[i] != u.password[i] {
			return false
		}
	}

	return true
}

func encrypt(s string) string {
	var newS string

	for i := range len(s) {
		newS += string(s[i] - 32)
	}

	return newS
}
