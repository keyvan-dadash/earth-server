package user

import "github.com/gocql/gocql"

type UserRepo struct {
	*gocql.Session
}

//InsertUser
// func (ur *UserRepo) InsertUser(user *User) (bool, error) {

// }
