package domain

import "testing"

func TestVerifyPassword_WhenCorrect_ThenReturnTrue(t *testing.T) {
	// Arrange.
	const pwd = "0passwordzZ"

	u, _ := NewUser(
		"123",
		"name",
		pwd,
		PermissionNone)

	// Act.
	isValid := u.VerifyPassword(pwd)

	// Assert.
	if !isValid {
		t.Error("Should be valid")
	}
}

func TestVerifyPassword_WhenIncorrect_ThenReturnFalse(t *testing.T) {
	// Arrange.
	const pwd = "0passwordzZ"

	u, _ := NewUser(
		"123",
		"name",
		"abc",
		PermissionNone)

	// Act.
	isValid := u.VerifyPassword(pwd)

	// Assert.
	if isValid {
		t.Error("Should be invalid")
	}
}

func TestHasPermission_WhenHasPermission_ThenReturnsTrue(t *testing.T) {
	// Arrange.
	u, _ := NewUser(
		"123",
		"name",
		"abc",
		PermissionManageEvents|PermissionDeleteRegisterEntry)

	// Act.
	hasP := u.HasPermission(PermissionManageEvents)
	hasP = hasP && u.HasPermission(PermissionDeleteRegisterEntry)

	// Assert.
	if !hasP {
		t.Error("Should have all permissions")
	}
}

func TestHasPermission_WhenDoesNotHavePermission_ThenReturnsFalse(t *testing.T) {
	// Arrange.
	u, _ := NewUser(
		"123",
		"name",
		"abc",
		PermissionManageEvents|PermissionDeleteRegisterEntry)

	// Act.
	hasP := u.HasPermission(PermissionLogin)

	// Assert.
	if hasP {
		t.Error("Should not have permission")
	}
}
