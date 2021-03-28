package user

import (
	"github.com/scylladb/gocqlx/v2/table"
	"github.com/sod-lol/earth-server/libs/cassandra_table_managment"
)

var (
	userByUsernameTableName = "user_by_username"
	userByEmailTableName    = "user_by_email"
	userByUUIDTableName     = "user_by_uuid"
)

//table user with username as priamry key
var userTableByUsernameMetaData = cassandra_table_managment.TableMetaData{
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

//table user with email as priamry key
var userTableByEmailMetaData = cassandra_table_managment.TableMetaData{
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

//table user with uuid as priamry key
var userTableByUUIDMetaData = cassandra_table_managment.TableMetaData{
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

//user table orm by username
var userTableByUsernameMeta = table.Metadata{
	Name:    "user",
	Columns: []string{"username", "password", "email", "nickname", "uuid", "joineddate"},
	PartKey: []string{"username"},
	SortKey: []string{"username"},
}

var userTableByUsername = table.New(userTableByUsernameMeta)

//user table orm by emailuserTableByUsernameMeta
var userTableByEmailMeta = table.Metadata{
	Columns: []string{"username", "password", "email", "nickname", "uuid", "joineddate"},
	PartKey: []string{"username"},
	SortKey: []string{"username"},
}

var userTableByEmail = table.New(userTableByEmailMeta)

//user table orm by uuid
var userTableByUUIDMeta = table.Metadata{
	Name:    "user",
	Columns: []string{"username", "password", "email", "nickname", "uuid", "joineddate"},
	PartKey: []string{"username"},
	SortKey: []string{"username"},
}

var userTableByUUID = table.New(userTableByUUIDMeta)
