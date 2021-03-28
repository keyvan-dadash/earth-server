package user

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/scylladb/gocqlx/v2/table"
	"github.com/sod-lol/earth-server/libs/cassandra_table_managment"
	"golang.org/x/crypto/bcrypt"
)

var (
	userByUsernameTableName = "user_by_username"
	userByEmailTableName    = "user_by_email"
	userByUUIDTableName     = "user_by_uuid"
)

var userTableByUsernameMeta = table.Metadata{
	Name:    "user",
	Columns: []string{"username", "password", "email", "nickname", "uuid", "joineddate"},
	PartKey: []string{"username"},
	SortKey: []string{"username"},
}

var userTable = table.New(userTableByUsernameMeta)

var userTableByUsername = cassandra_table_managment.TableMetaData{
	Name: userByUsernameTableName,
	ColumnsAndTypes: map[string]string{
		"username":   "text",
		"password":   "text",
		"email":      "text",
		"nickname":   "text",
		"uuid":       "text",
		"joineddate": "timestamp"},
	PartKey: []string{"username"},
	SortKey: []string{},
}

var userTableByEmail = cassandra_table_managment.TableMetaData{
	Name: userByEmailTableName,
	ColumnsAndTypes: map[string]string{
		"username":   "text",
		"password":   "text",
		"email":      "text",
		"nickname":   "text",
		"uuid":       "text",
		"joineddate": "timestamp"},
	PartKey: []string{"email"},
	SortKey: []string{},
}

var userTableByUUID = cassandra_table_managment.TableMetaData{
	Name: userByUUIDTableName,
	ColumnsAndTypes: map[string]string{
		"username":   "text",
		"password":   "text",
		"email":      "text",
		"nickname":   "text",
		"uuid":       "text",
		"joineddate": "timestamp"},
	PartKey: []string{"uuid"},
	SortKey: []string{},
}

//User model
type User struct { //TODO: we should clean a little bit this akward field names and provide better field nameing
	Username   string    `participate:"true" kind:"pk" type:"text"`
	Password   string    `participate:"true" type:"text"`
	Email      string    `participate:"true" type:"text"`
	Nickname   string    `participate:"true" type:"text"`
	Uuid       string    `participate:"true" type:"text"`
	Joineddate time.Time `participate:"true" type:"timestamp"`
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
		Uuid:       uuid.NewV4().String(),
		Joineddate: time.Now(),
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
