package domain

import (
	"errors"
	"strings"
)

type UserPermissions int

const (
	PermissionNone                  = 0
	PermissionLogin UserPermissions = 1 << iota
	PermissionViewAllEvents
	PermissionManageEvents
	PermissionDeleteRegisterEntry
)

const guestId string = "guest"

type User struct {
	Id          string
	Name        string
	Permissions UserPermissions
	IsGuest     bool
	password    string
}

func NewUser(
	id string,
	name string,
	password string,
	permissions UserPermissions,
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
		Id:          id,
		Name:        name,
		Permissions: permissions,
		IsGuest:     id == guestId,
		password:    encrypt(password),
	}

	return &o, nil
}

func NewGuestUser() User {
	return User{
		Id:          guestId,
		Name:        "Guest",
		Permissions: PermissionNone,
		IsGuest:     true,
		password:    encrypt("guest"),
	}
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

func (u *User) HasPermission(p UserPermissions) bool {
	return u.Permissions&p == p
}

func encrypt(s string) string {
	var newS string

	for i := range len(s) {
		newS += string(s[i] - 32)
	}

	return newS
}
