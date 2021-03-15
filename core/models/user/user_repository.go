package user

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
)

type UserRepoInterface interface {
	InsertUser(*User) error
	RetrieveUser(string) (*User, error)
	UpdateUser(string, *User) error
	DeleteUser(string) error
}

type UserRepo struct {
	*gocql.Session
}

//InsertUser function will try to insert to db
func (ur *UserRepo) InsertUser(user *User) error {
	err := ur.Session.Query(`INSERT INTO user (username, password, email, nickname, uuid, isonline, joineddate) VALUES 
						(?, ?, ?, ?, ?, ?, now()) IF NOT EXISTS`).Exec()

	if err != nil {
		logrus.Debugf("User with username %v already exist", user.username)
		return err
	}

	return err
}

func (ur *UserRepo) RetrieveUser(username string) (*User, error) {
	retrievedMap, err := ur.Session.Query(`SELECT * FROM user WHERE username = ?`, username).Iter().SliceMap()

	if err != nil {
		return nil, err
	}

	if len(retrievedMap) > 1 {
		logrus.Fatal("[Fatal](RetrieveUser) Found two user with same username")
		return nil, fmt.Errorf("Found two user with same username")
	}

	user := &User{}

	if err := ParseUserFromMap(retrievedMap[0], user); err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepo) UpdateUser(username string, updatedUser *User) error {

	err := ur.Session.Query(`
			UPDATE user 
				SET username=?,
					password=?,
					email = ?,
					nickname = ?,
					isOnline = ?
			WHERE username = ?`, updatedUser.username, updatedUser.password, updatedUser.Email, updatedUser.Nickname, updatedUser.IsOnline, username)

	if err != nil {
		return fmt.Errorf("updating user faild. err: %v", err)
	}

	return nil
}

func (ur *UserRepo) DeleteUser(username string) error {
	err := ur.Session.Query(`
			DELETE FROM user WHERE username = ?`, username)

	if err != nil {
		return fmt.Errorf("deleting user faild. err: %v", err)
	}

	return nil
}
