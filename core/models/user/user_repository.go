package user

import (
	"errors"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
)

type UserRepoInterface interface {
	InsertUser(*User) error
	RetrieveUser(string) (*User, error)
	UpdateUser(string, *User) error
	DeleteUser(string) error
}

type UserRepo struct {
	*gocqlx.Session
}

//InsertUser function will try to insert to db
func (ur *UserRepo) InsertUser(user *User, inserdUnique bool) error {
	if !inserdUnique {
		err := ur.Session.Query(userTableByUsername.Insert()).BindStruct(*user).ExecRelease()

		if err != nil {
			return err
		}

		err = ur.Session.Query(userTableByEmail.Insert()).BindStruct(*user).ExecRelease()

		if err != nil {
			return err
		}

		err = ur.Session.Query(userTableByUUID.Insert()).BindStruct(*user).ExecRelease()

		if err != nil {
			return err
		}

		return nil

	}

	if err := ur.RetrieveUser(user); err != nil {
		if !errors.Is(err, gocql.ErrNotFound) {
			return err
		}

		err := ur.Session.Query(userTableByUsername.Insert()).BindStruct(*user).ExecRelease()

		if err != nil {
			return err
		}

		err = ur.Session.Query(userTableByEmail.Insert()).BindStruct(*user).ExecRelease()

		if err != nil {
			return err
		}

		err = ur.Session.Query(userTableByUUID.Insert()).BindStruct(*user).ExecRelease()

		if err != nil {
			return err
		}

		return nil

	}

	return ErrDublicateUser

}

func (ur *UserRepo) RetrieveUser(user *User) error {
	q := qb.Select(userByUsernameTableName).Where(qb.Eq("username")).Query(*ur.Session).BindStruct(user)

	return q.Get(user)
}

func (ur *UserRepo) UpdateUser(username string, updatedUser *User, updatedColumns ...string) (bool, error) {
	if ok, err := ur.Session.Query(userTableByUsername.Update(updatedColumns...)).BindStruct(*updatedUser).ExecCASRelease(); err != nil || !ok {
		return ok, err
	}

	if ok, err := ur.Session.Query(userTableByEmail.Update(updatedColumns...)).BindStruct(*updatedUser).ExecCASRelease(); err != nil || !ok {
		return ok, err
	}

	if ok, err := ur.Session.Query(userTableByUUID.Update(updatedColumns...)).BindStruct(*updatedUser).ExecCASRelease(); err != nil || !ok {
		return ok, err
	}

	return true, nil
}

func (ur *UserRepo) DeleteUser(username string) error {
	user := User{
		Username: username,
	}
	ur.RetrieveUser(&user)

	if err := ur.Session.Query(userTableByUsername.Delete()).BindStruct(user).ExecRelease(); err != nil {
		return err
	}

	if err := ur.Session.Query(userTableByEmail.Delete()).BindStruct(user).ExecRelease(); err != nil {
		return err
	}

	if err := ur.Session.Query(userTableByUUID.Delete()).BindStruct(user).ExecRelease(); err != nil {
		return err
	}

	return nil
}
