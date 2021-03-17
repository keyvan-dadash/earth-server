package config

import (
	"github.com/scylladb/gocqlx/v2"
	"github.com/sirupsen/logrus"
	"github.com/sod-lol/earth-server/core/models/user"
)

func InitializeTables(session *gocqlx.Session) (bool, error) {

	if _, err := user.CheckOrCreateUserTable(session); err != nil {
		logrus.Fatalf("[Fatal](InitializeTables) Something went wrong when cheking for existence of user table. error: %v", err)
		return false, err
	}

	return true, nil
}
