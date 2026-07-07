package domain

import "testing"

func TestVerifyPassword_WhenCorrect_ThenReturnTrue(t *testing.T) {
	// Arrange.
	const pwd = "0passwordzZ"

	u, _ := NewUser(
		"123",
		"name",
		pwd)

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
		"abc")

	// Act.
	isValid := u.VerifyPassword(pwd)

	// Assert.
	if isValid {
		t.Error("Should be invalid")
	}
}
