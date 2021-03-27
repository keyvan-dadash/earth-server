package user

import (
	"reflect"
	"strings"

	"github.com/scylladb/gocqlx/v2"
	"github.com/sirupsen/logrus"
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

//CheckOrCreateUserTable is function that check if user table exist if not then going to cerate table
func CheckOrCreateUserTable(session *gocqlx.Session) (bool, error) {

	//TODO: we should invent mechanism to exmine filed of struct and
	// add field's to table without every time configuration

	rawSession := session.Session

	scanner := rawSession.Query(`SELECT table_name
	FROM system_schema.tables WHERE keyspace_name='earth' And table_name='user';`).Iter().Scanner()

	if !scanner.Next() {
		logrus.Info("[Info] table user does not exist begin ato create user table")
		results := rawSession.Query(
			buildCreateTableQuery(
				extractUserDesireFiledsFromUserStruct())).Exec()

		return results == nil, results
	}

	sliceMap, err := rawSession.Query(`SELECT column_name, type FROM system_schema.columns 
	WHERE keyspace_name = 'earth' AND table_name = 'user';`).Iter().SliceMap()

	if err != nil {
		logrus.Fatalf("[Fatal] cannot retrived sliceMap from user table columns. err: %v", err)
		return false, err
	}

	filedsTypes := extractUserTableFieldsFrom(sliceMap)

	desireFieldsTypes := extractUserDesireFiledsFromUserStruct()

	addColumnQuery, updateTypeQueries := buildUpdateQuery(
		filedsTypes, extractTypeOfFieldsFromDesireFields(desireFieldsTypes))

	//adding columns

	if addColumnQuery != "" {
		results := rawSession.Query(addColumnQuery).Exec()

		if results != nil {
			logrus.Fatalf("[Fatal] cannot execute add column query with query %v. err: %v", addColumnQuery, results)
			return results == nil, results
		}
	}

	//update columns types
	if len(updateTypeQueries) != 0 {
		for updateTypeQueryIndex := range updateTypeQueries {
			updateTypeQuery := updateTypeQueries[updateTypeQueryIndex]
			results := rawSession.Query(updateTypeQuery).Exec()

			if results != nil {
				logrus.Fatalf("[Fatal] cannot execute update column query with query %v. err: %v", updateTypeQuery, results)
				return results == nil, results
			}
		}
	}

	// userRep := UserRepo{
	// 	Session: session,
	// }

	// err = userRep.InsertUser(&User{
	// 	Username: "ali",
	// 	Password: "haqha",
	// })

	// if err != nil {
	// 	fmt.Printf("we faced to error: %v", err)
	// }

	return true, nil

}

func extractUserDesireFiledsFromUserStruct() map[string]reflect.StructField {

	userStruct := User{}

	userStructType := reflect.TypeOf(userStruct)

	userFieldsMap := map[string]reflect.StructField{}

	for i := 0; i < userStructType.NumField(); i++ {
		f := userStructType.Field(i)

		participate := f.Tag.Get("participate")
		capturedTag := f.Tag.Get("type")

		if len(capturedTag) == 0 || len(participate) == 0 || strings.ToLower(participate) == "false" {
			continue
		}

		userFieldsMap[strings.ToLower(f.Name)] = f

	}

	return userFieldsMap
}

func extractTypeOfFieldsFromDesireFields(fields map[string]reflect.StructField) map[string]string {

	mappedFieldTypes := make(map[string]string)

	for fieldName, field := range fields {
		fieldType := field.Tag.Get("type")

		mappedFieldTypes[fieldName] = fieldType
	}

	return mappedFieldTypes
}

func buildCreateTableQuery(fields map[string]reflect.StructField) string {

	query := "CREATE TABLE user (\n"

	primaryKeys := ""

	for fieldName, field := range fields {
		query += fieldName
		query += " "
		query += field.Tag.Get("type")

		kind := field.Tag.Get("kind")

		if len(kind) != 0 && strings.ToLower(kind) != "regular" {
			if len(primaryKeys) != 0 {
				primaryKeys += ", "
			}
			primaryKeys += fieldName
		}
		query += ",\n"
	}

	query += "PRIMARY KEY("
	query += primaryKeys
	query += ")"

	return query

}

func buildUpdateQuery(extractedFields, desireFields map[string]string) (string, []string) {
	addColumnQuery := "ALTER TABLE user ADD ("

	hasColumnUpdate := false

	updateColumnTypes := []string{}

	for columnName, desireType := range desireFields {

		extractedType, ok := extractedFields[columnName]

		if !ok {
			hasColumnUpdate = true
			addColumnQuery += columnName
			addColumnQuery += " "
			addColumnQuery += desireType
			addColumnQuery += ", "
			continue
		}

		if extractedType != desireType {
			updateColumm := ""
			updateColumm += "ALTER TABLE user ALTER " + columnName + " TYPE " + desireType + ";"

			updateColumnTypes = append(updateColumnTypes, updateColumm)
		}

	}

	addColumnQuery = addColumnQuery[:len(addColumnQuery)-2] //erase ', '

	addColumnQuery += ")"

	if !hasColumnUpdate {
		return "", updateColumnTypes
	}

	return addColumnQuery, updateColumnTypes

}

func extractUserTableFieldsFrom(columnsAndTypes []map[string]interface{}) map[string]string {

	filedsTypes := make(map[string]string)
	for row := range columnsAndTypes {
		retrRow := columnsAndTypes[row]
		rowName := strings.ToLower(retrRow["column_name"].(string))
		filedsTypes[rowName] = retrRow["type"].(string)
	}

	return filedsTypes
}
