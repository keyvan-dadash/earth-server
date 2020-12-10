package models

import (
	"errors"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

//User model
type User struct {
	id       int
	username string
	password string
	uuid     uuid.UUID
}

//CreateUser create user instance based on given username and password
func CreateUser(username string, password string) (*User, error) {
	hashedPassword, err := hashAndSaltPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		username: username,
		password: hashedPassword,
		uuid:     uuid.NewV4(),
	}, nil
}

//GetUsername return user username
func (u *User) GetUsername() string {
	return u.username
}

//GetPassword return user password
func (u *User) GetPassword() string {
	return u.password
}

//VerifyPassword is function that verfiy given password
func (u *User) VerifyPassword(givenPassword string) bool {

	actualPassword := u.GetPassword()

	err := bcrypt.CompareHashAndPassword([]byte(actualPassword), []byte(givenPassword))
	return err == nil
}

func hashAndSaltPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", errors.New("cannot generate hash from password")
	}

	return string(hash), nil
}
