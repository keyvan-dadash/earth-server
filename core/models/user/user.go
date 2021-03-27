package user

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/scylladb/gocqlx/v2/table"
	"golang.org/x/crypto/bcrypt"
)

var userTableMeta = table.Metadata{
	Name:    "user",
	Columns: []string{"username", "password"},
	PartKey: []string{"username"},
	SortKey: []string{"username"},
}

var userTable = table.New(userTableMeta)

//User model
type User struct {
	Username   string    `participate:"true" kind:"pk" type:"text"`
	Password   string    `participate:"true" type:"text"`
	Email      string    `participate:"true" type:"text"`
	Nickname   string    `participate:"true" type:"text"`
	UUID       uuid.UUID `participate:"true" type:"varchar"`
	JoinedDate time.Time
}

//CreateUser create user instance based on given username and password
func CreateUser(username string, password string) (*User, error) {
	hashedPassword, err := hashAndSaltPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		Username:   username,
		Password:   hashedPassword,
		Nickname:   username,
		UUID:       uuid.NewV4(),
		JoinedDate: time.Now(),
	}, nil
}

//GetUsername return user username
func (u *User) GetUsername() string {
	return u.Username
}

//GetPassword return user password
func (u *User) GetPassword() string {
	return u.Password
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
