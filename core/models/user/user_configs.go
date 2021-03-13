package user

import (
	"github.com/gocql/gocql"
)

//CheckOrCreateUserTable is function that check if user table exist if not then going to cerate table
func CheckOrCreateUserTable(session *gocql.Session) (bool, error) {

	scanner := session.Query(`SELECT table_name 
	FROM system_schema.tables WHERE keyspace_name='earth' And table_name='user';`).Iter().Scanner()

	if scanner.Next() {
		return true, nil
	}

	results := session.Query(`CREATE TABLE user (
		username text,
		password text,
		uuid text,
		PRIMARY KEY(username)
		)`).Exec()

	return results == nil, results

}
