package main

import (
	"errors"
	"slices"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Name           string
	Email          string
	HashedPassword []byte
}

var usersDB = []user{
	{
		Name:           "bob",
		Email:          "bob@google.com",
		HashedPassword: hashPassword("bobpassword"),
	},
	{
		Name:           "alice",
		Email:          "alice@google.com",
		HashedPassword: hashPassword("alicepassword"),
	},
}

var ErrInvalidUserOrPassword = errors.New("invalid user or password")

var dummyHash = []byte("$2a$10$7EqJtq98hPqEX7fNZaFWoOhi5pQ9qvZ1Wg0R8rFqYzjO3y8Yp0m3W")

func authUser(email, password string) (user, error) {
	// Ищем пользователя
	idx := slices.IndexFunc(usersDB, func(u user) bool {
		return strings.EqualFold(email, u.Email)
	})

	// По умолчанию проверяем пароль на dummyHash
	hashToCheck := dummyHash
	found := false

	var usr user
	if idx != -1 {
		usr = usersDB[idx]
		hashToCheck = usr.HashedPassword
		found = true
	}

	// ВАЖНО: checkPassword вызываем всегда
	// (даже если email не найден), чтобы закрыть CWE-208.
	if err := checkPassword(hashToCheck, password); err != nil || !found {
		return user{}, ErrInvalidUserOrPassword
	}

	return usr, nil
}

func hashPassword(password string) []byte {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hash
}

func checkPassword(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}
