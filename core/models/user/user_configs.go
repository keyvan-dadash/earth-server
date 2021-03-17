package user

import (
	"fmt"

	"github.com/scylladb/gocqlx/v2"
)

var CassandraAllowedSimpleDataType = []string{
	"boolean",
	"blob",
	"ascii",
	"bigint",
	"counter",
	"decimal",
	"double",
	"float",
	"int",
	"text",
	"varchar",
	"timestamp",
}

var userStruct = User{}

//CheckOrCreateUserTable is function that check if user table exist if not then going to cerate table
func CheckOrCreateUserTable(session *gocqlx.Session) (bool, error) {

	//TODO: we should invent mechanism to exmine filed of struct and
	// add field's to table without every time configuration

	// scanner := session.Query(`SELECT table_name
	// FROM system_schema.tables WHERE keyspace_name='earth' And table_name='user';`).Iter().Scanner()

	// if scanner.Next() {
	// 	return true, nil
	// }

	// results := session.Query(`CREATE TABLE user (
	// 	username text,
	// 	password text,
	// 	uuid text,
	// 	PRIMARY KEY(username)
	// 	)`).Exec()

	// return results == nil, results

	userRep := UserRepo{
		Session: session,
	}

	err := userRep.InsertUser(&User{
		Username: "ali",
		Password: "haqha",
	})

	if err != nil {
		fmt.Printf("we faced to error: %v", err)
	}

	return true, nil

}
