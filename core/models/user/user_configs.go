package user

import (
	"github.com/scylladb/gocqlx/v2"
	"github.com/sirupsen/logrus"
	"github.com/sod-lol/earth-server/libs/cassandra_table_managment"
)

var CassandraAllowedSimpleDataType = []string{
	"ascii",
	"bigint",
	"blob",
	"boolean",
	"counter",
	"date",
	"decimal",
	"double",
	"duration",
	"float",
	"inet",
	"int",
	"smallint",
	"text",
	"time",
	"timestamp",
	"timeuuid",
	"tinyint",
	"uuid",
	"varchar",
	"varint",
}

var UserRepository *UserRepo

//CheckOrCreateUserTable is function that check if user table exist if not then going to cerate table
func CheckOrCreateUserTable(session *gocqlx.Session) (bool, error) {

	rawSession := session.Session

	UserRepository = &UserRepo{
		Session: session,
	}

	tableManagment := cassandra_table_managment.TableManagment{}

	tableManagment.SetSession(rawSession)

	tableManagment.SetKeySpace("earth")

	tableManagment.AddTableMetaData(userTableByUsernameMetaData)
	tableManagment.AddTableMetaData(userTableByEmailMetaData)
	tableManagment.AddTableMetaData(userTableByUUIDMetaData)

	err := tableManagment.CheckExistanceOfTables()

	if err != nil {
		logrus.Fatalf("[Fatal] falild to iniliaze user tables during checking existance of tables. err: %v", err)
		return false, err
	}

	err = tableManagment.UpdateTables()

	if err != nil {
		logrus.Fatalf("[Fatal] falild to iniliaze user tables during updating tables. err: %v", err)
		return false, err
	}

	return true, nil

}
