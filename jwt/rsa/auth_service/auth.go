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

func authUser(email, password string) (user, error) {
	// ищем пользователя в БД
	idx := slices.IndexFunc(usersDB, func(u user) bool {
		return strings.EqualFold(email, u.Email)
	})
	if idx == -1 {
		return user{}, ErrInvalidUserOrPassword
	}

	usr := usersDB[idx]

	// сравниваем пароли
	if err := checkPassword(usr.HashedPassword, password); err != nil {
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
