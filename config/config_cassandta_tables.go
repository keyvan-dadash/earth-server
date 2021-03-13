package config

import (
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
	"github.com/sod-lol/earth-server/core/models/user"
)

func InitializeTables(session *gocql.Session) (bool, error) {

	if _, err := user.CheckOrCreateUserTable(session); err != nil {
		logrus.Fatalf("[Fatal](InitializeTables) Something went wrong when cheking for existence of user table. error: %v", err)
		return false, err
	}

	return true, nil
}
